/*

NewDispatch:
	Register   x cb:
	Handle     x   : "schedule cb(x)"

*/
package svc

import (
	"reflect"
)

type Dispatch struct {
	Fun
	handlers map[reflect.Type]func(interface{})
}

func NewDispatch() (v *Dispatch) {
	v = &Dispatch{
		handlers: make(map[reflect.Type]func(interface{})),
	}
	v.Fun = *NewFun(func(msg interface{}) {
		handle := v.handlers[reflect.TypeOf(msg)]
		if handle == nil {
			return
		}
		handle(msg)
	})
	return v
}

func (o *Dispatch) Register(msg interface{}, handle func(interface{})) {
	o.handlers[reflect.TypeOf(msg)] = handle
}

func (o *Dispatch) Handle(msg interface{}) {
	o.Fun.Call(msg)
}
