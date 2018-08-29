/*

NewHandlerService:
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

	if keyType != nil {
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
func (o Handler) Register(key interface{}, fun func(interface{})) {
	// assert invalid key type
	if ValidateMapKey(reflect.TypeOf(key)) == false {
		golib.Panic("%t isn't a valid map key type!\n", key)
		return
	}

	if fun == nil {
		// delete handler
		delete(o, key)
	} else {
		if o[key] != nil {
			golib.Warn("%v is already registered and will be overwritten!\n", key)
		}
		o[key] = fun
	}
}

func (o Handler) Handle(key interface{}, arg interface{}) {
	// assert invalid key type
	if ValidateMapKey(reflect.TypeOf(key)) == false {
		golib.Panic("%t isn't a valid map key type!\n", key)
		return
	}

	fun := o[key]
	if fun == nil {
		golib.Warn("%v doesn't has a handler!\n", key)
		return
	}

	fun(key)
}

/*

// 1. define a new type derive HandlerService
//
type MyHandlerService struct {
	HandlerService
}

// 2. override Handle() with pre-defined key() function
//
func (o MyHandlerService) Handle(arg interface{}) {
	o.HandlerService.Handle(key(arg), arg)
}

*/
type HandlerService struct {
	*FuncService
	Handler
}

func NewHandlerService(bufferCapacity uint) (v HandlerService) {
	v.Handler = NewHandler()
	v.FuncService = NewFuncService(bufferCapacity, func(anyKeyArg interface{}) {
		keyArg := anyKeyArg.([]interface{})
		key := keyArg[0]
		arg := keyArg[1]
		v.Handler.Handle(key, arg)
	})
	return
}

func (o *HandlerService) Handle(key interface{}, arg interface{}) {
	o.FuncService.Call([]interface{}{key, arg})
}
