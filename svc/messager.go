/*

NewDispatch:
	Register   x cb:
	Handle     x   : "sched cb(x)"

*/
package svc

import (
	"reflect"
)

type Messager struct {
	Function
	handlers map[reflect.Type]func(interface{})
}

func NewDispatch() (v *Messager) {
	v = &Messager{
		handlers: make(map[reflect.Type]func(interface{})),
	}
	v.Function = *NewFunction(func(msg interface{}) {
		handle := v.handlers[reflect.TypeOf(msg)]
		if handle == nil {
			// todo: issue warning
			return
		}
		handle(msg)
	})
	return v
}

func (o *Messager) Register(msg interface{}, handle func(interface{})) {
	o.handlers[reflect.TypeOf(msg)] = handle
}

func (o *Messager) Handle(msg interface{}) {
	o.Function.Call(msg)
}
