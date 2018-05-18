/*
New (f & (:)): 	"Loop f()."
	Start: 		"Loop in background."
	Stop : 		"Signal to stop."
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

func New(thunk func()) *Service {
	return &Service{
		thunk:     thunk,
		state:     STOPPED,
		stateLock: *new(sync.Mutex),
	}
}

// Run the service in background.
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

	//defer func() { recover() }()	// stop panic if occurred
	go o.loop()
}

// Signal the service to stop.
// Service stopped after the thunk() done.
func (o *Service) Stop() {
	o.state = STOPPED
}

// require thunk != nil
func (o *Service) loop() {
	for o.state == RUNNING {
		o.thunk()
	}
}
