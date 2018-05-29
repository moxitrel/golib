/*

NewFunction f	: derive Loop
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
	"time"
)

type Fun struct {
	fun      func(interface{})
	arg      chan interface{}
	maxDelay time.Duration //arg must be handled after maxDelay
}

// Return a started fun-service
// fun: apply with arg passed from Call()
func NewFun(delay time.Duration, fun func(arg interface{})) (v *Fun) {
	v = &Fun{
		fun:      fun,
		arg:      make(chan interface{}),
		maxDelay: delay,
	}
	if fun != nil {
		go func() {
			for {
				v.fun(<-v.arg)
			}
		}()
	}
	return
}

func (o *Fun) Call(arg interface{}) {
	select {
	case o.arg <- arg:
	case <-time.After(o.maxDelay):
		go func() {
			for {
				select {
				case arg := <-o.arg:
					o.fun(arg)
				case <-time.After(time.Minute):
					break
				}
			}
		}()
	}
}
