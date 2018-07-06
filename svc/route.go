/*

NewRoute:
	Register   x cb:
	Call     	   x   : "sched cb(x)"

*/
package svc

import (
	"reflect"
)

type Route struct {
	*Func
	funs map[reflect.Type]func(interface{})
}

func NewRoute() (v Route) {
	v = Route{
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

func (o *Route) Register(arg interface{}, fun func(interface{})) {
	o.funs[reflect.TypeOf(arg)] = fun
}

func (o *Route) Call(arg interface{}) {
	o.Func.Call(arg)
}
