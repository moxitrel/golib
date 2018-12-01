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

// pre: run once when start
// do: loop do() when RUNNING
// post: run once after stop
func NewSvc(pre func(), post func(), do func()) (o *Svc) {
	o = &Svc{
		state: RUNNING,
	}
	go func() {
		// update state if panic or nil
		defer o.Stop()

		if pre != nil {
			pre()
		}
		// run post() even do() panic, but not pre() panic
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
