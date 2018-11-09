package golib

import (
	"sync/atomic"
)

// fixed size, no deletion
type ArrayDispatch struct {
	pool    []func(interface{})
	poolLen uintptr
}

func NewArrayDispatch(size uint) *ArrayDispatch {
	return &ArrayDispatch{
		pool:    make([]func(interface{}), size),
		poolLen: 0,
	}
}

// Add an handler into dispatcher, return the handler's key
// handler: shouldn't be nil
func (o *ArrayDispatch) Add(handler func(interface{})) (index uintptr) {
	if handler == nil {
		Panic("handler == nil, want !nil")
	}
	poolLen := atomic.AddUintptr(&o.poolLen, 1)
	if poolLen > uintptr(len(o.pool)) {
		Panic("pool:%v is full", len(o.pool))
	}
	index = poolLen - 1
	o.pool[index] = handler
	return
}

// Apply the <arg> with the function has the key <index>
// index: should be the value return from Add(), or panic
func (o *ArrayDispatch) UnsafeCall(index uintptr, arg interface{}) {
	handler := o.pool[index]
	handler(arg)
}

//func (o *SliceDispatch) Call(index uintptr, arg interface{}) {
//	poolLen := atomic.LoadUintptr(&o.poolLen)
//	if index >= poolLen {
//		Panic("index:%v is out of range:%v", index, poolLen)
//	}
//	handler := o.pool[index]
//	if handler == nil {
//		Panic("%v, hasn't inited", index)
//	}
//	handler(arg)
//}
