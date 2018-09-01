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
)

type FuncService struct {
	*LoopService
	fun  func(interface{})
	args chan interface{}
}

type _StopSignal struct{}

func NewFuncService(bufferSize uint, fun func(interface{})) (v *FuncService) {
	v = &FuncService{
		fun:  fun,
		args: make(chan interface{}, bufferSize),
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
	if fun == nil {
		golib.Warn("^fun shouldn't be nil!\n")
		v.Stop()
	}
	return
}

func (o *FuncService) Stop() {
	if o.state == RUNNING {
		o.LoopService.Stop()
		o.args <- _StopSignal{}
	}
}

func (o *FuncService) Call(arg interface{}) {
	if o.state == RUNNING {
		o.args <- arg
	}
}

// XXX: may be unsafe, element order in channel my be change while Call()
func (o *FuncService) SetSize(n uint) {
	oldArgs := o.args
	o.args = make(chan interface{}, n)
	oldArgs <- _StopSignal{}
	for {
		select {
		case arg := <-oldArgs:
			o.args <- arg
		default:
			break
		}
	}
}
