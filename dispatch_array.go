package golib

import (
	"sync/atomic"
)

// fixed size, data race exists if not care
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

// index: it's better to not mixed use with Add(). Use the different index to avoid data race
func (o *ArrayDispatch) Set(index uintptr, handler func(interface{})) {
	if index >= uintptr(len(o.pool)) {
		Panic("index:%v is out of range:%v", index, len(o.pool))
	}
	atomic.StoreUintptr(&o.poolLen, index)
	o.pool[index] = handler
	return
}

// Apply the <arg> with the function has the key <index>
// index: must be the value returned from Add()
func (o *ArrayDispatch) Call(index uintptr, arg interface{}) {
	poolLen := atomic.LoadUintptr(&o.poolLen)
	if index >= poolLen {
		Panic("index:%v is out of range:%v", index, poolLen)
	}
	switch handler := o.pool[index]; handler {
	case nil:
		//Panic("%v isn't registered", index)
	default:
		handler(arg)
	}
}

// index: must be the value returned from Add()
func (o *ArrayDispatch) Get(index uintptr) func(interface{}) {
	poolLen := atomic.LoadUintptr(&o.poolLen)
	if index >= poolLen {
		Panic("index:%v is out of range:%v", index, poolLen)
	}
	return o.pool[index]
}
