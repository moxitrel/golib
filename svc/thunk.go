/*

NewThunk    	    : "Finish f() one by one in background."
	Add thunk		: "schedule thunk()"
	Stop 			:

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

func (o *Thunk) Add(thunk func()) {
	if thunk == nil {
		return
	}
	o.Fun.Call(thunk)
}
