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

// Loop running do() in a new goroutine.
type Loop struct {
	//do    func()
	state uintptr
	wg    sync.WaitGroup
}

// Make a new Loop.
// thunk: panic if nil.
func NewLoop(thunk func()) (o *Loop) {
	if thunk == nil {
		golib.Panic("thunk == nil, want !nil")
	}
	return NewHookedLoop(nil, thunk, nil)
}

func NewHookedLoop(pre func(), do func(), post func()) (o *Loop) {
	o = &Loop{
		state: RUNNING,
		wg:    sync.WaitGroup{},
	}

	o.wg.Add(1)
	go func() {
		defer o.wg.Done()
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

// Get current running state.
func (o *Loop) State() uintptr {
	return atomic.LoadUintptr(&o.state)
}

// Signal to stop running. May not stop immediately.
func (o *Loop) Stop() {
	atomic.StoreUintptr(&o.state, STOPPED)
}

// Block the current goroutine until stopped.
func (o *Loop) Join() {
	o.wg.Wait()
}
