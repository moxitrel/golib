package golib

import (
	"sync/atomic"
)

// zero value is not a valid key
type ArrayDispatchKey struct {
	uintptr
}

// fixed size, add only
type ArrayDispatch struct {
	pool    []func(interface{})
	poolLen uintptr
}

func NewArrayDispatch(size uintptr) *ArrayDispatch {
	return &ArrayDispatch{
		pool:    make([]func(interface{}), size),
		poolLen: 0,
	}
}

// Add an handler into dispatcher, return the handler's key
// handler: shouldn't be nil
func (o *ArrayDispatch) Add(handler func(interface{})) ArrayDispatchKey {
	if handler == nil {
		Panic("handler == nil, want !nil")
	}
	poolLen := atomic.AddUintptr(&o.poolLen, 1)
	if poolLen > uintptr(len(o.pool)) {
		Panic("pool:%v is full", len(o.pool))
	}
	o.pool[poolLen] = handler
	return ArrayDispatchKey{
		uintptr: poolLen,
	}
}

// Apply the <arg> with the function has the key <index>
// index: must be the value returned from Add()
func (o *ArrayDispatch) Call(key ArrayDispatchKey, arg interface{}) {
	if key == (ArrayDispatchKey{}) {
		Panic("key:%v isn't valid", key)
	}
	handler := o.pool[key.uintptr]
	handler(arg)
}

// don't use mix with .Add()
// exist data race if index is the same
// index: require > 0
//func (o *ArrayDispatch) Set(index uintptr, handler func(interface{})) ArrayDispatchKey {
//	if index == 0 {
//		Panic("index == 0, want !0")
//	}
//	if index >= uintptr(len(o.pool)) {
//		Panic("index:%v is out of range:%v", index, len(o.pool))
//	}
//	if handler == nil {
//		Panic("handler == nil, want !nil")
//	}
//	atomic.StoreUintptr(&o.poolLen, index)
//	o.pool[index] = handler
//	return ArrayDispatchKey{
//		uintptr: index,
//	}
//}
