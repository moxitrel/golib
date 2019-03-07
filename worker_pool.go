/*

NewWorkerPool min max delay timeout bufferSize fun -> *WorkerPool
	.Stop
	.Wait
	.Submit     arg
	.setTimeout delay idle
	.setCount   min   max

* See Also
- http://blog.taohuawu.club/article/42

*/
package gosvc

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// Start [min, max] goroutines of <WorkerPool.fun> to process args in <WorkerPool.arg>
type WorkerPool struct {
	//
	// make WorkerPool 64-bit aligned (put fields 64-bit in first)
	//

	// the current number of workers
	cur int64
	// destroy the workers idle for <idle> ns
	idle time.Duration

	// at least <min> workers will be created and live all the time
	min uint32
	// the max number of workers can be created
	max uint32
	// the current number of idle workers
	freeCount uint32

	fun func(interface{})
	arg chan interface{}

	// send timeout-signal to worker periodically
	idleTicker *Svc
	// the max number of timeout-signal to sent per tick
	maxSignal uint32
	// decrease worker's life when timeout-signal received, exit when life <= 0
	defaultLife uint32

	// lock cur when updating
	curLock sync.Mutex

	stopOnce sync.Once
	wg       sync.WaitGroup
}

type _WorkerPoolTimeoutSignal struct{}

func (o *WorkerPool) getCur() int64 {
	return atomic.LoadInt64(&o.cur)
}
func (o *WorkerPool) incCur() int64 {
	return atomic.AddInt64(&o.cur, 1)
}
func (o *WorkerPool) decCur() int64 {
	return atomic.AddInt64(&o.cur, -1)
}
func (o *WorkerPool) getMin() int64 {
	return int64(atomic.LoadUint32(&o.min))
}
func (o *WorkerPool) getMax() int64 {
	return int64(atomic.LoadUint32(&o.max))
}
func (o *WorkerPool) getFreeCount() uint32 {
	return atomic.LoadUint32(&o.freeCount)
}
func (o *WorkerPool) incFreeCount() uint32 {
	return atomic.AddUint32(&o.freeCount, 1)
}
func (o *WorkerPool) decFreeCount() uint32 {
	return atomic.AddUint32(&o.freeCount, ^uint32(0))
}
func (o *WorkerPool) getDefaultLife() uint32 {
	return atomic.LoadUint32(&o.defaultLife)
}
func (o *WorkerPool) getMaxSignal() uint32 {
	return atomic.LoadUint32(&o.maxSignal)
}

func NewWorkerPool(min, max uint, idleTimeout time.Duration, fun func(interface{})) (o *WorkerPool) {
	if fun == nil {
		panic(fmt.Errorf("fun == nil, want !nil"))
	}

	o = &WorkerPool{
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

// Change how many goroutines WorkerPool can create.
func (o *WorkerPool) setCount(min, max uint) {
	if max > math.MaxUint32 {
		panic(fmt.Errorf("max:%v > math.MaxInt32, want <= math.MaxUint32", max))
	}
	if min > max {
		panic(fmt.Errorf("min:%v > max:%v, want min <= max", min, max))
	}
	atomic.StoreUint32(&o.min, uint32(min))
	atomic.StoreUint32(&o.max, uint32(max))

	for n := int64(min) - o.getCur(); n > 0; n-- {
		o.newWorker()
	}
}

// A goroutine will be killed after idle for <idle> ns.
// idle: idle forever if == 0
func (o *WorkerPool) setTimeout(idle time.Duration) {
	if idle < 0 {
		panic(fmt.Errorf("idle:%v < 0, want >= 0", idle))
	}
	atomic.StoreInt64((*int64)(&o.idle), idle.Nanoseconds())

	// destroy old ticker
	if o.idleTicker != nil {
		o.idleTicker.Stop()
		o.idleTicker = nil
	}
	if idle > 0 {
		//o.wg.Add(1)
		ticker := time.NewTicker(idle)
		o.idleTicker = NewSvc(nil, func() {
			ticker.Stop()
			//o.wg.Done()
		}, func() {
			<-ticker.C

			// stop 1/2 idle workers
			n := o.getFreeCount() / 2
			if n > o.getMaxSignal() {
				n = o.getMaxSignal()
			}
			for ; n > 0; n-- {
				select {
				case o.arg <- _WorkerPoolTimeoutSignal{}:
				default:
					return
				}
			}
		})
	}
}

// Create a new worker goroutine.
func (o *WorkerPool) newWorker() {
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
			case _WorkerPoolTimeoutSignal{}:
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

func (o *WorkerPool) Stop() {
	o.stopOnce.Do(func() {
		o.setCount(0, 0)
		o.setTimeout(0)	// close idleTicker
		atomic.StoreUint32(&o.defaultLife, 1)	// stop immediately when receive a timeout-signal

		o.wg.Add(1)
		go func() {
			defer o.wg.Done()
			for o.getCur() > 0 {
				time.Sleep(100 * time.Millisecond)
				for freeCount := o.getFreeCount(); freeCount > 0; freeCount-- {
					select {
					case o.arg <- _WorkerPoolTimeoutSignal{}:
					default:
						goto STOP_SEND
					}
				}
			STOP_SEND:
			}
		}()
	})
}

func (o *WorkerPool) Wait() {
	o.wg.Wait()
}

func (o *WorkerPool) Submit(arg interface{}) {
	if o.getMax() > 0 {
		o.arg <- arg
	}
}
