package golib

import (
	"sync/atomic"
	"unsafe"
)

// zero value is not a valid key
type DispatchKey struct {
	dispatcher unsafe.Pointer
	key        uintptr
}

// fixed size, add only
type ArrayDispatch struct {
	handers []func(interface{})
	key     uintptr
}

func NewArrayDispatch(size uintptr) *ArrayDispatch {
	return &ArrayDispatch{
		handers: make([]func(interface{}), size),
		key:     0,
	}
}

func (o *ArrayDispatch) Add(handler func(interface{})) (dispatchKey DispatchKey) {
	index := atomic.AddUintptr(&o.key, 1)
	if index >= uintptr(len(o.handers)) {
		Panic("No key left.")
	}
	dispatchKey = DispatchKey{
		dispatcher: unsafe.Pointer(o),
		key:        index,
	}
	o.handers[dispatchKey.key] = handler
	return
}

func (o *ArrayDispatch) Call(key DispatchKey, arg interface{}) {
	if key == (DispatchKey{}) {
		Panic("key:%v isn't valid", key)
	}
	handler := o.handers[key.key]
	handler(arg)
}

func (o *ArrayDispatch) Get(key DispatchKey) func(interface{}) {
	return o.handers[key.key]
}
