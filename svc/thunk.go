/*

NewThunk    	    :
	Do thunk		: "sched thunkService()"
	Stop 			: "stop the service"

*/
package svc

type Thunk struct {
	Fun
}

func NewThunk() *Thunk {
	f := func(x interface{}) {
		thunk := x.(func())
		thunk()
	}
	return &Thunk{*NewFun(f)}
}

func (o *Thunk) Do(thunk func()) {
	if thunk == nil {
		return
	}
	o.Fun.Call(thunk)
}
