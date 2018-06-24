package svc

import (
	"sync/atomic"
	"time"
)

type Pool struct {
	fun func(interface{})
	arg chan interface{}

	// at least <min> coroutines will be created and live all the time
	min uint32
	// the current number of coroutines
	cur uint32
	// the max number of coroutines can be created
	max uint32
	// create a new coroutine if arg is blocked for <delay> ns
	delay time.Duration
	// destroy the coroutine if it's idle for <timeout> ns
	timeout time.Duration
}

func NewPool(fun func(interface{})) (v *Pool) {
	v = &Pool{
		fun:     fun,
		arg:     make(chan interface{}),
		min:     PoolMin,
		cur:     0,
		max:     PoolMax,
		delay:   PoolDelay,
		timeout: PoolTimeOut,
	}
	for i := uint32(0); i < v.min; i++ {
		v.newProcess()
	}
	return
}

func (o *Pool) Call(arg interface{}) {
	select {
	case o.arg <- arg:
		// ensure <o.arg> is blocked
	default:
		select {
		case o.arg <- arg:
		case <-time.After(o.delay):
			// If <delay> is too small, select may be interfered by the delay caused by gc,
			// and Go may select this case even <o.arg> isn't blocked.
			//
			// A proper value of <o.delay> should at least 0.1s
			if o.newProcess() {
				// todo: ensure the new process is started before try again
				o.Call(arg)
			} else {
				o.arg <- arg
			}
		}
	}
}

// The created coroutine won't quit unless time out. Set min to 0 if want to quit all
func (o *Pool) newProcess() bool {
	if atomic.AddUint32(&o.cur, 1) > o.max {
		atomic.AddUint32(&o.cur, ^uint32(0))
		return false
	}

	var loop *Loop
	loop = NewLoop(func() {
		select {
		case arg := <-o.arg:
			o.fun(arg)
		case <-time.After(o.timeout): //if idle for <timeout> ns, quit
			if atomic.AddUint32(&o.cur, ^uint32(0)) >= o.min {
				loop.Stop()
			} else {
				atomic.AddUint32(&o.cur, 1)
			}
		}
	})

	return true
}
