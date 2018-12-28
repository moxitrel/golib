/*

NewSvc pre post do: *Svc
	.Stop
	.State: Int

*/
package svc

import (
	"sync/atomic"
)

// Svc.state
const (
	NIL = iota
	RUNNING
	PAUSED
	STOPPED
)

type Svc struct {
	state int32
}

// Start a goroutine loop running do().
//
// pre: run once before start
// post: run once after stop
// do: loop do() when RUNNING
func NewSvc(pre func(), post func(), do func()) (o *Svc) {
	o = &Svc{
		state: RUNNING,
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
				case RUNNING:
					do()
				case PAUSED:
					// TODO
				default:
					return
				}
			}
		}
	}()
	return
}

// Signal service to exit. May not stop immediately.
func (o *Svc) Stop() {
	atomic.StoreInt32(&o.state, STOPPED)
}

// Get current running state.
func (o *Svc) State() int32 {
	return atomic.LoadInt32(&o.state)
}
