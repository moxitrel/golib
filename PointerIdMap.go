package golib

import (
	"sync"
	"sync/atomic"
)

///
/// Map go pointer object to an unique number
///
/// C code may not keep a copy of a Go pointer after the call returns.
/// Type always include Go pointers: string, channel, function types, interface, slice, map, _GoString_
///
type PointerIdMap struct {
	pointers sync.Map // map[uintptr]GoPointer
	id       uintptr  // increase only
}

func NewPointerIdMap() (o *PointerIdMap) {
	o = &PointerIdMap{
		pointers: sync.Map{},
		id:       0,
	}
	// add nil interface
	o.pointers.Store(o.id, nil)
	return
}

func (o *PointerIdMap) Add(ptr interface{}) (id uintptr) {
	switch ptr {
	case nil:
		id = 0
	default:
		id = atomic.AddUintptr(&o.id, 1)
		if id == 0 {
			Panic("id overflow")
		}
		o.pointers.Store(id, ptr)
	}
	return
}

func (o *PointerIdMap) Delete(id uintptr) {
	o.pointers.Delete(id)
}

func (o *PointerIdMap) Get(id uintptr) (interface{}, bool) {
	return o.pointers.Load(id)
}
