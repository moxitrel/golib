package svc

import (
	"testing"
	"time"
)

func TestFun_Example(t *testing.T) {
	o := NewFun(func(x interface{}) {
		t.Logf("%v", x)
	})
	defer o.Stop()
	o.Call(1)
	o.Call(2)
	o.Call(3)
	time.Sleep(time.Millisecond)
}
