/*

NewFunction f	:
	Call arg	: "sched f(arg)"
	Stop        : "stop the service"

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
	fun      func(interface{})
	args     chan interface{}
	stopOnce sync.Once
}

// Return a started fun-service
// fun: apply with arg passed from Call()
func NewFunction(fun func(arg interface{})) (v *Function) {
	v = &Function{
		fun:      fun,
		args:     make(chan interface{}, FUN_BUFFER_SIZE),
		stopOnce: sync.Once{},
	}
	go func() {
		if fun == nil {
			// todo: issue warning or panic
			for range v.args {
				// do nothing
			}
		} else {
			for arg := range v.args {
				fun(arg)
			}
		}
	}()
	return
}

func (o *Function) Stop() {
	o.stopOnce.Do(func() {
		close(o.args)
	})
}

func (o *Function) Call(arg interface{}) {
	o.args <- arg
}
