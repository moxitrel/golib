package golib

import (
	"github.com/moxitrel/golib/svc"
	"testing"
	"time"
	"unsafe"
)

func TestMapDispatch_DataRace(t *testing.T) {
	t.Logf("uintptr.size: %v", unsafe.Sizeof(uintptr(0)))

	o := NewMapDispatch()
	for i := 0; i < 2; i++ {
		svc.NewLoop(func() {
			o.Set(0, func(interface{}) {})
		})
		svc.NewLoop(func() {
			o.Get(0)
		})
		svc.NewLoop(func() {
			o.Call(0, nil)
		})
	}
	time.Sleep(time.Second)
}
