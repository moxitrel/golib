package svc

import (
	"github.com/moxitrel/golib"
	"math"
)

type Func struct {
	*Loop
	args chan interface{} // argument buffer
}

type _StopSignal struct{}

var stopSignal = _StopSignal{}

// Loop running fun(arg) in a new goroutine.
// argsCap: the max number of argument can be buffered.
// fun: panic if nil.
func NewFunc(argsCap uint, fun func(interface{})) (o *Func) {
	if fun == nil {
		golib.Panic("fun == nil, want !nil")
	}

	o = &Func{
		args: make(chan interface{}, argsCap),
	}
	o.Loop = NewLoop(func() {
		var arg interface{}
		for arg = range o.args {
			switch arg {
			case stopSignal:
				goto handleStop
			default:
				fun(arg)
			}
		}

	handleStop:
		o.Loop.Stop()
		stopTimer := NewTimer()
		// when .Stop(), continue to handle delivered args,
		// or client may be blocked at .Call()
		for {
			select {
			case arg = <-o.args:
			default:
				stopTimer.Start(_STOP_DELAY)
				select {
				case arg = <-o.args:
					stopTimer.Stop()
				case <-stopTimer.C:
					return
				}
			}
			if arg != stopSignal {
				fun(arg)
			}
		}
	})
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
