/*
func NewApply(poolSize uint) Apply
func (*Apply) Start()
func (*Apply) Stop()
func (*Apply) Add(func())			: add new task
*/
package svc

type Apply struct {
	*Service
	thunks chan func() //should not be closed
}

func NewApply(poolSize uint) (v Apply) {
	v.thunks = make(chan func(), poolSize)
	v.Service = New(func() {
		//thunk, ok := <-v.thunks
		//if !ok {
		//	v.Stop()
		//	return
		//}
		thunk, _ := <-v.thunks
		if thunk == nil {
			return
		}
		thunk()
	})
	return
}

func (o *Apply) Stop() {
	o.Service.Stop()

	// quit from thunk if blocked by receiving
	if len(o.thunks) == 0 {
		// parallel Stop() without block
		select {
		case o.thunks <- func() {}:
		default:
		}
	}
}

// Blocked if pool is full.
func (o *Apply) Add(thunk func()) {
	if thunk == nil {
		return
	}
	o.thunks <- thunk
}
