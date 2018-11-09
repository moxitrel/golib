package golib

import (
	"github.com/moxitrel/golib"
	"math"
	"sync/atomic"
	"testing"
)

// fixed size, no deletion
type _AtomicArrayDispatch struct {
	pool    []*func(interface{})
	poolLen uintptr
}

func newAtomicArrayDispatch(size uint) *_AtomicArrayDispatch {
	return &_AtomicArrayDispatch{
		pool:    make([]*func(interface{}), size),
		poolLen: 0,
	}
}

// handler: shouldn't be nil
func (o *_AtomicArrayDispatch) Add(handler func(interface{})) (index uintptr) {
	if handler == nil {
		golib.Panic("handler == nil, want !nil")
	}
	poolLen := atomic.AddUintptr(&o.poolLen, 1)
	if poolLen > uintptr(len(o.pool)) {
		golib.Panic("pool:%v is full", len(o.pool))
	}
	index = poolLen - 1
	o.pool[index] = &handler
	return
}

// index: should be the value return from Add(), or panic
func (o *_AtomicArrayDispatch) UnsafeCall(index uintptr, arg interface{}) {
	handler := o.pool[index]
	(*handler)(arg)
}

func BenchmarkArrayDispatch(b *testing.B) {
	o := NewArrayDispatch(math.MaxInt8)

	index := o.Add(func(interface{}) {})
	for i := 0; i < b.N; i++ {
		o.UnsafeCall(index, nil)
	}
}

func BenchmarkAtomicArrayDispatch(b *testing.B) {
	o := newAtomicArrayDispatch(math.MaxInt8)

	index := o.Add(func(interface{}) {})
	for i := 0; i < b.N; i++ {
		o.UnsafeCall(index, nil)
	}
}
