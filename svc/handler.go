/*

NewHandlerService:
	Set   	 x cb: "add handler for x"
	Handle   x   : "sched cb(x)"

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

// key: key's type shoudn't be function, slice, map or struct contains function, slice or map field
// fun: nil, delete the handler for key
func (o Handler) Set(key interface{}, fun func(arg interface{})) {
	// assert invalid key type
	if ValidateMapKey(reflect.TypeOf(key)) == false {
		golib.Panic("%t isn't a valid map key type!\n", key)
		return
	}

	if fun == nil {
		// delete handler
		delete(o, key)
	} else {
		o[key] = fun
	}
}

func (o Handler) Get(key interface{}) (fun func(arg interface{})) {
	if ValidateMapKey(reflect.TypeOf(key)) == false {
		golib.Warn("%t isn't a valid map key type!\n", key)
		return
	}
	return o[key]
}

// Skip check, may be panic if things ars not expected
func (o Handler) HandleWithoutCheckout(key interface{}, arg interface{}) {
	fun := o[key]
	fun(arg)
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
	o.HandlerService.Handle(arg.Key(), arg)
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
		v.Handler.HandleWithoutCheckout(key, arg)
	})
	return
}

func (o HandlerService) Handle(key interface{}, arg interface{}) {
	if o.Get(key) == nil {
		golib.Warn("%v doesn't has a handler!\n", key)
		return
	}
	o.HandleWithoutCheck(key, arg)
}

func (o HandlerService) HandleWithoutCheck(key interface{}, arg interface{}) {
	o.FuncService.Call([]interface{}{key, arg})
}