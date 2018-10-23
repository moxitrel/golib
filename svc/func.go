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
	"math"
)

const defaultArgsSize = math.MaxUint16

type Func struct {
	*Loop
	fun  func(interface{})
	args chan interface{}
}

type _StopSignal struct{}

func NewFunc(fun func(interface{})) (v *Func) {
	return NewFuncWithSize(defaultArgsSize, fun)
}

func NewFuncWithSize(argsCap uint, fun func(interface{})) (v *Func) {
	if fun == nil {
		golib.Panic("^fun shouldn't be nil!\n")
	}
	v = &Func{
		fun:  fun,
		args: make(chan interface{}, argsCap),
	}
	v.Loop = NewLoop(func() {
		arg := <-v.args
	APPLY:
		if arg != (_StopSignal{}) {
			v.fun(arg)
		}
		select {
		case arg = <-v.args:
			// when Stop(), continue to handle delivered args,
			// or client may be blocked at .Apply()
			goto APPLY
		default:
			// return
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

func (o *Func) Apply(arg interface{}) {
	if o.state == RUNNING {
		o.args <- arg
	}
}

func (o *Func) WithSize(argsCap uint) *Func {
	if argsCap == uint(cap(o.args)) {
		return o
	}
	old := *o
	*o = *NewFuncWithSize(argsCap, o.fun)
	old.Stop()
	return o
}
