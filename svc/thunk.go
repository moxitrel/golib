/*

NewThunk    	    :
	Do thunk		: "sched thunk()"

*/
package svc

type Thunk struct {
	Function
}

func NewThunk() *Thunk {
	return &Thunk{*NewFunction(func(thunkAny interface{}) {
		thunk := thunkAny.(func())
		thunk()
	})}
}

func (o *Thunk) Do(thunk func()) {
	if thunk == nil {
		return
	}
	o.Function.Call(thunk)
}
