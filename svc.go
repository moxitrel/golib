/*

NewSvc pre post do -> *Svc : make a new service which loop running <do>
	.Stop          	: signal the service to stop
	.State -> Int	: return the current running state

*/
package gosvc

import (
	"fmt"
	"sync/atomic"
)

// Extend goroutine with state reference.
type Svc struct {
	state int32
}

// Svc.state
const (
	ST_NIL = iota
	ST_RUNNING
	ST_STOPPED
	ST_NA // not available
)

// Start a new goroutine loop running do().
//
// pre: run once before start.
// post: run once after stop.
// do: loop do() when running.
func NewSvc(pre func(), post func(), do func()) (o *Svc) {
	o = &Svc{
		state: ST_RUNNING,
	}
	go func() {
		// update state when finish or panic
		defer o.Stop()

		if pre != nil {
			pre()
		}
		// register post() to run when do() finish or panic
		if post != nil {
			defer post()
		}

		if do != nil {
			for {
				switch o.State() {
				case ST_RUNNING:
					do()
				case ST_STOPPED:
					goto DO_END
				default:
					panic(fmt.Errorf("state:%v isn't valid.", o.State()))
				}
			}
		DO_END:
		}
	}()
	return
}

// Signal service to exit. May not stop immediately.
func (o *Svc) Stop() {
	atomic.StoreInt32(&o.state, ST_STOPPED)
}

// Get current running state.
func (o *Svc) State() int {
	return int(atomic.LoadInt32(&o.state))
}
