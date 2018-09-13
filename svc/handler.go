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

// Cache the check result
var validateMapKeyCache = new(sync.Map)

// Check whether <keyType> is a valid key type of map
func ValidateMapKey(keyType reflect.Type) (v bool) {
	v = true

	if keyType != nil {
		switch keyType.Kind() {
		case reflect.Invalid, reflect.Func, reflect.Slice, reflect.Map:
			// shoudn't be a function, slice or map
			v = false
		case reflect.Struct:	// shoudn't be a struct contains function, slice or map field
			// fetch result from cache if exists
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

/* Usage:

// 1. define a new type derive HandlerService
//
type MyHandlerService struct {
	HandlerService
}

// 2. define method Call() with arg has key() attribute
//
func (o MyHandlerService) Call(arg T) {
	o.HandlerService.Handle(arg.Key(), arg)
}

*/
// Process arg by handlers[arg.key()] in goroutine pool
type HandlerService struct {
	*FuncService
	*Pool
	handlers map[interface{}]func(interface{})
}

func NewHandlerService(bufferSize uint, poolMin uint) (v *HandlerService) {
	v = new(HandlerService)
	v.handlers = make(map[interface{}]func(interface{}))
	v.Pool = NewPool(poolMin, func(anyKeyArg interface{}) {
		keyArg := anyKeyArg.([]interface{})
		key := keyArg[0]
		arg := keyArg[1]

		fun := v.handlers[key]
		fun(arg)
	})
	v.FuncService = NewFuncService(bufferSize, v.Pool.Call)
	return
}

// key: key's type shoudn't be function, slice, map or struct contains function, slice or map field
// fun: nil, delete the handler according to key
func (o *HandlerService) Set(key interface{}, fun func(arg interface{})) {
	// assert invalid key type
	if ValidateMapKey(reflect.TypeOf(key)) == false {
		golib.Panic("%t isn't a valid map key type!\n", key)
		return
	}

	if fun == nil {
		// delete handler
		delete(o.handlers, key)
	} else {
		o.handlers[key] = fun
	}
}
func (o *HandlerService) Handle(key interface{}, arg interface{}) {
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

func (o *HandlerService) HandleWithoutCheck(key interface{}, arg interface{}) {
	o.FuncService.Call([]interface{}{key, arg})
}
