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
)

type Function struct {
	*Loop
	fun      func(interface{})
	args     chan interface{}
	stopOnce *sync.Once
}

type _StopSignal struct{}

var _STOP_SIGNAL = _StopSignal{}

//func FunctionOf(f func(interface{})) (func(interface{}), func()) {
//	xs := make(chan interface{}, math.MaxInt32)
//	o := NewLoop(func() {
//		// apply args until emtpy
//		for x := range xs {
//			if x != _STOP_SIGNAL {
//				f(x)
//			}
//
//			if len(xs) == 0 {
//				break
//			}
//		}
//	})
//	return func(x interface{}) {
//			if o.state == RUNNING {
//				xs <- x
//			}
//		}, func() {
//			v.stopOnce.Do(func() {
//				o.Stop()
//				xs <- _STOP_SIGNAL
//			})
//		}
//}

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
