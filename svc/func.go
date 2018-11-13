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
	"sync"
	"sync/atomic"
)

// Loop running fun(arg) in a new goroutine.
type Func struct {
	*Loop
	//fun  func(interface{})
	args chan interface{} // argument buffer
}

func (o *Func) getState() uintptr {
	return atomic.LoadUintptr(&o.state)
}
func (o *Func) setState(state uintptr) {
	atomic.StoreUintptr(&o.state, state)
}

// Make a new Func service.
// argsCap: the max number of argument can be buffered.
// fun: panic if nil.
func NewFunc(argsCap uint, fun func(interface{})) (o Func) {
	if fun == nil {
		golib.Panic("fun == nil, want !nil")
	}

	o = Func{
		Loop: &Loop{
			state: RUNNING,
			wg:    sync.WaitGroup{},
		},
		args: make(chan interface{}, argsCap),
	}
	o.wg.Add(1)
	go func() {
		defer o.wg.Done()
		stopTimer := NewTimer()
		for {
			switch o.getState() {
			case STOPPED:
				// when .Stop(), continue to handle delivered args,
				// or client may be blocked at .Call()
				stopTimer.Start(_STOP_DELAY)
				select {
				case arg := <-o.args:
					stopTimer.Stop()
					if arg != stopSignal {
						fun(arg)
					}
				case <-stopTimer.C:
					return
				}
			default:
				switch arg := <-o.args; arg {
				case stopSignal:
					o.setState(STOPPED)
				default:
					fun(arg)
				}
			}
		}
	}()
	return
}

func FuncWrapper(fun func(interface{})) (func(interface{}), func()) {
	v := NewFunc(math.MaxUint16, fun)
	return v.Call, v.Stop
}

func (o *Func) Stop() {
	if o.State() != STOPPED {
		o.args <- stopSignal
	}
}

func (o *Func) Call(arg interface{}) {
	if o.State() != STOPPED {
		o.args <- arg
	}
}
