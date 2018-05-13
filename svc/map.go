/*

NewMap:
	Register   arg f:
	Call       arg  : "run f(arg)"

*/
package svc

import (
	"reflect"

	"gitlab.com/clogwire/v4/log"
)

type Map struct {
	Fun
	handlers map[reflect.Type]func(interface{})
}

func NewMap() (v *Map) {
	v = &Map{
		handlers: make(map[reflect.Type]func(interface{})),
	}
	v.Fun = *NewFun(func(argv []interface{}) {
		msg := argv[0]
		handler := v.handlers[reflect.TypeOf(msg)]
		handler(msg)
	})
	return v
}

func (o *Map) Register(x interface{}, f func(interface{})) {
	o.handlers[reflect.TypeOf(x)] = f
}

func (o *Map) Call(x interface{}) {
	if o.handlers[reflect.TypeOf(x)] == nil {
		log.Warn("handler for %t doesn't exist", x)
		return
	}
	o.Fun.Call(x)
}
