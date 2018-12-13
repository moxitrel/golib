package golib

import (
	"github.com/moxitrel/golib/svc"
	"testing"
	"time"
	"unsafe"
)

func TestMapDispatcher_DataRace(t *testing.T) {
	t.Logf("uintptr.size: %v", unsafe.Sizeof(uintptr(0)))

	o := new(MapDispatcher)
	o.Set(0, func(interface{}) {})
	for i := 0; i < 2; i++ {
		svc.NewLoop(func() {
			o.Set(0, func(interface{}) {})
		})
		svc.NewLoop(func() {
			o.Get(0)
		})
		svc.NewLoop(func() {
			f := o.Get(0)
			f(nil)
		})
	}
	time.Sleep(time.Second)
}
