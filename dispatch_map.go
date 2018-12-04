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

func (o *MapDispatch) Add(handler interface{}) (dispatchKey DispatchKey) {
	dispatchKey.key = atomic.AddUintptr(&o.key, 1)
	if dispatchKey.key == 0 {
		Panic("No key left.")
	}
	dispatchKey.dispatcher = unsafe.Pointer(o)
	o.Store(dispatchKey, handler)
	return
}

func (o *MapDispatch) Set(key interface{}, handler interface{}) {
	switch handler {
	case nil:
		o.Delete(key)
	default:
		o.Store(key, handler)
	}
}

func (o *MapDispatch) Get(key interface{}) interface{} {
	handler, _ := o.Load(key)
	return handler
}
