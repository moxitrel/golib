/*

func         PoolWrap (func (interface{})                        ) *Pool
func (*Pool) WithCount(min   uint         , max     uint         ) *Pool
func (*Pool) WithTime (delay time.Duration, timeout time.Duration) *Pool

func (*Pool) Call    (interface{})

*/
package svc

import (
	"github.com/moxitrel/golib"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

const (
	_POOL_MIN     = 2
	_POOL_MAX     = math.MaxUint16
	_POOL_DELAY   = 100 * time.Millisecond
	_POOL_TIMEOUT = time.Minute
)

// Start [min, max] goroutines of <Pool.fun> to process <Pool.arg>
//
// * Example
// f := func(x interface{}) { time.Sleep(time.Second) }
// p := PoolWrap(f)	    // start 1 goroutines of f
// p.Call("1")			// run f("1") in background and return immediately
// p.Call("2")			// run f("2") in background after block <Pool.delay> ns
// p.Call("3")			// run f("3") in background after block <Pool.delay> ns
//
type Pool struct {
	// the current number of goroutines
	// put in head to make <cur> 64-bit aligned
	cur int64
	// at least <min> goroutines will be created and live all the time
	min uint32
	// the max number of goroutines can be created
	max uint32
	// create a new goroutine if <arg> is blocked for <delay> ns
	// a proper value should be >= 0.1s
	delay time.Duration
	// destroy the goroutines which idle for <timeout> ns
	timeout time.Duration

	fun func(interface{})
	arg chan interface{}

	// protect cur when update
	curLock sync.Mutex
}

func (o *Pool) getCur() int64 {
	return atomic.LoadInt64(&o.cur)
}
func (o *Pool) getMin() int64 {
	return int64(atomic.LoadUint32(&o.min))
}
func (o *Pool) getMax() int64 {
	return int64(atomic.LoadUint32(&o.max))
}

func NewPool(min, max uint, delay, timeout time.Duration, fun func(interface{})) (o *Pool) {
	if delay < 0 {
		golib.Panic("delay < 0, want >= 0")
	}
	if timeout < 0 {
		golib.Panic("timeout < 0, want >= 0")
	}
	if fun == nil {
		golib.Panic("fun == nil, want !nil")
	}

	o = &Pool{
		cur:     0,
		min:     uint32(min),
		max:     uint32(max),
		delay:   delay,
		timeout: timeout,
		fun:     fun,
		arg:     make(chan interface{}),
		curLock: sync.Mutex{},
	}
	for o.cur < int64(o.min) {
		o.newProcess()
	}
	return
}

func PoolWrap(fun func(interface{})) (func(interface{}) /* stop */, func()) {
	v := NewPool(_POOL_MIN, _POOL_MAX, _POOL_DELAY, _POOL_TIMEOUT, fun)
	return v.GetCall(), v.Stop
}

// The created goroutine won't quit unless time out. Set min to 0 if want to quit all.
func (o *Pool) newProcess() {
	o.curLock.Lock()
	// NOTE: use atomic for o.cur to avoid data race from .GetCall()
	if o.getCur() >= o.getMax() {
		o.curLock.Unlock()
		return
	}
	atomic.AddInt64(&o.cur, 1)
	o.curLock.Unlock()

	// init timer
	timeoutTimer := time.NewTimer(o.timeout)
	if !timeoutTimer.Stop() {
		<-timeoutTimer.C
	}

	var loop *Loop
	loop = NewLoop(func() {
		select {
		case arg := <-o.arg:
			// the outer select: skip creating timer when busy
			o.fun(arg)
		default:
			// start timer
			timeoutTimer.Reset(o.delay)

			select {
			case arg := <-o.arg:
				// reset timer
				if !timeoutTimer.Stop() {
					<-timeoutTimer.C
				}
				o.fun(arg)
			case <-timeoutTimer.C:
				// quit if idle for <timeout> ns
				o.curLock.Lock()
				if o.getCur() > o.getMin() {
					atomic.AddInt64(&o.cur, -1)
					loop.Stop()
				}
				o.curLock.Unlock()
			}
		}
	})
}

func (o *Pool) Stop() {
	atomic.StoreUint32(&o.min, 0)
	atomic.StoreUint32(&o.max, 0)
	//atomic.StoreUint32(&o.timeout, 0)
}

func (o *Pool) GetCall() func(interface{}) {
	// init timer
	delayTimer := time.NewTimer(o.delay)
	if !delayTimer.Stop() {
		<-delayTimer.C
	}

	return func(arg interface{}) {
		for {
			// NOTE: use atomic instead of lock for o.cur for performance
			if o.getCur() >= o.getMax() {
				// skip creating timer, when busy and no more goroutine can be created
				o.arg <- arg
			} else {
				delayTimer.Reset(o.delay) // start timer
				select {
				case o.arg <- arg:
					// stop and clean timer
					if !delayTimer.Stop() {
						<-delayTimer.C
					}
				case <-delayTimer.C:
					// If <delay> is too small, select may choose this case even <o.arg> isn't blocked.
					o.newProcess()
					continue
				}
			}
			break
		}
	}
}

// Set when to create or kill a goroutine.
// A new goroutine will be created after the argument blocked for ^delay ns.
// A goroutine will be killed after idle for ^timeout ns
//func (o *Pool) WithTime(delay time.Duration, timeout time.Duration) *Pool {
//	o.delay = delay
//	o.timeout = timeout
//	return o
//}

// Change how many goroutines the Pool can create, ^min <= count <= ^max.
//func (o *Pool) WithCount(min uint, max uint) *Pool {
//	if min > max {
//		golib.Panic("min:%v > max:%v !", min, max)
//	}
//
//	atomic.StoreUint32(&o.min, uint32(min))
//	atomic.StoreUint32(&o.max, uint32(max))
//
//	for o.getCur() < int64(o.min) {
//		o.newProcess()
//	}
//
//	return o
//}
