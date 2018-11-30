package svc

import (
	"sync/atomic"
)

// Svc.state
const (
	NIL = iota
	RUNNING
	STOPPED
)

type Svc struct {
	state int32
}

func NewSvc(pre func(), post func(), do func()) (o *Svc) {
	o = &Svc{
		state: RUNNING,
	}
	go func() {
		// update state if panic
		defer o.Stop()
		// run post() even panic
		if post != nil {
			defer post()
		}

		if pre != nil {
			pre()
		}
		if do != nil {
			for {
				switch o.State() {
				case RUNNING:
					do()
				default:
					goto BREAK_FOR
				}
			}
		BREAK_FOR:
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
