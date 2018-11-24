package golib

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// zero value is not a valid key
type DispatchKey struct {
	ptr unsafe.Pointer
	uintptr
}

type MapDispatch struct {
	sync.Map
	uintptr
}

func NewMapDispatch() *MapDispatch {
	return &MapDispatch{
		Map:     sync.Map{},
		uintptr: 0,
	}
}

func (o *MapDispatch) NewKey() DispatchKey {
	index := atomic.AddUintptr(&o.uintptr, 1)
	if index == 0 {
		Panic("All key is used.")
	}
	return DispatchKey{
		ptr:     unsafe.Pointer(o),
		uintptr: index,
	}
}

func (o *MapDispatch) Set(key interface{}, handler func(interface{})) {
	switch handler {
	case nil:
		o.Delete(key)
	default:
		o.Store(key, handler)
	}
}

func (o *MapDispatch) Call(key interface{}, arg interface{}) {
	handler, _ := o.Load(key)
	switch handler {
	case nil:
		//Warn("%#v, the handler doesn't exist!", key)
	default:
		handler.(func(interface{}))(arg)
	}
}

func (o *MapDispatch) Get(key interface{}) (v func(interface{})) {
	handler, _ := o.Load(key)
	if handler != nil {
		v = handler.(func(interface{}))
	}
	return
}
