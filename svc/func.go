/*

FuncWrap (Any -> ())	: "loop f(arg)"
	Call Any : "sched f(arg)"

*** e.g.

// 1. define a new type derive Func
type T struct {
	*Func
}

// 2. define construction
func NewF() T {
	return &F{FuncWrap(func(argAny interface{}) {
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
	"github.com/moxitrel/golib"
	"math"
)

// Loop running fun(arg) in a new goroutine.
type Func struct {
	*Loop
	fun  func(interface{})
	args chan interface{} // argument buffer
}

type _StopSignal struct{}

// Make a new Func service.
// argsCap: the max number of argument can be buffered.
// fun: panic if nil.
func NewFunc(argsCap uint, fun func(interface{})) (o *Func) {
	if fun == nil {
		golib.Panic("fun == nil, want !nil")
	}

	o = &Func{
		fun:  fun,
		args: make(chan interface{}, argsCap),
	}
	o.Loop = NewLoop(func() {
		arg := <-o.args
	CALL:
		if arg != (_StopSignal{}) {
			o.fun(arg)
		}
		select {
		case arg = <-o.args:
			// when .Stop(), continue to handle delivered args,
			// or client may be blocked at .Call()
			goto CALL
		default:
			// return
		}
	})
	return
}

func FuncWrap(fun func(interface{})) (func(interface{}), func()) {
	v := NewFunc(math.MaxUint16, fun)
	return v.Call, v.Stop
}

func (o *Func) Stop() {
	if o.State() != STOPPED {
		o.Loop.Stop()
		o.args <- _StopSignal{}
	}
}

func (o *Func) Call(arg interface{}) {
	if o.State() == STOPPED {
		golib.Warn("%v is stopped", o)
		return
	}
	o.args <- arg
}

//func (o *Func) WithSize(argsCap uint) *Func {
//	if argsCap == uint(cap(o.args)) {
//		return o
//	}
//	old := *o
//	*o = *NewFunc(argsCap, o.fun)
//	old.Stop()
//	return o
//}
