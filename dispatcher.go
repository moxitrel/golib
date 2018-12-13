package golib

import (
	"reflect"
	"unsafe"
)

type DispatchMessage interface {
	DispatchKey() interface{}
}

// Generated by ArrayDispatcher or MapDispatcher
type AutoDispatchKey struct {
	dispatcher unsafe.Pointer
	key        uintptr
}

func (o AutoDispatchKey) DispatchKey() interface{} {
	return o
}

type TypeDispatchKey struct {
	Key reflect.Type
}

func (o TypeDispatchKey) DispatchKey() interface{} {
	return o.Key
}
