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
	thunk func()
	state uintptr
	wg    sync.WaitGroup
}

// Make a new Loop.
// thunk: panic if nil.
func NewLoop(thunk func()) (o *Loop) {
	if thunk == nil {
		golib.Panic("thunk == nil, want !nil")
	}

	o = &Loop{
		thunk: thunk,
		state: RUNNING,
		wg:    sync.WaitGroup{},
	}
	go func() {
		o.wg.Add(1)
		defer o.wg.Done()
		for atomic.LoadUintptr(&o.state) == RUNNING {
			o.thunk()
		}
	}()
	return
}

// Get current running state
func (o *Loop) State() uintptr {
	return atomic.LoadUintptr(&o.state)
}

// Signal to stop running.
// May not stop immediately.
func (o *Loop) Stop() {
	atomic.StoreUintptr(&o.state, STOPPED)
}

// Block the current goroutine until stopped.
func (o *Loop) Join() {
	o.wg.Wait()
}
