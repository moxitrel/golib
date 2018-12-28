/*

NewLoop do: *Loop
	.Stop
	.Wait
	.State: Int

*/
package svc

import (
	"sync"
)

type Loop struct {
	*Svc
	wg sync.WaitGroup
}

// Loop running thunk() in a new goroutine.
// If thunk is nil, stop immediately.
func NewLoop(thunk func()) (o *Loop) {
	o = new(Loop)
	o.wg.Add(1)
	o.Svc = NewSvc(
		nil,
		func() { o.wg.Done() },
		thunk)
	return
}

// Block current goroutine until Loop stopped.
func (o *Loop) Wait() {
	o.wg.Wait()
}
