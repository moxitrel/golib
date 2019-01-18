/*

NewPool min max delay timeout bufferSize fun -> *Pool
	.Stop
	.Wait
	.Call       arg
	.setTimeout delay idle
	.setCount   min   max

*/
package gosvc

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Start [min, max] goroutines of <Pool.fun> to process <Pool.arg>
type Pool struct {
	//
	// make Pool 64-bit aligned, put fields 64-bit long in first
	//

	// the current number of workers
	cur int64
	// destroy the workers idle for <idle> ns
	idle time.Duration

	// at least <min> workers will be created and live all the time
	min int32
	// the max number of workers can be created
	max int32
	// the number of idle workers
	freeCount int32

	fun func(interface{})
	arg chan interface{}

	// send timeout-signal to worker periodically
	idleTicker *Svc
	// the max number of timeout-signal to sent after a tick
	maxSignal int32
	// decrease worker's life when timeout-signal received, exit when life <= 0
	defaultLife int32

	// lock cur when updating
	curLock sync.Mutex

	stopOnce sync.Once
	wg       sync.WaitGroup
}

const _POOL_STOP_DELAY = 200 * time.Millisecond

type _PoolTimeoutSignal struct{}

var poolTimeoutSignal _PoolTimeoutSignal

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
func (o *Pool) incFreeCount() int32 {
	return atomic.AddInt32(&o.freeCount, 1)
}
func (o *Pool) decFreeCount() int32 {
	return atomic.AddInt32(&o.freeCount, -1)
}
func (o *Pool) getMaxSignal() int32 {
	return atomic.LoadInt32(&o.maxSignal)
}
func (o *Pool) getDefaultLife() int32 {
	return atomic.LoadInt32(&o.defaultLife)
}

func NewPool(min, max uint, timeout time.Duration, bufferSize uint, fun func(interface{})) (o *Pool) {
	//if bufferSize > math.MaxInt32 {
	//	panic(fmt.Sprintf("bufferSize:%v > math.MaxInt32, want <= math.MaxInt32", max))
	//}
	if fun == nil {
		panic("fun == nil, want !nil")
	}

	o = &Pool{
		//min:       	.setCount(),
		//max:       	.setCount(),
		//cur:       	0,
		//freeCount: 	0,

		fun: fun,
		arg: make(chan interface{}, bufferSize),

		//idle:  		.setTimeout(),
		//idleTicker: 	.setTimeout(),

		defaultLife: 2,
		maxSignal:   math.MaxInt16,

		//curLock:  sync.Mutex{},
		//stopOnce: sync.Once{},
		//wg:       sync.WaitGroup{},
	}
	// init min, max
	o.setCount(min, max)
	// init idle, idleTicker
	o.setTimeout(timeout)
	return
}

func PoolWrapper(fun func(interface{})) (func(interface{}), func()) {
	v := NewPool(uint(runtime.NumCPU()+1), 1<<23,  10*time.Minute, 1<<20, fun)
	return v.Call, v.Stop
}

// Create a new worker goroutine.
func (o *Pool) newWorker() {
	// NOTE: PLEASE DO NOT FIX THE FOLLOW PIECE OF CODE UNTIL ENCOUNTER ISSUES.
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
		o.decCur() >= o.getMax() { // detect the changes of o.cur between the interval
		return
	}

	o.incFreeCount()
	o.wg.Add(1)
	go func() {
		defer o.wg.Done()
		life := o.getDefaultLife()

		for ; ; o.incFreeCount() {
			arg := <-o.arg
			o.decFreeCount()

			switch arg {
			case poolTimeoutSignal:
				life--
				if life <= 0 {
					// Use Mutex will be more robust:
					//		o.curLock.Lock()
					//		switch {
					//		case o.getCur() <= o.getMin(): 	// continue to work
					//			o.curLock.Unlock()
					//			life = o.getDefaultLife()
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
						life = o.getDefaultLife()
					}
				}
			default:
				life = o.getDefaultLife()
				o.fun(arg)
			}
		}
	}()
}

func (o *Pool) Stop() {
	o.stopOnce.Do(func() {
		o.setCount(0, 0)
		atomic.StoreInt32(&o.defaultLife, 1) // stop immediately when receive a timeout-signal

		if o.idleTicker != nil {
			o.idleTicker.Stop()
			o.idleTicker = nil
		}

		o.wg.Add(1)
		go func() {
			defer o.wg.Done()
			for o.getCur() > 0 {
				time.Sleep(_POOL_STOP_DELAY)
				for freeCount := o.getFreeCount(); freeCount > 0; freeCount-- {
					select {
					case o.arg <- poolTimeoutSignal:
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


func (o *Pool) Call(arg interface{}) {
CALL:
	switch {
	case o.getMax() <= 0: // stopped?
		// nop
	case o.getCur() >= o.getMax():
		o.arg <- arg
	default:
		select {
		case o.arg <- arg:
		default:
			// NOTE: Call() may create more than 1 worker
			o.newWorker()
			goto CALL
		}
	}
}

// Change when to create or kill a goroutine.
// A goroutine will be killed after idle for <idle> ns.
// idle: idle forever if <= 0
func (o *Pool) setTimeout(idle time.Duration) {
	atomic.StoreInt64((*int64)(&o.idle), idle.Nanoseconds())

	// clean old ticker
	if o.idleTicker != nil {
		o.idleTicker.Stop()
		o.idleTicker = nil
	}
	if o.idle > 0 {
		ticker := time.NewTicker(o.idle)
		o.wg.Add(1)
		o.idleTicker = NewSvc(nil, func() {
			o.wg.Done()
			ticker.Stop()
		}, func() {
			<-ticker.C

			// stop idle workers
			n := o.getFreeCount()
			if n > o.getMaxSignal() {
				n = o.getMaxSignal()
			}
			for ; n > 0; n-- {
				select {
				case o.arg <- poolTimeoutSignal:
				default:
					return
				}
			}
		})
	}
}

// Change how many goroutines Pool can create.
func (o *Pool) setCount(min uint, max uint) {
	if min > math.MaxInt32 {
		panic(fmt.Sprintf("min:%v > math.MaxInt32, want <= math.MaxInt32", min))
	}
	if max > math.MaxInt32 {
		panic(fmt.Sprintf("max:%v > math.MaxInt32, want <= math.MaxInt32", max))
	}
	if min > max {
		panic(fmt.Sprintf("min:%v > max:%v, want min <= max", min, max))
		//min = max
	}
	atomic.StoreInt32(&o.min, int32(min))
	atomic.StoreInt32(&o.max, int32(max))
	// if idle is too small, new process may exit quicker than creating.
	for n := int64(min) - o.getCur(); n > 0; n-- {
		o.newWorker()
	}
}
