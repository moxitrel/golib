/*
func NewApply(poolSize uint) Apply
func (Apply) Start()
func (Apply) Stop()
func (Apply) Add(func())
*/

package svc

type Apply struct {
	Service
	thunks chan func()
}

func NewApply(poolSize uint) (v Apply) {
	v.thunks = make(chan func(), poolSize)
	v.Service = New(func() {
		f, ok := <-v.thunks
		if !ok {
			v.Stop()
			return
		}
		if f == nil {
			return
		}
		f()
	})
	return
}

func (o Apply) Add(thunk func()) {
	if thunk == nil {
		return
	}
	o.thunks <- thunk
}
