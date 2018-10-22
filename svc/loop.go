/*
NewLoop ^f	: "Loop running f() in background."
		Stop: "Signal service to stop."
		Join: "Wait service to stop."
*/
package svc

import (
	"github.com/moxitrel/golib"
	"sync"
)

// State
const (
	STOPPED = iota
	RUNNING
)

// Loop Service: loop running thunk() in a new goroutine.
type Loop struct {
	thunk func()
	state int
	wg    *sync.WaitGroup
}

// Return a new loop service.
// thunk: panic if nil.
func NewLoop(thunk func()) (v *Loop) {
	if thunk == nil {
		golib.Panic("^thunk shouldn't be nil!")
	}

	v = &Loop{
		thunk: thunk,
		state: RUNNING,
		wg:    new(sync.WaitGroup),
	}
	go func() {
		v.wg.Add(1)
		defer v.wg.Done()
		for v.state == RUNNING {
			v.thunk()
		}
	}()
	return
}

// Signal the loop to stop running.
// May not stop immediately.
func (o *Loop) Stop() {
	o.state = STOPPED
}

// Block the current goroutine until the loop stopped.
func (o *Loop) Join() {
	o.wg.Wait()
}
