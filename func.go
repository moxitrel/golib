/*

NewFunc size fun -> *Func : loop processing fun(arg) in a new goroutine
	.Stop			: signal the service to stop
	.Wait			: wait until stopped
	.Call arg		: send ^arg to service to process
	.State -> Int	: return the current running state

*/
package gosvc

import (
	"time"
)

type Func struct {
	*Loop
	args chan interface{} // argument buffer
}

// time to wait for sent args when received stop-signal
const _FUNC_STOP_DELAY = 200 * time.Millisecond

type _FuncStopSignal struct{}

var funcStopSignal = _FuncStopSignal{}

// Loop running fun(arg) in a new goroutine.
//
// argsCap: the max number of argument can be buffered.
// fun: panic if nil.
func NewFunc(argsCap uint, fun func(interface{})) (o *Func) {
	if fun == nil {
		panic("fun == nil, want !nil")
	}

	o = &Func{
		args: make(chan interface{}, argsCap),
	}
	o.Loop = NewLoop(func() {
		var arg interface{}
		for arg = range o.args {
			switch arg {
			case funcStopSignal:
				goto HANDLE_STOP
			default:
				fun(arg)
			}
		}

	HANDLE_STOP:
		o.Loop.Stop()
		stopTimer := NewTimer()
		// continue to handle delivered args when .Stop(), or client may be blocked
		for {
			select {
			case arg = <-o.args:
			default:
				stopTimer.Start(_FUNC_STOP_DELAY)
				select {
				case arg = <-o.args:
				case <-stopTimer.C: // quit if timeout
					goto HANDLE_STOP_END
				}
				stopTimer.Stop()
			}

			if arg != funcStopSignal {
				fun(arg)
			}
		}
	HANDLE_STOP_END:
	})
	return
}

// Signal service to exit. May not stop immediately.
func (o *Func) Stop() {
	if o.State() != ST_STOPPED {
		o.args <- funcStopSignal
	}
}

// Send ^arg to process
func (o *Func) Call(arg interface{}) {
	if o.State() == ST_RUNNING {
		o.args <- arg
	}
}
