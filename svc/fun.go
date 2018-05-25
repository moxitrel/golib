/*

NewFun f		:
	Call arg	: "sched f(arg)"
	Stop        : "stop the service"

*** e.g.

// 1. define a new type derive Fun
type T struct {
	Fun
}

// 2. define construction
func NewF() *T {
	// 2.1. define a function with signature func(interface{})
	f := func (arg interface{}) {
		x := arg.(ArgT)		//recover type first
		...					//do the things
	}

	// 2.2. return
	return &F{ *NewFun(f) }
}

// 3. override Call() with desired argument type
func (o *T) Call(x ArgT) {
	o.Fun.Call(x)
}

*/
package svc

import (
	"sync"
)

type Fun struct {
	args     chan interface{}
	stopOnce sync.Once
}

// Return a started fun-service
// fun: apply with arg passed from Call()
func NewFun(fun func(arg interface{})) (v *Fun) {
	v = &Fun{
		args:     make(chan interface{}, FUN_BUFFER_SIZE),
		stopOnce: sync.Once{},
	}
	go func() {
		if fun == nil {
			// todo: issue warning
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

func (o *Fun) Stop() {
	o.stopOnce.Do(func() {
		close(o.args)
	})
}

func (o *Fun) Call(arg interface{}) {
	o.args <- arg
}
