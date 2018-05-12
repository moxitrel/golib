package svc

import (
	"reflect"

	"gitlab.com/clogwire/v4/log"
)

type Messager struct {
	*Service
	handlers map[reflect.Type]func(interface{})
	storage  chan interface{}
}

func NewMessager() (v *Messager) {
	v = new(Messager)
	v.handlers = make(map[reflect.Type]func(interface{}))
	v.storage = make(chan interface{}, 256)
	v.Service = New(func() {
		msg := <-v.storage
		v.handlers[reflect.TypeOf(msg)](msg)
	})
	return v
}

func (o *Messager) Register(msg interface{}, handler interface{}) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil {
		log.Fatal("msg shouldn't nil interface{}")
		return
	}
	th := reflect.TypeOf(handler)
	if !(th.Kind() == reflect.Func &&
		th.NumIn() == 1 &&
		(th.In(0) == msgType || th.In(0).Kind() == reflect.Interface && msgType.Implements(th.In(0)))) {

		log.Fatal("%t isn't a valid func", handler)
		return
	}

	o.handlers[msgType] = func(x interface{}) {
		reflect.ValueOf(handler).Call([]reflect.Value{reflect.ValueOf(x)})
	}
}

func (o *Messager) AddMessage(message interface{}) {
	if o.handlers[reflect.TypeOf(message)] == nil {
		log.Warn("handler for %t doesn't exist", message)
		return
	}
	o.storage <- message
}
