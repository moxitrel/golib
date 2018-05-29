/*

NewThunk    	    :
	Do thunk		: "sched thunk()"

*/
package svc

type Thunk struct {
	Function
	thunk chan func()
}

func NewThunk() *Thunk {
	return &Thunk{
		Function: *NewFunction(FunctionBufferSize, func(thunkAny interface{}) {
			thunk := thunkAny.(func())
			thunk()
		}),
		thunk: make(chan func()),
	}
}

func (o *Thunk) Do(thunk func()) {
	if thunk == nil {
		return
	}
	o.Function.Call(thunk)
}
