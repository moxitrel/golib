/*
NewLoop ^f	: "Loop running f() in background."
		Stop: "Signal service to stop."
		Join: "Wait service to stop."
*/
package svc

import (
	"github.com/moxitrel/golib"
	"sync"
	"sync/atomic"
)

// State
const (
	STOPPED = iota
	RUNNING
)

// Loop running thunk() in a new goroutine.
type Loop struct {
	//thunk func()
	state uintptr
	wg    sync.WaitGroup
}

func (o *Loop) getState() uintptr {
	return atomic.LoadUintptr(&o.state)
}
func (o *Loop) setState(state uintptr) {
	atomic.StoreUintptr(&o.state, state)
}

// Make a new Loop.
// thunk: panic if nil.
func NewLoop(thunk func()) (o *Loop) {
	if thunk == nil {
		golib.Panic("thunk == nil, want !nil")
	}

	o = &Loop{
		state: RUNNING,
		wg:    sync.WaitGroup{},
	}

	o.wg.Add(1)
	go func() {
		for o.getState() == RUNNING {
			thunk()
		}
		o.wg.Done()
	}()

	return
}

// Get current running state.
func (o *Loop) State() uintptr {
	return o.getState()
}

// Signal to stop running. May not stop immediately.
func (o *Loop) Stop() {
	o.setState(STOPPED)
}

// Block the current goroutine until stopped.
func (o *Loop) Join() {
	o.wg.Wait()
}
