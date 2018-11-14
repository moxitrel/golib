package golib

import (
	"sync"
)

type MapDispatch struct {
	sync.Map
	sync.Mutex
}

func NewMapDispatch() *MapDispatch {
	return &MapDispatch{
		Map:   sync.Map{},
		Mutex: sync.Mutex{},
	}
}

// Panic if key is added.
func (o *MapDispatch) Add(key interface{}, handler func(interface{})) {
	if handler == nil {
		Panic("handler == nil, want !nil")
	}

	o.Lock()
	if _, ok := o.Load(key); ok {
		Panic("key:%v has been registered.", key)
	}
	o.Store(key, handler)
	o.Unlock()
}

func (o *MapDispatch) Set(key interface{}, handler func(interface{})) {
	switch handler {
	case nil:
		o.Delete(key)
	default:
		o.Store(key, handler)
	}
}

func (o *MapDispatch) Call(key interface{}, arg interface{}) {
	handler, _ := o.Load(key)
	switch handler {
	case nil:
		//Warn("%#v, the handler doesn't exist!", key)
	default:
		handler.(func(interface{}))(arg)
	}
}

func (o *MapDispatch) Get(key interface{}) (v func(interface{})) {
	handler, _ := o.Load(key)
	if handler != nil {
		v = handler.(func(interface{}))
	}
	return
}
