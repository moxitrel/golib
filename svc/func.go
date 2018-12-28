/*

NewFunc size fun: *Func
	.Stop
	.Wait
	.Call arg
	.State: Int

*/

package svc

import (
	"github.com/moxitrel/golib"
	"time"
)

const (
	// time to wait for receiving sent args when receive stop signal
	_STOP_DELAY = 200 * time.Millisecond
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
		// continue to handle delivered args when .Stop(), or client may be blocked
		for {
			select {
			case arg = <-o.args:
			default:
				stopTimer.Start(_STOP_DELAY)
				select {
				case arg = <-o.args:
					stopTimer.Stop()
				case <-stopTimer.C: // quit if timeout
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

//
// Wrap time.Timer
//
type Timer struct {
	*time.Timer
}

func NewTimer() (o Timer) {
	o = Timer{
		Timer: time.NewTimer(time.Second),
	}
	o.Stop()
	return
}

func (o Timer) Start(timeout time.Duration) {
	o.Reset(timeout)
}

func (o Timer) Stop() {
	if !o.Timer.Stop() {
		<-o.C
	}
}
