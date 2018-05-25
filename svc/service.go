/*
New f	: "Loop f() in background."
	Stop: "Signal service to stop."
*/
package svc

const (
	STOPPED = iota
	RUNNING
)

type Service struct {
	state int
}

func New(thunk func()) (v Service) {
	v = Service{
		state: RUNNING,
	}
	if thunk == nil {
		// todo: issue warning
		return
	}
	go func() {
		for v.state == RUNNING {
			thunk()
		}
	}()
	return
}

func (o *Service) Stop() {
	o.state = STOPPED
}
