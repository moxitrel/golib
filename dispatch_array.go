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

type DispatchMessage struct {
	DispatchKey DispatchKey
	Message interface{}
}

// fixed size, add only
type ArrayDispatch struct {
	handers []func(interface{})
	index   uintptr
}

func NewArrayDispatch(size uintptr) *ArrayDispatch {
	return &ArrayDispatch{
		handers: make([]func(interface{}), size),
		index:   0,
	}
}

func (o *ArrayDispatch) Add(handler func(interface{})) (dispatchKey DispatchKey) {
	dispatchKey.key = atomic.AddUintptr(&o.index, 1)
	if dispatchKey.key >= uintptr(len(o.handers)) {
		Panic("Exceed the max size.")
	}
	dispatchKey.dispatcher = unsafe.Pointer(o)
	o.handers[dispatchKey.key] = handler
	return
}

func (o *ArrayDispatch) Get(dispatchKey DispatchKey) func(interface{}) {
	if dispatchKey.dispatcher != unsafe.Pointer(o) {
		Panic("dispatchKey:%v isn't valid", dispatchKey)
	}
	if dispatchKey.key == 0 {
		Panic("dispatchKey:%v isn't valid", dispatchKey)
	}
	return o.handers[dispatchKey.key]
}

func (o *ArrayDispatch) Call(dispatchKey DispatchKey, arg interface{}) {
	fun := o.Get(dispatchKey)
	fun(arg)
}
