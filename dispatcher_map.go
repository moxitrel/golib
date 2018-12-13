package golib

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type MapDispatcher struct {
	sync.Map
	key uintptr
}

func NewMapDispatcher() *MapDispatcher {
	return new(MapDispatcher)
}

func (o *MapDispatcher) Add(handler func(interface{})) (dispatchKey AddDispatchKey) {
	dispatchKey.key = atomic.AddUintptr(&o.key, 1)
	if dispatchKey.key == 0 {
		Panic("No key left.")
	}
	dispatchKey.dispatcher = unsafe.Pointer(o)
	o.Store(dispatchKey, handler)
	return
}

func (o *MapDispatcher) Set(key interface{}, handler func(interface{})) {
	switch handler {
	case nil:
		o.Delete(key)
	default:
		o.Store(key, handler)
	}
}

func (o *MapDispatcher) Get(key interface{}) func(interface{}) {
	handler, _ := o.Load(key)
	if handler == nil {
		return nil
	}
	return handler.(func(interface{}))
}

func (o *MapDispatcher) Call(key interface{}, arg interface{}) {
	fun := o.Get(key)
	if fun == nil {
		Warn("%v: handler doesn't exist", key)
		return
	}
	fun(arg)
}
