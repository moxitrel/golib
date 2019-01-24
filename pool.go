/*

NewPool min max delay timeout bufferSize fun -> *Pool
	.Stop
	.Wait
	.Call       arg
	.setTimeout delay idle
	.setCount   min   max

* See also
- Go并发调度器解析之实现一个协程池: https://zhuanlan.zhihu.com/p/37754274

*/
package gosvc

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// Start [min, max] goroutines of <Pool.fun> to process args in <Pool.arg>
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
	// the max number of timeout-signal to sent
	maxSignal int32
	// decrease worker's life when timeout-signal received, exit when life <= 0
	defaultLife int32

	// lock cur when updating
	curLock sync.Mutex

	stopOnce sync.Once
	wg       sync.WaitGroup
}

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
func (o *Pool) getDefaultLife() int32 {
	return atomic.LoadInt32(&o.defaultLife)
}

func NewPool(min, max uint, idleTimeout time.Duration, fun func(interface{})) (o *Pool) {
	if fun == nil {
		panic(fmt.Errorf("fun == nil, want !nil"))
	}

	o = &Pool{
		//min:       	.setCount(),
		//max:       	.setCount(),
		//cur:       	0,
		//freeCount: 	0,

		fun: fun,
		arg: make(chan interface{}, 1<<6),

		//idle:  		.setTimeout(),
		//idleTicker: 	.setTimeout(),
		maxSignal:   100000,
		defaultLife: 2,

		//curLock:  sync.Mutex{},
		//stopOnce: sync.Once{},
		//wg:       sync.WaitGroup{},
	}
	// init min, max
	o.setCount(min, max)
	// init idle, idleTicker
	o.setTimeout(idleTimeout)
	return
}

// Change how many goroutines Pool can create.
func (o *Pool) setCount(min, max uint) {
	if max > math.MaxInt32 {
		panic(fmt.Errorf("max:%v > math.MaxInt32, want <= math.MaxInt32", max))
	}
	if min > max {
		panic(fmt.Errorf("min:%v > max:%v, want min <= max", min, max))
	}
	atomic.StoreInt32(&o.min, int32(min))
	atomic.StoreInt32(&o.max, int32(max))

	for n := o.getMin() - o.getCur(); n > 0; n-- {
		o.newWorker()
	}
}

// A goroutine will be killed after idle for <idle> ns.
// idle: idle forever if <= 0
func (o *Pool) setTimeout(idle time.Duration) {
	atomic.StoreInt64((*int64)(&o.idle), idle.Nanoseconds())

	// destroy old ticker
	if o.idleTicker != nil {
		o.idleTicker.Stop()
		o.idleTicker = nil
	}
	if o.idle > 0 {
		//o.wg.Add(1)
		ticker := time.NewTicker(o.idle)
		o.idleTicker = NewSvc(nil, func() {
			ticker.Stop()
			//o.wg.Done()
		}, func() {
			<-ticker.C

			// stop 1/2 idle workers
			n := o.getFreeCount() / 2
			if maxSignal := atomic.LoadInt32(&o.maxSignal); n > maxSignal {
				n = maxSignal
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

			if o.getCur() < o.getMax() && o.getFreeCount() < 1 {
				o.newWorker()
			}

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
		o.setTimeout(0)
		atomic.StoreInt32(&o.maxSignal, math.MaxInt32)
		atomic.StoreInt32(&o.defaultLife, 1) // stop immediately when receive a timeout-signal

		o.wg.Add(1)
		go func() {
			defer o.wg.Done()
			for o.getCur() > 0 {
				time.Sleep(100 * time.Millisecond)
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
	if o.getMax() > 0 {
		o.arg <- arg
	}
}
