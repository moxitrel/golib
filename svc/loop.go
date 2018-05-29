/*
NewLoop f	: "Loop f() in background."
		Stop: "Signal service to stop."
		Join: "Wait service to stop."
*/
package svc

import "sync"

const (
	STOPPED = iota
	RUNNING
)

type Loop struct {
	thunk func()
	state int
	wg    sync.WaitGroup
}

func NewLoop(thunk func()) (v *Loop) {
	v = &Loop{
		thunk: thunk,
		state: RUNNING,
		wg:    sync.WaitGroup{},
	}
	if v.thunk == nil {
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

func (o *Loop) Stop() {
	o.state = STOPPED
}

func (o *Loop) Join() {
	o.wg.Wait()
}
