/*

NewMap:
	Register   x cb:
	Call     	   x   : "sched cb(x)"

*/
package svc

import (
	"reflect"
)

type Map struct {
	*Func
	funs map[reflect.Type]func(interface{})
}

func NewMap() (v Map) {
	v = Map{
		funs: make(map[reflect.Type]func(interface{})),
	}
	v.Func = NewFunc(func(arg interface{}) {
		fun := v.funs[reflect.TypeOf(arg)]
		if fun == nil {
			// todo: issue warning
			return
		}
		fun(arg)
	})
	return v
}

func (o *Map) Register(arg interface{}, fun func(interface{})) {
	o.funs[reflect.TypeOf(arg)] = fun
}

func (o *Map) Call(arg interface{}) {
	o.Func.Call(arg)
}
