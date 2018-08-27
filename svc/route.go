/*

NewRoute:
	Register   	x cb:
	Apply     	x   : "sched cb(x)"

*/
package svc

import (
	"reflect"
	"github.com/moxitrel/golib"
	"sync"
)

var validateMapKeyCache = new(sync.Map)

func ValidateMapKey(keyType reflect.Type) (v bool) {
	v = true

	switch keyType.Kind() {
	case reflect.Invalid, reflect.Func, reflect.Slice, reflect.Map:
		v = false
	case reflect.Struct:
		if anyV, ok := validateMapKeyCache.Load(keyType); ok {
			v = anyV.(bool)
		} else {
			for i := 0; i < keyType.NumField(); i++ {
				if ValidateMapKey(keyType.Field(i).Type) == false {
					v = false
					break
				}
			}
			validateMapKeyCache.Store(keyType, v)
		}
	}

	return
}

// not thread-safe
type Handler map[interface{}]func(interface{})

func NewHandler() Handler {
	return make(Handler)
}

// arg: arg's type shoudn't be function, slice, map or struct contains function, slice or map field
// fun: nil, delete the handler for arg
func (o Handler) Register(arg interface{}, fun func(interface{})) {
	// assert invalid key type
	if ValidateMapKey(reflect.TypeOf(arg)) == false {
		golib.Panic("%t isn't a valid map key type!\n", arg)
		return
	}

	if fun == nil {
		// delete handler
		delete(o, arg)
	} else {
		if o[arg] != nil {
			golib.Warn("%v is already registered and will be overridden!\n", arg)
		}
		o[arg] = fun
	}
}

func (o Handler) Apply(arg interface{}) {
	// skip invalid key type
	if ValidateMapKey(reflect.TypeOf(arg)) == false {
		golib.Warn("%t isn't a valid map key type!\n", arg)
		return
	}

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
