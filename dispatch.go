package golib

import "unsafe"

// zero value is not a valid key
type DispatchKey struct {
	dispatcher unsafe.Pointer
	key        uintptr
}
func (o DispatchKey) DispatchKey() interface{} {
	return o
}

type DispatchMessage interface {
	DispatchKey() interface{}
}

type DispatchKeyMessage struct {
	dispatchKey DispatchKey
	message     interface{}
}

func WithDispatchKey(key DispatchKey, message interface{}) DispatchKeyMessage {
	return DispatchKeyMessage{
		dispatchKey: key,
		message:     message,
	}
}
func (o DispatchKeyMessage) DispatchKey() interface{} {
	return o.dispatchKey
}
func (o DispatchKeyMessage) Message() interface{} {
	return o.message
}
