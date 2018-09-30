/*

func         NewPool (func (interface{})                        ) *Pool
func (*Pool) SetCount(min   uint         , max     uint         ) *Pool
func (*Pool) SetTime (delay time.Duration, timeout time.Duration) *Pool

func (*Pool) Call    (interface{})

*/
package svc

import (
	golib ".."
	"fmt"
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
	POOL_MIN     = 0
	POOL_MAX     = math.MaxUint16
	POOL_DELAY   = 200 * time.Millisecond
	POOL_TIMEOUT = time.Minute
)

type Pool struct {
	fun func(interface{})
	arg chan interface{}

	// at least <min> coroutines will be created and live all the time
	min uint32
	// the current number of coroutines
	cur int64
	// the max number of coroutines can be created
	max uint32
	// create a new coroutine if <arg> is blocked for <delay> ns
	// a proper value should be >= 0.1s
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
		delay:   POOL_DELAY,
		timeout: POOL_TIMEOUT,
	}
	for v.cur < int64(v.min) {
		v.newProcess()
	}
	return
}

// Set when to create or kill a goroutine.
// A new goroutine will be created after the argument blocked for ^delay ns.
// A goroutine will be killed after idle for ^timeout ns
func (o *Pool) SetTime(delay time.Duration, timeout time.Duration) *Pool {
	o.delay = delay
	o.timeout = timeout
	return o
}

// Change how many goroutines the Pool can create, ^min <= count <= ^max.
func (o *Pool) SetCount(min uint, max uint) *Pool {
	if min > max {
		golib.Warn(fmt.Sprintf("min:%v > max:%v !", min, max))
		min = max
	}

	o.min = uint32(min)
	o.max = uint32(max)
	for o.cur < int64(o.min) {
		o.newProcess()
	}
	return o
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
	if atomic.AddInt64(&o.cur, 1) > int64(o.max) {
		// no coroutine created, restore the value
		atomic.AddInt64(&o.cur, -1)
		return false
	}

	var loop *Loop
	loop = NewLoop(func() {
		select {
		case arg := <-o.arg:
			o.fun(arg)
		case <-time.After(o.timeout): // quit if idle for <timeout> ns
			if atomic.AddInt64(&o.cur, -1) >= int64(o.min) {
				loop.Stop()
			} else {
				// coroutine isn't killed, restore the value
				atomic.AddInt64(&o.cur, 1)
			}
		}
	})

	return true
}
