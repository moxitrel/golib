/*
NewLoop ^f	: "Loop running f() in background."
		Stop: "Signal service to stop."
		Join: "Wait service to stop."
*/
package svc

import (
	"sync/atomic"
)

// State
const (
	NIL = iota
	RUNNING
	STOPPED
)

type Svc struct {
	state uintptr
}

func NewSvc(pre func(), do func(), post func()) (o *Svc) {
	o = &Svc{
		state: RUNNING,
	}
	go func() {
		if pre != nil {
			pre()
		}
		if do != nil {
			for o.State() == RUNNING {
				do()
			}
		}
		if post != nil {
			post()
		}
	}()
	return
}

// Signal to stop running. May not stop immediately.
func (o *Svc) Stop() {
	atomic.StoreUintptr(&o.state, STOPPED)
}

// Get current running state.
func (o *Svc) State() uintptr {
	return atomic.LoadUintptr(&o.state)
}
