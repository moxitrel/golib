package svc

import (
	"testing"
	"reflect"
	"time"
)

func TestHandler_Register_Dup(t *testing.T) {
	h := NewHandler()
	flag := 0
	f1 := func(_ interface{}) { flag = 1 }
	f2 := func(_ interface{}) { flag = 2 }

	h.Set(1, f1)
	h.Set(1, f2)
	h.HandleWithoutCheckout(1, struct{}{})
	if flag != 2 {
		t.Errorf("flag = %v, want 2", flag)
	}
}

func TestHandlerService_Example(t *testing.T) {
	key := func(arg interface{}) reflect.Type {
		return reflect.TypeOf(arg)
	}
	v := ""

	o := NewHandlerService(8)
	o.Set(key(""), func(arg interface{}) {
		v = arg.(string)
	})
	time.Sleep(time.Millisecond)

	arg:= "11:56"
	o.Handle(key(arg), arg)
	time.Sleep(time.Millisecond)
	if v != arg {
		t.Errorf("v = %v, want %v", v, arg)
	}
}