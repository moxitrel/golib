package golib

import (
	"github.com/moxitrel/golib/svc"
	"sync/atomic"
	"testing"
	"time"
)

func TestMapDispatch_DataRace(t *testing.T) {
	o := NewMapDispatch()
	n := int64(0)
	for i := 0; i < 3; i++ {
		svc.NewLoop(func() {
			o.Add(atomic.AddInt64(&n, 1), func(interface{}) {})
		})
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
