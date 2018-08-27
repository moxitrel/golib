/*

NewRoute:
	Register   	x cb:
	Apply     	x   : "sched cb(x)"

*/
package svc

import (
	"reflect"
	"github.com/moxitrel/golib"
)

var validateMapKeyCache = make(map[reflect.Type]bool)

func validateMapKey(keyType reflect.Type) (v bool) {
	v = true
	switch keyType.Kind() {
	case reflect.Invalid:
		v = false
	case reflect.Func:
		v = false
	case reflect.Slice:
		v = false
	case reflect.Map:
		v = false
	case reflect.Struct:
		if v_, ok := validateMapKeyCache[keyType]; ok {
			v = v_
		} else {
			for i := 0; i < keyType.NumField(); i++ {
				if validateMapKey(keyType.Field(i).Type) == false {
					v = false
					break
				}
			}
		}
	}

	validateMapKeyCache[keyType] = v
	return
}

type Handler map[interface{}]func(interface{})

func NewHandler() Handler {
	return make(Handler)
}

func (o Handler) Register(arg interface{}, fun func(interface{})) {
	// ignore invalid key type
	if validateMapKey(reflect.TypeOf(arg)) == false {
		golib.Panic("%t isn't a valid map key type!\n", arg)
		return
	}
	if fun == nil {
		golib.Warn("<fun> shouldn't be nil!\n")
	}

	o[arg] = fun
}

func (o Handler) Apply(arg interface{}) {
	// skip invalid key type
	//if validateMapKey(reflect.TypeOf(arg)) == false {
	//	golib.Warn("%t isn't a valid map key type!\n", arg)
	//	return
	//}

	fun := o[arg]
	if fun == nil {
		golib.Warn("%v doesn't has a handler!\n", arg)
		return
	}

	fun(arg)
}

type Route struct {
	*Func
	Handler
}

func NewRoute(bufferCapacity uint) (v Route) {
	v.Handler = NewHandler()
	v.Func = NewFunc(bufferCapacity, v.Handler.Apply)
	return
}

func (o *Route) Apply(arg interface{}) {
	o.Func.Apply(arg)
}
