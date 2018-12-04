package golib

import "unsafe"

///
/// Map go pointer object to an unique number
///
/// C code may not keep a copy of a Go pointer after the call returns.
/// Type always include Go pointers: string, channel, function types, interface, slice, map, _GoString_
///
type IdMap struct {
	MapDispatch
}

func (o *IdMap) Add(ptr interface{}) (id uintptr) {
	return o.MapDispatch.Add(ptr).key
}

func (o *IdMap) Delete(id uintptr) {
	dispatchKey := DispatchKey{
		dispatcher: unsafe.Pointer(&o.MapDispatch),
		key:        id,
	}
	o.Set(dispatchKey, nil)
}

func (o *IdMap) Get(id uintptr) interface{} {
	dispatchKey := DispatchKey{
		dispatcher: unsafe.Pointer(&o.MapDispatch),
		key:        id,
	}
	return o.MapDispatch.Get(dispatchKey)
}
