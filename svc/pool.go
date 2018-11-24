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
	_POOL_WORKER_MIN     = 1
	_POOL_WORKER_MAX     = math.MaxUint16
	_POOL_WORKER_DELAY   = 250 * time.Millisecond
	_POOL_WORKER_TIMEOUT = 45 * time.Second

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
	timeout time.Duration

	fun func(interface{})
	arg chan interface{}

	// send timeout-signal to worker periodically
	timeoutTicker *Svc
	// How many signals to sent out after a tick.
	// Control how many idle workers to close after a timeout.
	workerDelta int32
	// decrease worker's life when tick-signal received, timeout when life <= 0
	life int32
	// send timeout-signal to submitter periodically through delayChannel
	delayTicker *Svc

	// lock cur when updating
	curLock  sync.Mutex
	stopOnce sync.Once
	wg       sync.WaitGroup
}

type _TimeoutSignal struct{}

var timeoutSignal _TimeoutSignal

func (o *Pool) getMin() int32 {
	return atomic.LoadInt32(&o.min)
}
func (o *Pool) getMax() int32 {
	return atomic.LoadInt32(&o.max)
}
func (o *Pool) getCur() int32 {
	return atomic.LoadInt32(&o.cur)
}
func (o *Pool) getFreeCount() int32 {
	return atomic.LoadInt32(&o.freeCount)
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
		golib.Panic("delay:%v <= 0, want > 0", delay)
	}
	if timeout <= 0 {
		golib.Panic("timeout:%v <= 0, want > 0", timeout)
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
		timeout:   timeout,

		fun: fun,
		arg: make(chan interface{}, bufferSize),

		timeoutTicker: nil, // inited below
		workerDelta:   512, // default number of timeoutSignal to send per tick
		life:          2,   // default life
		delayTicker:   nil, // inited below

		curLock:  sync.Mutex{},
		stopOnce: sync.Once{},
		wg:       sync.WaitGroup{},
	}
	o.SetTimeout(delay, timeout)

	// if timeout is too small, new process may exit quicker than creating.
	for i := o.min; i > 0; i-- {
		o.newProcess()
	}
	return
}

func PoolWrapper(fun func(interface{})) (func(interface{}), func()) {
	v := NewPool(_POOL_WORKER_MIN, _POOL_WORKER_MAX, _POOL_WORKER_DELAY, _POOL_WORKER_TIMEOUT, 0, fun)
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
		life := o.getLife()
		for {
			atomic.AddInt32(&o.freeCount, 1)
			arg := <-o.arg
			atomic.AddInt32(&o.freeCount, -1)
			switch arg {
			case timeoutSignal:
				life--
				if life <= 0 {
					o.curLock.Lock()
					switch {
					case o.getCur() <= atomic.LoadInt32(&o.min): // continue to work
						o.curLock.Unlock()
						life = o.getLife()
					default: // quit
						atomic.AddInt32(&o.cur, -1)
						o.curLock.Unlock()
						return
					}
				}
			default:
				life = o.getLife()
				o.fun(arg)
			}
		}
	}()
}

func (o *Pool) Stop() {
	o.stopOnce.Do(func() {
		o.SetWorkerLimit(0, 0, uint(o.workerDelta))
		atomic.StoreInt32(&o.life, 1) // stop when receive a timeoutSignal

		if o.timeoutTicker != nil {
			o.timeoutTicker.Stop()
		}
		if o.delayTicker != nil {
			o.delayTicker.Stop()
		}

		o.wg.Add(1)
		go func() {
			defer o.wg.Done()
			for o.getCur() > 0 {
				time.Sleep(_STOP_DELAY)
				for freeCount := o.getFreeCount(); freeCount > 0; freeCount-- {
					select {
					case o.arg <- timeoutSignal:
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

// FIXME: Call() creates more than 1 worker if delay = 0, add pause after newProcess()
func (o *Pool) Call(arg interface{}) {
CALL:
	switch {
	case o.getMax() <= 0:
		// return
	case o.getDelay() > 0:
		o.arg <- arg
	case o.getCur() >= o.getMax():
		o.arg <- arg
	default:
		select {
		case o.arg <- arg:
			if len(o.arg) > 0 {
				o.newProcess()
			}
		default:
			o.newProcess()
			goto CALL
		}
	}
}

// Set when to create or kill a goroutine.
// A new goroutine will be created after the argument blocked for <delay> ns.
// A goroutine will be killed after idle for <timeout> ns
func (o *Pool) SetTimeout(delay time.Duration, timeout time.Duration) {
	if delay < 0 {
		golib.Panic("delay:%v <= 0, want > 0", delay)
	}
	if timeout <= 0 {
		golib.Panic("timeout:%v <= 0, want > 0", timeout)
	}
	atomic.StoreInt64((*int64)(&o.delay), delay.Nanoseconds())
	atomic.StoreInt64((*int64)(&o.timeout), timeout.Nanoseconds())

	// init timeout-ticker
	if o.timeoutTicker != nil {
		o.timeoutTicker.Stop()
	}
	if o.timeout < math.MaxInt64 {
		timeoutTicker := time.NewTicker(o.timeout)
		o.wg.Add(1)
		o.timeoutTicker = NewSvc(nil, func() {
			<-timeoutTicker.C
			// stop idle workers
			n := o.getFreeCount()
			maxTimeoutSignal := atomic.LoadInt32(&o.workerDelta)
			if n > maxTimeoutSignal {
				n = maxTimeoutSignal
			}
			for ; n > 0; n-- {
				select {
				case o.arg <- timeoutSignal:
				default:
					return
				}
			}
		}, func() {
			o.wg.Done()
			timeoutTicker.Stop()
		})
	}

	// init delay-ticker
	if o.delayTicker != nil {
		o.delayTicker.Stop()
	}
	if o.delay > 0 {
		delayTicker := time.NewTicker(o.delay)
		o.wg.Add(1)
		o.delayTicker = NewSvc(nil, func() {
			<-delayTicker.C
			jobs := len(o.arg)
			idleWorkers := o.getFreeCount()
			avaliableWorkers := o.getMax() - o.getCur()
			if n := int32(jobs) - idleWorkers; avaliableWorkers > 0 && n > 0 {
				// TODO: compute a proper number
				if n > avaliableWorkers {
					n = avaliableWorkers
				}
				maxTimeoutSignal := atomic.LoadInt32(&o.workerDelta)
				if n > maxTimeoutSignal {
					n = maxTimeoutSignal
				}
				for ; n > 0; n-- {
					o.newProcess()
				}
			}
		}, func() {
			o.wg.Done()
			delayTicker.Stop()
		})
	}
}

// Change how many goroutines the Pool can create.
func (o *Pool) SetWorkerLimit(min uint, max uint, delta uint) {
	if min > math.MaxInt32 {
		golib.Panic("min:%v > math.MaxInt32, want <= math.MaxInt32", min)
	}
	if max > math.MaxInt32 {
		golib.Panic("max:%v > math.MaxInt32, want <= math.MaxInt32", max)
	}
	if min > max {
		golib.Panic("min:%v > max:%v, want min <= max", min, max)
	}
	if delta > math.MaxInt32 {
		golib.Panic("delta:%v > math.MaxInt32, want <= math.MaxInt32", delta)
	}
	atomic.StoreInt32(&o.min, int32(min))
	atomic.StoreInt32(&o.max, int32(max))
	atomic.StoreInt32(&o.workerDelta, int32(max))
	for n := o.getMin() - o.getCur(); n > 0; n-- {
		o.newProcess()
	}
}
