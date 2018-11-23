/*

MakeArrayDispatch bufferSize:
	Set    x cb: "add handler for x.Key"
	Call   x   : "sched cb(x)"

*/
package svc

import (
	"github.com/moxitrel/golib"
)

// Cache the check result
//var validateMapKeyCache = new(sync.Map)

// Check whether <keyType> is a valid key type of map
//func ValidateMapKey(keyType reflect.Type) (v bool) {
//	v = true
//
//	if keyType != nil {
//		switch keyType.Kind() {
//		case reflect.Invalid, reflect.Func, reflect.Slice, reflect.Map:
//			// shoudn't be a function, slice or map
//			v = false
//		case reflect.Struct: // shoudn't be a struct contains function, slice or map field
//			// fetch result from cache if exists
//			if anyV, ok := validateMapKeyCache.Load(keyType); ok {
//				v = anyV.(bool)
//			} else {
//				for i := 0; i < keyType.NumField(); i++ {
//					if ValidateMapKey(keyType.Field(i).Type) == false {
//						v = false
//						break
//					}
//				}
//				validateMapKeyCache.Store(keyType, v)
//			}
//		}
//	}
//
//	return
//}

/* Usage:

// 1. define a new type derive Dispatch
//
type MyHandlerService struct {
	Dispatch
}

// 2. define method Call() with arg has key() attribute
//
func (o MyHandlerService) Call(arg T) {
	o.Dispatch.Call(arg.Key(), arg)
}

*/
type DispatchMsg interface {
	DispatchKey() interface{}
}

// Process arg by handlers[arg.key()] in goroutine pool
type Dispatch struct {
	*golib.MapDispatch
	*Pool
	Func
}

func NewDispatch(bufferSize uint, poolMax uint) (o Dispatch) {
	if poolMax <= 0 {
		golib.Panic("poolMax == 0, want > 0")
	}
	o.MapDispatch = golib.NewMapDispatch()
	o.Pool = NewPool(1, poolMax, 0, _POOL_TICKER_INTVL, 0, func(arg interface{}) {
		msg := arg.(DispatchMsg)
		fun := o.MapDispatch.Get(msg.DispatchKey())
		if fun == nil {
			golib.Warn("%#v, handler doesn't exist!", msg)
			return
		}
		fun(msg)
	})
	o.Func = NewFunc(bufferSize, o.Pool.Call)
	return
}

func (o *Dispatch) Stop() {
	if o.state == RUNNING {
		o.Pool.Stop()
		o.Func.Stop()
	}
}

// key: key's type shoudn't be function, slice, map or struct contains function, slice or map field
// fun: nil, delete the handler according to key
func (o *Dispatch) Set(key interface{}, fun func(interface{})) {
	o.MapDispatch.Set(key, fun)
}

//func (o *Dispatch) Get(key interface{}) func(interface{}) {
//	v, _ := o.handlers.Load(key)
//	if v == nil {
//		return nil
//	}
//	return v.(func(interface{}))
//}

func (o *Dispatch) Call(msg DispatchMsg) {
	if msg == nil {
		return
	}
	o.Func.Call(msg)
}
