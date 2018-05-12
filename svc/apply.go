/*

NewApply    	    : "Finish f() one by one in background."
	Add (f & (:))	: "add new task"
	Start			:
	Stop 			: "Service.thunk() won't quit if blocked by receiving"

*/
package svc

type Apply struct {
	Fun
}

func NewApply() (v *Apply) {
	v = &Apply{
		*NewFun(func(thunk interface{}) {
			thunk.(func())()
		}),
	}
	return
}

// Blocked if pool is full.
func (o *Apply) Add(thunk func()) {
	if thunk == nil {
		return
	}
	o.Call(thunk)
}
