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
	handers []interface{}
	key     uintptr
}

func NewArrayDispatch(size uintptr) *ArrayDispatch {
	return &ArrayDispatch{
		handers: make([]interface{}, size),
		key:     0,
	}
}

func (o *ArrayDispatch) Add(handler interface{}) (dispatchKey DispatchKey) {
	dispatchKey.key = atomic.AddUintptr(&o.key, 1)
	if dispatchKey.key >= uintptr(len(o.handers)) {
		Panic("No key left.")
	}
	dispatchKey.dispatcher = unsafe.Pointer(o)
	o.handers[dispatchKey.key] = handler
	return
}

func (o *ArrayDispatch) Get(dispatchKey DispatchKey) interface{} {
	if dispatchKey.dispatcher != unsafe.Pointer(o) {
		Panic("dispatchKey:%v isn't valid", dispatchKey)
	}
	if dispatchKey.key == 0 {
		Panic("dispatchKey:%v isn't valid", dispatchKey)
	}
	return o.handers[dispatchKey.key]
}
