/*

func         PoolWrapper (func (interface{})                     ) *Pool
func (*Pool) WithCount(min   uint         , max     uint         ) *Pool
func (*Pool) WithTime (delay time.Duration, timeout time.Duration) *Pool

func (*Pool) Call     (interface{})

*/
package svc

import (
	"github.com/moxitrel/golib"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// Wrap time.Timer
type Timer struct {
	*time.Timer
}

func NewTimer() (o Timer) {
	o = Timer{
		Timer: time.NewTimer(-1),
	}
	o.Stop()
	return
}
func (o Timer) Start(timeout time.Duration) {
	o.Reset(timeout)
}
func (o Timer) Stop() {
	if !o.Timer.Stop() {
		<-o.C
	}
}

const (
	_POOL_WORKER_MIN   = 0
	_POOL_WORKER_MAX   = math.MaxUint16
	_POOL_CALL_DELAY   = 200 * time.Millisecond
	_POOL_TICKER_INTVL = 45 * time.Second

	// time to wait for receiving sent args when receive stop signal
	_STOP_DELAY = 200 * time.Millisecond
)

// Start [min, max] goroutines of <Pool.fun> to process <Pool.arg>
//
// * Example
// f := func(x interface{}) { time.Sleep(time.Second) }
// p := PoolWrapper(f)	    // start 1 goroutines of f
// p.Call("1")			// run f("1") in background and return immediately
// p.Call("2")			// run f("2") in background after block <Pool.delay> ns
// p.Call("3")			// run f("3") in background after block <Pool.delay> ns
//
type Pool struct {
	// at least <min> workers will be created and live all the time
	min int32
	// the max number of workers can be created
	max int32
	// the current number of workers
	cur int32
	// the number of idle workers
	freeCount int32
	// create a new worker if <arg> is blocked for <delay> ns, a proper value should be >= 0.1s
	delay time.Duration
	// destroy the workers which idle for <timeout> ns
	//timeout time.Duration

	fun func(interface{})
	arg chan interface{}

	// send timeout-signal to worker periodically
	workerTicker *Svc
	// How many signals to sent out after a tick.
	// Control how many idle workers to close after a timeout.
	maxTickSignal int32
	// decrease worker's life when tick-signal received, timeout when life <= 0
	life int32
	// send timeout-signal to submitter periodically through submitterTickerSignal
	submitterTicker       *Svc
	submitterTickerSignal chan struct{}
	// how many .Call() blocked by *ticker* (when cur < max)
	blockCount int32

	// lock cur when updating
	curLock  sync.Mutex
	stopOnce sync.Once
	wg       sync.WaitGroup
}

type _TickSignal struct{}

var tickSignal _TickSignal

func (o *Pool) getMax() int32 {
	return atomic.LoadInt32(&o.max)
}
func (o *Pool) getCur() int32 {
	return atomic.LoadInt32(&o.cur)
}
func (o *Pool) getDelay() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&o.delay)))
}
func (o *Pool) getLife() int32 {
	return atomic.LoadInt32(&o.life)
}

func NewPool(min, max uint, delay, timeout time.Duration, bufferSize uint, fun func(interface{})) (o *Pool) {
	if min > math.MaxInt32 {
		golib.Panic("min:%v > math.MaxInt32, want <= math.MaxInt32", min)
	}
	if max > math.MaxInt32 {
		golib.Panic("max:%v > math.MaxInt32, want <= math.MaxInt32", max)
	}
	if min > max {
		golib.Panic("min:%v > max:%v, want min <= max", min, max)
	}
	if delay < 0 {
		golib.Panic("delay:%v < 0, want >= 0", delay)
	}
	if timeout < 0 {
		golib.Panic("timeout:%v < 0, want >= 0", timeout)
	}
	if fun == nil {
		golib.Panic("fun == nil, want !nil")
	}

	o = &Pool{
		min:       int32(min),
		max:       int32(max),
		cur:       0,
		freeCount: 0,
		delay:     delay,

		fun: fun,
		arg: make(chan interface{}, bufferSize),

		workerTicker:          nil,          // inited below
		maxTickSignal:         math.MaxInt8, // default number of tickSignal to send per tick
		life:                  2,            // default life
		submitterTicker:       nil,          // inited below
		submitterTickerSignal: make(chan struct{}),
		blockCount:            0,

		curLock:  sync.Mutex{},
		stopOnce: sync.Once{},
		wg:       sync.WaitGroup{},
	}
	// init worker ticker
	workerTicker := time.NewTicker(timeout)
	o.wg.Add(1)
	o.workerTicker = NewSvc(nil, func() {
		<-workerTicker.C
		// stop idle workers
		n := atomic.LoadInt32(&o.freeCount)
		maxTickSignal := atomic.LoadInt32(&o.maxTickSignal)
		if n > maxTickSignal {
			n = maxTickSignal
		}
		for ; n > 0; n-- {
			select {
			case o.arg <- tickSignal:
			default:
				return
			}
		}
	}, func() {
		o.wg.Done()
		workerTicker.Stop()
	})

	// init submitter ticker
	if o.delay > 0 {
		submitterTicker := time.NewTicker(delay)
		o.wg.Add(1)
		o.submitterTicker = NewSvc(nil, func() {
			<-submitterTicker.C
			for n := atomic.LoadInt32(&o.blockCount); n > 0; n-- {
				select {
				case o.submitterTickerSignal <- struct{}{}:
				default:
					return
				}
			}
		}, func() {
			o.wg.Done()
			submitterTicker.Stop()
		})
	}

	// if timeout is 0, new process will exit immediately which decrease the cur.
	for i := o.min; i > 0; i-- {
		o.newProcess()
	}
	return
}

func PoolWrapper(fun func(interface{})) (func(interface{}), func()) {
	v := NewPool(_POOL_WORKER_MIN, _POOL_WORKER_MAX, 0, _POOL_TICKER_INTVL, 0, fun)
	return v.Call, v.Stop
}

// Create a new goroutine.
func (o *Pool) newProcess() {
	// NOTE: Additionally use atomic operator for o.cur to avoid data race from .Submitter() which has no lock.
	//       Only lock o.cur when updating.
	o.curLock.Lock()
	if o.getCur() >= o.getMax() {
		o.curLock.Unlock()
		return
	}
	atomic.AddInt32(&o.cur, 1)
	o.curLock.Unlock()

	o.wg.Add(1)
	go func() {
		defer o.wg.Done()
		life := atomic.LoadInt32(&o.life)
		for {
			atomic.AddInt32(&o.freeCount, 1)
			arg := <-o.arg
			atomic.AddInt32(&o.freeCount, -1)
			switch arg {
			case tickSignal:
				life--
				if life <= 0 {
					o.curLock.Lock()
					switch {
					case o.getCur() <= atomic.LoadInt32(&o.min): // continue to work
						o.curLock.Unlock()
						life = atomic.LoadInt32(&o.life)
					default: // quit
						atomic.AddInt32(&o.cur, -1)
						o.curLock.Unlock()
						return
					}
				}
			default:
				life = atomic.LoadInt32(&o.life)
				o.fun(arg)
			}
		}
	}()
}

func (o *Pool) Stop() {
	o.stopOnce.Do(func() {
		atomic.StoreInt32(&o.min, 0)
		atomic.StoreInt32(&o.max, 0)
		atomic.StoreInt32(&o.life, 1) // stop when receive a tickSignal

		o.workerTicker.Stop()
		if o.submitterTicker != nil {
			o.submitterTicker.Stop()
		}

		go func() {
			for o.getCur() > 0 {
				time.Sleep(_STOP_DELAY)
				for blockCount := atomic.LoadInt32(&o.blockCount); blockCount > 0; blockCount-- {
					select {
					case o.submitterTickerSignal <- struct{}{}:
					default:
					}
				}
				for freeCount := atomic.LoadInt32(&o.freeCount); freeCount > 0; freeCount-- {
					select {
					case o.arg <- tickSignal:
					default:
					}
				}
			}
		}()
	})
}

func (o *Pool) Join() {
	o.wg.Wait()
}

func (o *Pool) Call(arg interface{}) {
	life := atomic.LoadInt32(&o.life)
	for {
		switch max := o.getMax(); {
		case max <= 0: // stopped
			return
		case o.getCur() >= max: // wait forever
			// NOTE: use atomic instead of lock for o.cur for performance
			o.arg <- arg
			return
		default:
			select {
			case o.arg <- arg: // skip select if not blocked
				return
			default:
				switch {
				case o.getDelay() == 0: // skip wait if no delay
					// FIXME: .Call() may start more than 1 worker because for{} runs too fast
					o.newProcess()
				default:
					atomic.AddInt32(&o.blockCount, 1)
					select {
					case o.arg <- arg:
						atomic.AddInt32(&o.blockCount, -1)
						return
					case <-o.submitterTickerSignal:
						atomic.AddInt32(&o.blockCount, -1)
						// NOTE: If <delay> is too small, select may choose this case even <o.arg> isn't blocked,
						//       which may be caused by gc.  A proper value should be >= 0.1s in my test.
						life--
						if life <= 0 {
							life = atomic.LoadInt32(&o.life)
							o.newProcess()
						}
					}
				}
			}
		}
	}
}

func (o *Pool) SetMaxSignal(max uint) {
	if max > math.MaxInt32 {
		golib.Panic("max > math.MaxInt32, want <= math.MaxInt32")
	}
	atomic.StoreInt32(&o.maxTickSignal, int32(max))
}

//// Set when to create or kill a goroutine.
//// A new goroutine will be created after the argument blocked for <delay> ns.
//// A goroutine will be killed after idle for <timeout> ns
//func (o *Pool) SetTime(delay time.Duration, timeout time.Duration) {
//	o.setDelay(delay)
//	o.setTimeout(timeout)
//}
//
//// Change how many goroutines the Pool can create.
//func (o *Pool) SetCount(min uint, max uint) {
//	if min > max {
//		golib.Panic("min:%v > max:%v, want min <= max", min, max)
//	}
//
//	o.setMin(int32(min))
//	o.setMax(int32(max))
//
//	n := o.getCur() - o.getMin()
//	for i := int32(0); i < n; i++ {
//		o.newProcess()
//	}
//}
