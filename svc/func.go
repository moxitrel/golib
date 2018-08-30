/*

NewFuncService (Any -> ())	: "loop f(arg)"
	Apply Any : "sched f(arg)"

*** e.g.

// 1. define a new type derive FuncService
type T struct {
	*FuncService
}

// 2. define construction
func NewF() T {
	return &F{NewFuncService(func(argAny interface{}) {
		arg := argAny.(ArgT)	// recover the type
		...
	})}
}

// 3. override Apply() with desired type
func (o *T) Apply(x ArgT) {
	o.FuncService.Apply(x)
}

*/
package svc

import (
	"github.com/moxitrel/golib"
	"sync"
)

type FuncService struct {
	*LoopService
	fun      func(interface{})
	args     chan interface{}
	stopOnce *sync.Once
}

type _StopSignal struct{}

func NewFuncService(bufferCapacity uint, fun func(arg interface{})) (v *FuncService) {
	if fun == nil {
		golib.Warn("<fun> shouldn't be nil!\n")
		fun = func(_ interface{}) {}
	}

	v = &FuncService{
		fun:      fun,
		args:     make(chan interface{}, bufferCapacity),
		stopOnce: new(sync.Once),
	}
	v.LoopService = NewLoopService(func() {
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

func (o *FuncService) Stop() {
	o.stopOnce.Do(func() {
		o.LoopService.Stop()
		o.args <- _StopSignal{}
	})
}

func (o *FuncService) Call(arg interface{}) {
	if o.state == RUNNING {
		o.args <- arg
	}
}
