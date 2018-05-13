package svc

import (
	"testing"
	"time"
)

func nop(_ interface{}) {
	//nop
}

type NopFun struct {
	Fun
}

func NewNopFunSvc() (v *NopFun) {
	v = &NopFun{
		Fun: *NewFun(func(argv []interface{}) {
			x := argv[0]
			nop(x)
		}),
	}
	return
}
func (o *NopFun) Call(x interface{}) {
	o.Fun.Call(x)
}
func (o *NopFun) Stop() {
	for len(o.argvs) != 0 {
		time.Sleep(time.Millisecond)
	}
	o.Fun.Stop()
}
func TestFun(t *testing.T) {
	o := NewNopFunSvc()
	for i := 0; i < cap(o.argvs); i++ {
		o.Call(struct {}{})
	}
	o.Start()
	defer o.Stop()
}
