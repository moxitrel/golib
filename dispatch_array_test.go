package golib

import (
	"github.com/moxitrel/golib/svc"
	"math"
	"sync/atomic"
	"testing"
	"time"
)

func TestArrayDispatch_DataRace(t *testing.T) {
	o := NewArrayDispatch(math.MaxInt32)
	var last = o.Add(func(_ interface{}) {})

	svc.NewLoop(func() {
		i := o.Add(func(_ interface{}) {})
		atomic.StoreUintptr(&last, i)
	})
	svc.NewLoop(func() {
		o.UnsafeCall(atomic.LoadUintptr(&last), nil)
	})

	time.Sleep(time.Second)
}
