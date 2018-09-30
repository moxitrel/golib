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
	"github.com/moxitrel/golib"
)

type Func struct {
	*Loop
	fun  func(interface{})
	args chan interface{}
}

type _StopSignal struct{}

func NewFunc(bufferSize uint, fun func(interface{})) (v *Func) {
	if fun == nil {
		golib.Panic("^fun shouldn't be nil!\n")
	}

	v = &Func{
		fun:  fun,
		args: make(chan interface{}, bufferSize),
	}
	v.Loop = NewLoop(func() {
		arg := <-v.args
		for { // when Stop(), process remaining received args
			if arg != (_StopSignal{}) {
				v.fun(arg)
			}
			select {
			case arg = <-v.args:
			default:
				return // break can't quit for loop
			}
		}
	})
	return
}

func (o *Func) Stop() {
	if o.state == RUNNING {
		o.Loop.Stop()
		o.args <- _StopSignal{}
	}
}

func (o *Func) Call(arg interface{}) {
	if o.state == RUNNING {
		o.args <- arg
	}
}

// XXX: element order in channel can be changed while Call()
//func (o *Func) SetSize(n uint) {
//	oldArgs := o.args
//	o.args = make(chan interface{}, n)
//	oldArgs <- _StopSignal{}
//	for {
//		select {
//		case arg := <-oldArgs:
//			o.args <- arg
//		default:
//			return
//		}
//	}
//}
