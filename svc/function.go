/*

NewFunction n f	: "loop f(arg)"
	Call arg	: "sched f(arg)"

*** e.g.

// 1. define a new type derive Function
type T struct {
	Function
}

// 2. define construction
func NewF() *T {
	// 2.1. define the function
	f := func (arg ArgT) {
		...
	}

	// 2.2. wrap f with signature func(interface{})
	return &F{*NewFunction(func(arg interface{}) {
		f(arg.(ArgT))	//2.3. recover the type
	})}
}

// 3. override Call() with desired argument type
func (o *T) Call(x ArgT) {
	o.Function.Call(x)
}

*/
package svc

import (
	"sync"
	"sync/atomic"
	"time"
)

type Function struct {
	*Loop
	fun      func(interface{})
	args     chan interface{}
	stopOnce *sync.Once
}

type _StopSignal struct{}
var _STOP_SIGNAL = _StopSignal{}

func NewFunction(maxArgs uint, fun func(arg interface{})) (v *Function) {
	v = &Function{
		fun:      fun,
		args:     make(chan interface{}, maxArgs),
		stopOnce: new(sync.Once),
	}
	v.Loop = NewLoop(func() {
		// apply args until emtpy
		for arg := range v.args {
			if arg != _STOP_SIGNAL {
				v.fun(arg)
			}

			if len(v.args) == 0 {
				break
			}
		}
	})
	if fun == nil {
		v.Stop()
	}
	return
}

func (o *Function) Stop() {
	o.stopOnce.Do(func() {
		o.Loop.Stop()
		o.args <- _STOP_SIGNAL
	})
}

func (o *Function) Call(arg interface{}) {
	if o.state == RUNNING {
		o.args <- arg
	}
}

// min    : at least <min> coroutines will be created and live all the time
// max    : the max number of coroutines can be created
// delay  : create a new coroutine if arg is blocked for <delay> ns
// timeout: destroy the coroutine if it's idle for <timeout> ns
//
// created coroutine won't quit until time out. Set *min to 0 if want to quit all
func PoolOf(fun func(interface{}), min *uint, max *uint, delay *time.Duration, timeout *time.Duration) func(interface{}) {
	if min == nil {
		defaultMin := uint(0)
		min = &defaultMin
	}
	if max == nil {
		defaultMax := *min
		min = &defaultMax
	}
	if delay == nil {
		defaultDelay := time.Duration(0)
		delay = &defaultDelay
	}
	if timeout == nil {
		defaultTimeout := *delay
		timeout = &defaultTimeout
	}

	var argChannel = make(chan interface{})
	var cur int32 //current coroutines count
	var newCoroutine = func() {
		atomic.AddInt32(&cur, 1)
		var loop *Loop
		loop = NewLoop(func() {
			// if idle for <timeout> ns, quit
			select {
			case arg := <-argChannel:
				fun(arg)
			case <-time.After(*timeout):
				if atomic.LoadInt32(&cur) > int32(*min) {
					loop.Stop()
					atomic.AddInt32(&cur, -1)
				}
			}
		})
	}

	for i := int32(0); i < int32(*min); i++ {
		newCoroutine()
	}

	var limitFun func(interface{})
	limitFun = func(arg interface{}) {
		if atomic.LoadInt32(&cur) >= int32(*max) {
			argChannel <- arg
		} else {
			select {
			case argChannel <- arg:
				// The case here is to ensure <argChannel> is blocked
				//
				// Don't it seem the same as the case in default clause?
				// No. If <delay> is a small value, it would be interfered by the delay caused by gc,
				// and Go may choose the second case.
			default:
				select {
				case argChannel <- arg:
				case <-time.After(*delay): //a proper value should at least 0.1s, e.g. 0.5s
					newCoroutine()
					limitFun(arg)
				}
			}
		}
	}
	return limitFun
}
