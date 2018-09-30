package svc

import (
	"reflect"
	"testing"
	"time"
)

func TestDispatch_Example(t *testing.T) {
	key := func(arg interface{}) reflect.Type {
		return reflect.TypeOf(arg)
	}
	v := ""

	o := NewDispatch(8, 0)
	defer func() {
		o.Stop()
		o.Join()
	}()
	o.Set(key(""), func(arg interface{}) {
		v = arg.(string)
	})

	arg := "11:56"
	o.Call(key(arg), arg)
	time.Sleep(o.delay + 100*time.Millisecond)
	if v != arg {
		t.Errorf("v = %v, want %v", v, arg)
	}
}
