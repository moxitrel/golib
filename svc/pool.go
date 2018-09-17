package svc

import (
	"github.com/moxitrel/golib"
	"math"
	"sync/atomic"
	"time"
)

// Start [min, max] goroutines of <Pool.fun> to process <Pool.arg>
//
// * Example
// f := func(x interface{}) { time.Sleep(time.Second) }
// p := NewPool(f)	    // start 2 goroutines of f
// p.Call("1")			// run f("1") in background and return immediately
// p.Call("2")			// run f("2") in background and return immediately
// p.Call("3")			// run f("3") in background after block <Pool.delay> ns
//

const (
	POOL_MIN     = 2
	POOL_MAX     = math.MaxUint16
	POOL_DELAY   = 200 * time.Millisecond
	POOL_TIMEOUT = time.Minute
)

type Pool struct {
	fun func(interface{})
	arg chan interface{}

	// at least <min> coroutines will be created and live all the time
	min uint16
	// the current number of coroutines
	cur int32
	// the max number of coroutines can be created
	max uint16
	// create a new coroutine if <arg> is blocked for <delay> ns
	delay time.Duration
	// destroy the coroutine idle for <timeout> ns
	timeout time.Duration
}

func NewPool(fun func(interface{})) (v *Pool) {
	if fun == nil {
		golib.Panic("^fun shouldn't be nil!\n")
	}
	v = &Pool{
		fun:     fun,
		arg:     make(chan interface{}),
		min:     POOL_MIN,
		cur:     0,
		max:     POOL_MAX,
		delay:   POOL_DELAY, // a proper value should be at least 0.1s
		timeout: POOL_TIMEOUT,
	}
	for v.cur < int32(v.min) {
		v.newProcess()
	}
	return
}

func (o *Pool) SetCount(min uint, max uint) {
	if min > max {
		golib.Warn("min:%v > max:%v!\n", min, max)
		min = max
	}
	o.min = uint16(min)
	o.max = uint16(max)
	for o.cur < int32(o.min) {
		o.newProcess()
	}
}

func (o *Pool) SetTime(delay time.Duration, timeout time.Duration) {
	o.delay = delay
	o.timeout = timeout
}

func (o *Pool) Call(arg interface{}) {
	select {
	case o.arg <- arg:
		// ensure <o.arg> is blocked
	default:
		select {
		case o.arg <- arg:
		case <-time.After(o.delay):
			// If <delay> is too small, select may choose this case even <o.arg> isn't blocked.
			if o.newProcess() {
				// NOTE: expect the new process has started by delay when try again
				o.Call(arg)
			} else {
				// wait if no more coroutine can be created
				o.arg <- arg
			}
		}
	}
}

// The created coroutine won't quit unless time out. Set min to 0 if want to quit all.
func (o *Pool) newProcess() bool {
	if atomic.AddInt32(&o.cur, 1) > int32(o.max) {
		// no coroutine created, restore the value
		atomic.AddInt32(&o.cur, -1)
		return false
	}

	var loop *Loop
	loop = NewLoop(func() {
		select {
		case arg := <-o.arg:
			o.fun(arg)
		case <-time.After(o.timeout): // quit if idle for <timeout> ns
			if atomic.AddInt32(&o.cur, -1) >= int32(o.min) {
				loop.Stop()
			} else {
				// coroutine isn't killed, restore the value
				atomic.AddInt32(&o.cur, 1)
			}
		}
	})

	return true
}
