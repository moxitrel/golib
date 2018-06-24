/*

NewHandler:
	Register   x cb:
	Do     	   x   : "sched cb(x)"

*/
package svc

import (
	"reflect"
)

type Handler struct {
	*Func
	handlers map[reflect.Type]func(interface{})
}

func NewHandler() (v *Handler) {
	v = &Handler{
		handlers: make(map[reflect.Type]func(interface{})),
	}
	v.Func = NewFunc(func(msg interface{}) {
		handle := v.handlers[reflect.TypeOf(msg)]
		if handle == nil {
			// todo: issue warning
			return
		}
		handle(msg)
	})
	return v
}

func (o *Handler) Register(msg interface{}, handle func(interface{})) {
	o.handlers[reflect.TypeOf(msg)] = handle
}

func (o *Handler) Do(msg interface{}) {
	o.Func.Call(msg)
}
