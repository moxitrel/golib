package golib

import (
	"reflect"
	"unsafe"
)

type DispatchMessage interface {
	DispatchKey() interface{}
}

// Generated by ArrayDispatcher.Add() or MapDispatcher.Add()
type AddDispatchKey struct {
	dispatcher unsafe.Pointer
	key        uintptr
}

func (o AddDispatchKey) DispatchKey() interface{} {
	return o
}

type TypeDispatchKey struct {
	Key reflect.Type
}

func (o TypeDispatchKey) DispatchKey() interface{} {
	return o.Key
}