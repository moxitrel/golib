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

// Loop running thunk() in a new goroutine.
type Loop struct {
	*Svc
	wg sync.WaitGroup
}

// Make a new Loop.
// thunk: panic if nil.
func NewLoop(thunk func()) (o *Loop) {
	if thunk == nil {
		golib.Panic("thunk == nil, want !nil")
	}
	o = new(Loop)
	o.wg.Add(1)
	o.Svc = NewSvc(nil, thunk, func() {
		o.wg.Done()
	})
	return
}

// Block the current goroutine until stopped.
func (o *Loop) Join() {
	o.wg.Wait()
}
