/*

NewLoop do -> *Loop	: make a new service which loop running ^do
	.Stop			: signal the service to stop
	.Wait			: wait until stopped
	.State -> Int	: return the current running state

*/
package gosvc

import (
	"sync"
)

// Extend goroutine with state and stop-wait support.
type Loop struct {
	*Svc
	wg sync.WaitGroup
}

// Loop running thunk() in a new goroutine.
//
// If thunk is nil, stop immediately.
func NewLoop(thunk func()) (o *Loop) {
	o = &Loop{}
	o.wg.Add(1)
	o.Svc = NewSvc(
		nil,
		o.wg.Done,
		thunk)
	return
}

// Block current goroutine until stopped.
func (o *Loop) Wait() {
	o.wg.Wait()
}
