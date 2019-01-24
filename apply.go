/*

NewApply size fun -> *Apply : loop processing fun(arg) in a new goroutine
	.Stop			: signal the service to stop
	.Wait			: wait until stopped
	.Call arg		: send ^arg to service to process
	.State -> Int	: return the current running state

*/
package gosvc

import (
	"fmt"
	"time"
)

type Apply struct {
	*Loop
	args chan interface{} // argument buffer
}

type _FuncStopSignal struct{}

var funcStopSignal = _FuncStopSignal{}

// Loop running fun(arg) in a new goroutine.
//
// argsCap: the max number of argument can be buffered.
// fun: panic if nil.
func NewApply(argsCap uint, fun func(interface{})) (o *Apply) {
	if fun == nil {
		panic(fmt.Errorf("fun == nil, want !nil"))
	}

	o = &Apply{
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
				stopTimer.Start(200 * time.Millisecond)
				select {
				case arg = <-o.args:
				case <-stopTimer.C: // quit if timeout
					return
				}
				stopTimer.Stop()
			}

			if arg != funcStopSignal {
				fun(arg)
			}
		}
	})
	return
}

// Signal service to exit. May not stop immediately.
func (o *Apply) Stop() {
	if o.State() != ST_STOPPED {
		o.args <- funcStopSignal
	}
}

// Send ^arg to process
func (o *Apply) Call(arg interface{}) {
	if o.State() == ST_RUNNING {
		o.args <- arg
	}
}
