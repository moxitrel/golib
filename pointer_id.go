package golib

import (
	"sync"
	"sync/atomic"
)

///
/// Map go pointer object to an unique number
///
/// C code may not keep a copy of a Go pointer after the call returns.
/// C code may store Go pointers in C memory, subject to the rule above: it must stop storing the Go pointer when the C function returns
///
/// Type always include Go pointers: string, channel, function types, interface, slice, map, _GoString_
/// A pointer type may hold a Go pointer or a C pointer
///
type IdPointer struct {
	size     uint64   // used as id, Increase only
	pointers sync.Map // map[uint64]PointerObject
}

func NewIdPointer() (o *IdPointer) {
	o = &IdPointer{
		size:     0,
		pointers: sync.Map{},
	}
	o.pointers.Store(o.size, nil)
	return
}

// XXX: require o.size never overflow
func (o *IdPointer) Add(ptr interface{}) (id uint64) {
	switch ptr {
	case nil:
		id = 0
	default:
		id = atomic.AddUint64(&o.size, 1)
		if id == 0 {
			Panic("id > math.MaxUint64, overflow")
		}
		o.pointers.Store(id, ptr)
	}
	return
}

func (o *IdPointer) Delete(id uint64) {
	o.pointers.Delete(id)
}

func (o *IdPointer) Get(id uint64) (interface{}, bool) {
	return o.pointers.Load(id)
}
