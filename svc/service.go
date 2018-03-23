/*
func New(func()) Service
func (*Service) Start()
func (*Service) Stop()
*/

package svc

import "sync"

const (
	STOPPED = iota
	RUNNING
)

type Service struct {
	thunk     func()
	state     int
	stateLock sync.Mutex
}

func New(thunk func()) Service {
	return Service{
		thunk:     thunk,
		state:     STOPPED,
		stateLock: sync.Mutex{},
	}
}

func (o *Service) Start() {
	// implement single instance
	o.stateLock.Lock()
	if o.state == RUNNING {
		return
	}
	o.state = RUNNING
	o.stateLock.Unlock()

	if o.thunk == nil {
		return
	}

	//defer recover()	// stop panic if occurred
	go o._loop()
}

func (o *Service) Stop() {
	o.state = STOPPED
}

// require: thunk != nil
func (o *Service) _loop() {
	for o.state == RUNNING {
		o.thunk()
	}
}
