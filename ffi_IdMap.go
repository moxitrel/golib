package golib

import (
	"sync"
	"sync/atomic"
)

///
/// Map Go pointer object to an unique number
///
/// C code may not keep a copy of a Go pointer after the call returns.
/// Type always include Go pointers: string, channel, function types, interface, slice, map, _GoString_
///
type IdMap struct {
	sync.Map
	uintptr
}

func NewIdMap() *IdMap {
	o := new(IdMap)
	o.Store(0, nil)
	return o
}

func (o *IdMap) Add(ptr interface{}) (id uintptr) {
	if ptr == nil {
		return 0
	}
	id = atomic.AddUintptr(&o.uintptr, 1)
	if id == 0 {
		Panic("id out of range.")
	}
	o.Store(id, ptr)
	return
}

func (o *IdMap) Get(id uintptr) interface{} {
	val, _ := o.Load(id)
	return val
}
