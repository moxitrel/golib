/*

NewThunk    	    : "Finish f() one by one in background."
	Add (f & (:))	: "add new task f()."
	Start			:
	Stop 			:

*/
package svc

type Thunk struct {
	Fun
}

func NewThunk() (v *Thunk) {
	v = &Thunk{
		*NewFun(func(argv interface{}) {
			thunk := argv.(func())
			thunk()
		}),
	}
	return
}

// Blocked if pool is full.
func (o *Thunk) Add(thunk func()) {
	if thunk == nil {
		return
	}
	o.Call(thunk)
}
