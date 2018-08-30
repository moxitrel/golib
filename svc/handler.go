/*

NewHandlerService bufferSize:
	Set   	 x cb: "add handler for x.Key"
	Handle   x   : "sched cb(x)"

*/
package svc

import (
	"github.com/moxitrel/golib"
	"reflect"
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
	handlers map[interface{}]func(interface{})
}

func NewHandlerService(bufferCapacity uint) (v HandlerService) {
	v.handlers = make(map[interface{}]func(interface{}))
	v.FuncService = NewFuncService(bufferCapacity, func(anyKeyArg interface{}) {
		keyArg := anyKeyArg.([]interface{})
		key := keyArg[0]
		arg := keyArg[1]

		fun := v.handlers[key]
		fun(arg)
	})
	return
}

// key: key's type shoudn't be function, slice, map or struct contains function, slice or map field
func (o HandlerService) Set(key interface{}, fun func(arg interface{})) {
	// assert invalid key type
	if ValidateMapKey(reflect.TypeOf(key)) == false {
		golib.Panic("%t isn't a valid map key type!\n", key)
		return
	}
	o.SetWithoutCheck(key, fun)
}
func (o HandlerService) Handle(key interface{}, arg interface{}) {
	if ValidateMapKey(reflect.TypeOf(key)) == false {
		golib.Warn("%t isn't a valid map key type!\n", key)
		return
	}
	if o.handlers[key] == nil {
		golib.Warn("%v, handler doesn't exist!\n", key)
		return
	}
	o.HandleWithoutCheck(key, arg)
}
func (o HandlerService) SetWithoutCheck(key interface{}, fun func(arg interface{})) {
	if fun == nil {
		// delete handler
		delete(o.handlers, key)
	} else {
		o.handlers[key] = fun
	}
}
func (o HandlerService) HandleWithoutCheck(key interface{}, arg interface{}) {
	o.FuncService.Call([]interface{}{key, arg})
}
