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

func (o *MapDispatch) Set(key interface{}, handler func(interface{})) {
	switch handler {
	case nil:
		o.Delete(key)
	default:
		o.Store(key, handler)
	}
}

func (o *MapDispatch) Get(key interface{}) (v func(interface{})) {
	handler, _ := o.Load(key)
	if handler != nil {
		v = handler.(func(interface{}))
	}
	return
}

func (o *MapDispatch) Call(key interface{}, arg interface{}) {
	handler, _ := o.Load(key)
	switch handler {
	case nil:
		Warn("%#v, handler doesn't exist!", key)
	default:
		handler.(func(interface{}))(arg)
	}
}
