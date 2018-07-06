/*

NewFunc (Any -> ())	: "loop f(arg)"
	Call Any : "sched f(arg)"

*** e.g.

// 1. define a new type derive Func
type T struct {
	*Func
}

// 2. define construction
func NewF() T {
	return &F{NewFunc(func(argAny interface{}) {
		arg := argAny.(ArgT)	// recover the type
		...
	})}
}

// 3. override Call() with desired type
func (o *T) Call(x ArgT) {
	o.Func.Call(x)
}

*/
package svc

import (
	"sync"
)

type Func struct {
	*Loop
	fun      func(interface{})
	args     chan interface{}
	stopOnce *sync.Once
}

type _StopSignal struct{}

func NewFunc(fun func(arg interface{})) (v *Func) {
	v = &Func{
		fun:      fun,
		args:     make(chan interface{}, FuncArgMax),
		stopOnce: new(sync.Once),
	}
	v.Loop = NewLoop(func() {
		for {
			arg := <-v.args
			if arg != (_StopSignal{}) {
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

func (o *Func) Stop() {
	o.stopOnce.Do(func() {
		o.Loop.Stop()
		o.args <- _StopSignal{}
	})
}

func (o *Func) Call(arg interface{}) {
	if o.state == RUNNING {
		o.args <- arg
	}
}
