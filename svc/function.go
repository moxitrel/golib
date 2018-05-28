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
	Loop
	args     chan interface{}
	stopOnce sync.Once
}

// Return a started fun-service
// fun: apply with arg passed from Call()
func NewFunction(fun func(arg interface{})) (v *Function) {
	v = &Function{
		args:     make(chan interface{}, FunctionBufferSize),
		stopOnce: sync.Once{},
	}
	v.Loop = *NewLoop(func() {
		// do {...} until (...);
		for {
			arg := <-v.args
			if arg != v.args {	//ignore quit recv flag sent by Stop()
				fun(arg)
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
		o.args <- o.args	//quit recv if blocked, unexported field args as a flag
	})
}

func (o *Function) Call(arg interface{}) {
	if o.state == RUNNING {
		o.args <- arg
	}
}
