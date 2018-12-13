package golib

import (
	"sync/atomic"
	"unsafe"
)

// fixed size, add only
type ArrayDispatcher struct {
	handers []func(interface{})
	index   uintptr
}

// NOTE: the first one is reserved. So the max number of handlers can be added is size - 1
func NewArrayDispatcher(size uintptr) *ArrayDispatcher {
	return &ArrayDispatcher{
		handers: make([]func(interface{}), size),
		index:   0,
	}
}

func (o *ArrayDispatcher) Add(handler func(interface{})) (dispatchKey AddDispatchKey) {
	dispatchKey.key = atomic.AddUintptr(&o.index, 1)
	if dispatchKey.key >= uintptr(len(o.handers)) {
		Panic("Exceed the max size.")
	}
	dispatchKey.dispatcher = unsafe.Pointer(o)
	o.handers[dispatchKey.key] = handler
	return
}

func (o *ArrayDispatcher) Get(dispatchKey AddDispatchKey) func(interface{}) {
	if dispatchKey.dispatcher != unsafe.Pointer(o) {
		Panic("dispatchKey.dispatcher:%v isn't valid, want %v", dispatchKey.dispatcher, o)
	}
	if dispatchKey.key == 0 {
		Panic("dispatchKey.key == 0, want !0", dispatchKey)
	}
	return o.handers[dispatchKey.key]
}

// Use the dispatchKey returned from Add()
func (o *ArrayDispatcher) Call(dispatchKey AddDispatchKey, arg interface{}) {
	fun := o.Get(dispatchKey)
	fun(arg)
}
