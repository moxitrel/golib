/*
NewLoop f	: "Loop f() in background."
		Stop: "Signal service to stop."
*/
package svc

const (
	STOPPED = iota
	RUNNING
)

type Loop struct {
	thunk func()
	state int
}

func NewLoop(thunk func()) (v Loop) {
	v = Loop{
		thunk: thunk,
		state: RUNNING,
	}
	if v.thunk == nil {
		// todo: issue warning or panic
		return
	}
	go func() {
		for v.state == RUNNING {
			v.thunk()
		}
	}()
	return
}

func (o *Loop) Stop() {
	o.state = STOPPED
}
