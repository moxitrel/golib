/*
func New(func()) Service
func (Service) Start()
func (Service) Stop()
*/

package svc

import (
	"github.com/moxitrel/golib/contract"
)

const (
	STOPPED = iota
	RUNNING
)

type Service struct {
	thunk     func()
	state     int
}

func New(thunk func()) Service {
	contract.Assert(thunk != nil, "thunk: should not be nil")
	return Service{
		thunk:     thunk,
		state:     STOPPED,
	}
}

func (o Service) Start() {
	// single instance, not thread-safe
	if o.state == RUNNING {
		return
	}
	o.state = RUNNING

	go func() {
		//defer recover()	// stop panic if occurred
		o.loop()
	}()
}

func (o Service) Stop() {
	o.state = STOPPED
}

func (o Service) loop() {
	for o.state == RUNNING {
		o.thunk()
	}
}
