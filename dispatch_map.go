package golib

import (
	"sync"
)

type MapDispatch struct {
	sync.Map
}

func NewMapDispatch() *MapDispatch {
	return &MapDispatch{
		Map: sync.Map{},
	}
}

func (o *MapDispatch) Add(key interface{}, handler func(interface{})) {
	if handler == nil {
		return
	}
	handlers, _ := o.Load(key)
	if handlers == nil {
		o.Store(key, []func(interface{}){handler})
	} else {
		o.Store(key, append(handlers.([]func(interface{})), handler))
	}
}

func (o *MapDispatch) Set(key interface{}, handler func(interface{})) {
	if handler == nil {
		o.Delete(key)
	} else {
		o.Store(key, []func(interface{}){handler})
	}
}

func (o *MapDispatch) Call(key interface{}, arg interface{}) {
	handlers, _ := o.Load(key)
	if handlers == nil {
		return
	}
	for _, handler := range handlers.([]func(interface{})) {
		handler(arg)
	}
}
