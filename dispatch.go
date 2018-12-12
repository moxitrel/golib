package golib

import "unsafe"

type DispatchMessage interface {
	DispatchKey() interface{}
}

// zero value is not a valid key
type EmbeddedDispatchKey struct {
	dispatcher unsafe.Pointer
	key        uintptr
}
func (o EmbeddedDispatchKey) DispatchKey() interface{} {
	return o
}


///// deprecated /////
type DispatchKeyMessage struct {
	dispatchKey EmbeddedDispatchKey
	message     interface{}
}

func NewDispatchKeyMessage(key EmbeddedDispatchKey, message interface{}) DispatchKeyMessage {
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
