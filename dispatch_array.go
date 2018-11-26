package golib

import (
	"sync/atomic"
	"unsafe"
)

// fixed size, add only
type ArrayDispatch struct {
	pool []func(interface{})
	uintptr
}

func NewArrayDispatch(size uintptr) *ArrayDispatch {
	return &ArrayDispatch{
		pool:    make([]func(interface{}), size),
		uintptr: 0,
	}
}

func (o *ArrayDispatch) NewKey() DispatchKey {
	index := atomic.AddUintptr(&o.uintptr, 1)
	if index >= uintptr(len(o.pool)) {
		Panic("All key is used.")
	}
	return DispatchKey{
		dispatcher: unsafe.Pointer(o),
		uintptr:    index,
	}
}

// Add an handler into dispatcher, return the handler's key
// handler: shouldn't be nil
func (o *ArrayDispatch) Set(key DispatchKey, handler func(interface{})) {
	if key == (DispatchKey{}) {
		Panic("key:%v isn't valid", key)
	}
	if handler == nil {
		Panic("handler == nil, want !nil")
	}
	o.pool[key.uintptr] = handler
}

// Apply the <arg> with the function has the key <index>
// index: must be the value returned from Add()
func (o *ArrayDispatch) Call(key DispatchKey, arg interface{}) {
	if key == (DispatchKey{}) {
		Panic("key:%v isn't valid", key)
	}
	handler := o.pool[key.uintptr]
	handler(arg)
}
