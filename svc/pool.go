package svc

import (
	"sync/atomic"
	"time"
)

// Start [min, max] goroutines of fun to process arg
//
// [Example]
// f := func(x interface{}) { time.Sleep(time.Second) }
// p := NewPool(2, f)	// start 2 goroutines of f
// p.Call("1")			// run f("1") in background
// p.Call("2")			// run f("2") in background
// p.Call("3")			// run f("3") in background after wait POOL_DELAY ns
//
type Pool struct {
	fun func(interface{})
	arg chan interface{}

	// at least <min> coroutines will be created and live all the time
	min uint16
	// the current number of coroutines
	cur int32
	// the max number of coroutines can be created
	max uint16
	// create a new coroutine if arg is blocked for <delay> ns
	delay time.Duration
	// destroy the coroutine that is idle for <timeout> ns
	timeout time.Duration
}

func NewPool(min uint, fun func(interface{})) (v *Pool) {
	v = &Pool{
		fun:     fun,
		arg:     make(chan interface{}),
		min:     uint16(min),
		cur:     0,
		max:     POOL_MAX,
		delay:   POOL_DELAY, // a proper value should at least 0.1s
		timeout: POOL_TIMEOUT,
	}
	for v.cur < int32(v.min) {
		v.newProcess()
	}
	return
}

func (o *Pool) SetTime(delay time.Duration, timeout time.Duration) {
	o.delay = delay
	o.timeout = timeout
}

func (o *Pool) SetCount(min uint16, max uint16) {
	o.min = min
	o.max = max
	for o.cur < int32(o.min) {
		o.newProcess()
	}
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
				// NOTE: use delay to ensure the new process is already started when try again
				o.Call(arg)
			} else {
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

	var loop *LoopService
	loop = NewLoopService(func() {
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
