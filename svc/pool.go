package svc

import (
	"github.com/moxitrel/golib"
	"math"
	"runtime"
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
		Timer: time.NewTimer(time.Minute),
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
	_POOL_WORKER_DELAY = 200 * time.Millisecond

	// time to wait for receiving sent args when receive stop signal
	_STOP_DELAY = 200 * time.Millisecond
)

// Start [min, max] goroutines of <Pool.fun> to process <Pool.arg>
type Pool struct {
	// the current number of workers
	// put cur in first to make it 64-bit aligned
	cur int64
	// at least <min> workers will be created and live all the time
	min int32
	// the max number of workers can be created
	max int32
	// the number of idle workers
	freeCount int32
	// how many Call() blocked when delay > 0
	blockCount int64
	// create a new worker if <arg> is blocked for <delay> ns, a proper value should be >= 0.1s
	delay time.Duration
	// destroy the workers which idle for <idle> ns
	idle time.Duration

	fun func(interface{})
	arg chan interface{}

	// send timeoutSignal to worker periodically
	idleTicker *Svc
	// send timeoutSignal to submitter periodically through delayChannel
	delayTicker *Svc
	// decrease worker's life when tick-signal received, idle when life <= 0
	life int32
	// How many signals to sent out after a tick.
	// Control how many workers to stop after idle or create after delay.
	maxSignal int32

	// lock cur when updating
	curLock sync.Mutex

	stopOnce sync.Once
	wg       sync.WaitGroup
}

type _TimeoutSignal struct{}

var timeoutSignal _TimeoutSignal

func (o *Pool) getCur() int64 {
	return atomic.LoadInt64(&o.cur)
}
func (o *Pool) incCur() int64 {
	return atomic.AddInt64(&o.cur, 1)
}
func (o *Pool) decCur() int64 {
	return atomic.AddInt64(&o.cur, -1)
}
func (o *Pool) getMin() int64 {
	return int64(atomic.LoadInt32(&o.min))
}
func (o *Pool) getMax() int64 {
	return int64(atomic.LoadInt32(&o.max))
}
func (o *Pool) getFreeCount() int32 {
	return atomic.LoadInt32(&o.freeCount)
}
func (o *Pool) incFreeCount() int64 {
	return int64(atomic.AddInt32(&o.freeCount, 1))
}
func (o *Pool) decFreeCount() int64 {
	return int64(atomic.AddInt32(&o.freeCount, -1))
}
func (o *Pool) getBlockCount() int64 {
	return atomic.LoadInt64(&o.blockCount)
}
func (o *Pool) incBlockCount() int64 {
	return int64(atomic.AddInt64(&o.blockCount, 1))
}
func (o *Pool) decBlockCount() int64 {
	return int64(atomic.AddInt64(&o.blockCount, -1))
}
func (o *Pool) getDelay() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&o.delay)))
}
func (o *Pool) getLife() int32 {
	return atomic.LoadInt32(&o.life)
}
func (o *Pool) getMaxSignal() int32 {
	return atomic.LoadInt32(&o.maxSignal)
}

func NewPool(min, max uint, delay, timeout time.Duration, bufferSize uint, fun func(interface{})) (o *Pool) {
	if bufferSize > math.MaxInt32 {
		golib.Panic("bufferSize:%v > math.MaxInt32, want <= math.MaxInt32", max)
	}
	if fun == nil {
		golib.Panic("fun == nil, want !nil")
	}

	o = &Pool{
		//cur:       0,
		//freeCount: 0,

		//min:       int32(min),
		//max:       int32(max),

		fun: fun,
		arg: make(chan interface{}, bufferSize),

		//delay:     	delay,
		//idle:   		idle,
		//idleTicker: 	nil,
		//delayTicker:	nil,
		maxSignal: math.MaxInt16, // default max number of signals to send per tick
		life:      2,             // default life

		//curLock:  sync.Mutex{},
		//stopOnce: sync.Once{},
		//wg:       sync.WaitGroup{},
	}
	// init min, max
	o.SetCount(min, max)
	// init delay, idle, delayTicker, idleTicker
	o.SetTimeout(delay, timeout)
	return
}

func PoolWrapper(fun func(interface{})) (func(interface{}), func()) {
	v := NewPool(uint(runtime.NumCPU()+1), 1<<23, 0, 2*time.Minute, 1<<19, fun)
	return v.Call, v.Stop
}

// Create a new goroutine.
func (o *Pool) newProcess() {
	// NOTE: PLEASE DO NOT FIX THE FOLLOW PIECE OF CODE UNTIL ENCOUNTER ISSUE.
	//
	// The following method (in testing) to update o.cur is not correct in theory, but ok in practice for int64.
	// (int128 is better, but int32 isn't good enough)
	//
	// Use Mutex will be more robust:
	//		o.curLock.Lock()
	//		if o.getCur() < o.getMax() {
	//		 	o.incCur()
	//			o.curLock.Unlock()
	//		} else {
	//			o.curLock.Unlock()
	//			return
	//		}
	//
	if o.incCur() > o.getMax() &&
		o.decCur() >= o.getMax() { // detect the changes of o.cur in the gap
		return
	}

	o.incFreeCount()
	o.wg.Add(1)
	go func() {
		defer o.wg.Done()
		life := o.getLife()

		for ; ; o.incFreeCount() {
			arg := <-o.arg
			o.decFreeCount()

			switch arg {
			case timeoutSignal:
				life--
				if life <= 0 {
					// Use Mutex will be more robust:
					//		o.curLock.Lock()
					//		switch {
					//		case o.getCur() <= o.getMin(): 	// continue to work
					//			o.curLock.Unlock()
					//			life = o.getLife()
					//		default: 						// quit
					//			o.decCur()
					//			o.curLock.Unlock()
					//			return
					//		}
					//
					switch {
					case o.decCur() >= o.getMin():
						return
					case o.incCur() > o.getMin():
						return
					default:
						life = o.getLife()
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
		o.SetCount(0, 0)
		atomic.StoreInt32(&o.life, 1) // stop immediately when receive a timeout-signal

		if o.idleTicker != nil {
			o.idleTicker.Stop()
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

func (o *Pool) Wait() {
	o.wg.Wait()
}

// NOTE: Call() may create more than 1 worker if delay = 0
func (o *Pool) Call(arg interface{}) {
CALL:
	switch {
	case o.getMax() <= 0: // stopped?
		// nop
	case o.getDelay() > 0:
		o.incBlockCount()
		o.arg <- arg
		o.decBlockCount()
	case o.getCur() >= o.getMax():
		o.arg <- arg
	default:
		select {
		case o.arg <- arg:
			if o.getFreeCount() < int32(len(o.arg)) {
				o.newProcess()
			}
		default:
			o.newProcess()
			goto CALL
		}
	}
}

// Changue when to create or kill a goroutine.
// A new goroutine will be created after the argument blocked for <delay> ns.
// A goroutine will be killed after idle for <idle> ns
func (o *Pool) SetTimeout(delay time.Duration, idle time.Duration) {
	if delay < 0 {
		golib.Panic("delay:%v < 0, want >= 0", delay)
	}
	if idle <= 0 {
		golib.Panic("idle:%v <= 0, want > 0", idle)
	}
	atomic.StoreInt64((*int64)(&o.delay), delay.Nanoseconds())
	atomic.StoreInt64((*int64)(&o.idle), idle.Nanoseconds())

	// init idle-ticker
	if o.idleTicker != nil {
		o.idleTicker.Stop()
	}
	if o.idle < math.MaxInt64 {
		timeoutTicker := time.NewTicker(o.idle)
		o.wg.Add(1)
		o.idleTicker = NewSvc(nil, func() {
			o.wg.Done()
			timeoutTicker.Stop()
		}, func() {
			<-timeoutTicker.C

			// stop idle workers
			n := o.getFreeCount()
			if maxSignal := o.getMaxSignal(); n > maxSignal {
				n = maxSignal
			}
			for ; n > 0; n-- {
				select {
				case o.arg <- timeoutSignal:
				default:
					return
				}
			}
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
			o.wg.Done()
			delayTicker.Stop()
		}, func() {
			<-delayTicker.C

			avaliableWorkers := o.getMax() - o.getCur()
			if avaliableWorkers <= 0 {
				return
			}
			jobs := o.getBlockCount() + int64(len(o.arg)) - int64(o.getFreeCount())
			if jobs <= 0 {
				return
			}
			if jobs > avaliableWorkers {
				jobs = avaliableWorkers
			}
			if maxSignal := int64(o.getMaxSignal()); jobs > maxSignal {
				jobs = maxSignal
			}
			for ; jobs > 0; jobs-- {
				o.newProcess()
			}
		})
	}
}

// Change how many goroutines Pool can create.
func (o *Pool) SetCount(min uint, max uint) {
	if min > math.MaxInt32 {
		golib.Panic("min:%v > math.MaxInt32, want <= math.MaxInt32", min)
	}
	if max > math.MaxInt32 {
		golib.Panic("max:%v > math.MaxInt32, want <= math.MaxInt32", max)
	}
	if min > max {
		golib.Warn("min:%v > max:%v, want min <= max", min, max)
		min = max
	}
	atomic.StoreInt32(&o.min, int32(min))
	atomic.StoreInt32(&o.max, int32(max))
	// if idle is too small, new process may exit quicker than creating.
	for n := int64(min) - o.getCur(); n > 0; n-- {
		o.newProcess()
	}
}
