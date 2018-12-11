package golib

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type MapDispatch struct {
	sync.Map
	key uintptr
}

func NewMapDispatch() *MapDispatch {
	return new(MapDispatch)
}

func (o *MapDispatch) Add(handler func(interface{})) (dispatchKey DispatchKey) {
	dispatchKey.key = atomic.AddUintptr(&o.key, 1)
	if dispatchKey.key == 0 {
		Panic("No key left.")
	}
	dispatchKey.dispatcher = unsafe.Pointer(o)
	o.Store(dispatchKey, handler)
	return
}

func (o *MapDispatch) Set(key interface{}, handler func(interface{})) {
	switch handler {
	case nil:
		o.Delete(key)
	default:
		o.Store(key, handler)
	}
}

func (o *MapDispatch) Get(key interface{}) func(interface{}) {
	handler, _ := o.Load(key)
	if handler == nil {
		return nil
	}
	return handler.(func(interface{}))
}

func (o *MapDispatch) Call(key interface{}, arg interface{}) {
	fun := o.Get(key)
	if fun == nil {
		Panic("%v: handler doesn't exist", key)
	}
	fun(arg)
}
