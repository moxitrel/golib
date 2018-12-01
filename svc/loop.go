package svc

import (
	"github.com/moxitrel/golib"
	"sync"
)

type Loop struct {
	*Svc
	wg sync.WaitGroup
}

// Loop running thunk() in a new goroutine.
func NewLoop(thunk func()) (o *Loop) {
	if thunk == nil {
		golib.Panic("thunk == nil, want !nil")
	}
	o = new(Loop)
	o.wg.Add(1)
	o.Svc = NewSvc(
		nil,
		func() { o.wg.Done() },
		thunk)
	return
}

// Block current goroutine until stopped.
func (o *Loop) Wait() {
	o.wg.Wait()
}
