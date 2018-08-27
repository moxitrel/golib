/*

NewFunc (Any -> ())	: "loop f(arg)"
	Apply Any : "sched f(arg)"

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

// 3. override Apply() with desired type
func (o *T) Apply(x ArgT) {
	o.Func.Apply(x)
}

*/
package svc

import (
	"sync"
	"github.com/moxitrel/golib"
)

type Func struct {
	*Loop
	fun      func(interface{})
	args     chan interface{}
	stopOnce *sync.Once
}

type _StopSignal struct{}

func NewFunc(bufferCapacity uint, fun func(arg interface{})) (v *Func) {
	if fun == nil {
		golib.Warn("<fun> shouldn't be nil!\n")
		fun = func(_ interface{}) {}
	}

	v = &Func{
		fun:      fun,
		args:     make(chan interface{}, bufferCapacity),
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
	return
}

func (o *Func) Stop() {
	o.stopOnce.Do(func() {
		o.Loop.Stop()
		o.args <- _StopSignal{}
	})
}

func (o *Func) Apply(arg interface{}) {
	if o.state == RUNNING {
		o.args <- arg
	}
}
