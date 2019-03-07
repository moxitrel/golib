/*

NewFunc size fun -> *Func : loop processing fun(arg) in a new goroutine
	.Call arg		: send arg to process
	.Stop			: signal the service to stop
	.Wait			: wait until stopped

*/
package gosvc

import (
	"fmt"
	"sync"
	"time"
)

type Func struct {
	*Svc
	args chan interface{} // argument buffer
	wg sync.WaitGroup
}

type _FuncStopSignal struct{}

// Loop running fun(arg) in a new goroutine.
//
// argsCap: the max number of argument can be buffered.
// fun: panic if nil.
func NewFunc(argsCap uint, fun func(interface{})) (o *Func) {
	if fun == nil {
		panic(fmt.Errorf("fun == nil, want !nil"))
	}

	o = &Func{
		args: make(chan interface{}, argsCap),
	}
	o.wg.Add(1)
	o.Svc = NewSvc(nil, o.wg.Done, func() {
		var arg interface{}
		for arg = range o.args {
			switch arg {
			case _FuncStopSignal{}:
				goto HANDLE_STOP
			default:
				fun(arg)
			}
		}

	HANDLE_STOP:
		o.Svc.Stop()
		stopTimer := NewTimer()
		// continue to handle delivered args when .Stop(), or client may be blocked
		for {
			select {
			case arg = <-o.args:
			default:
				stopTimer.Start(100 * time.Millisecond)
				select {
				case arg = <-o.args:
				case <-stopTimer.C: // quit if timeout
					return
				}
				stopTimer.Stop()
			}

			if arg != (_FuncStopSignal{}) {
				fun(arg)
			}
		}
	})
	return
}

// Submit arg.
func (o *Func) Call(arg interface{}) {
	if o.State() == ST_RUNNING {
		o.args <- arg
	}
}

// Signal service to exit. May not stop immediately.
func (o *Func) Stop() {
	if o.State() != ST_STOPPED {
		o.args <- _FuncStopSignal{}
	}
}

// Block current goroutine until stopped.
func (o *Func) Wait() {
	o.wg.Wait()
}
