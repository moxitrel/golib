package svc

import (
	"reflect"
	"testing"
	"time"
)

func TestHandlerService_Example(t *testing.T) {
	key := func(arg interface{}) reflect.Type {
		return reflect.TypeOf(arg)
	}
	v := ""

	o := NewHandlerService(8)
	defer func(){
		o.Stop()
		o.Join()
	}()
	o.Set(key(""), func(arg interface{}) {
		v = arg.(string)
	})
	time.Sleep(time.Millisecond)

	arg := "11:56"
	o.Handle(key(arg), arg)
	time.Sleep(time.Millisecond)
	if v != arg {
		t.Errorf("v = %v, want %v", v, arg)
	}
}
