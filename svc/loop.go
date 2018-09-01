/*
NewLoopService f	: "LoopService f() in background."
		Stop: "Signal service to stop."
		Join: "Wait service to stop."
*/
package svc

import (
	"github.com/moxitrel/golib"
	"sync"
)

const (
	STOPPED = iota
	RUNNING
)

type LoopService struct {
	thunk func()
	state int
	wg    *sync.WaitGroup
}

func NewLoopService(thunk func()) (v *LoopService) {
	v = &LoopService{
		thunk: thunk,
		state: RUNNING,
		wg:    new(sync.WaitGroup),
	}
	if v.thunk == nil {
		golib.Warn("^thunk shouldn't be nil!\n")
		return
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

func (o *LoopService) Stop() {
	o.state = STOPPED
}

func (o *LoopService) Join() {
	o.wg.Wait()
}
