package golib

import (
	"github.com/moxitrel/golib/svc"
	"math"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func TestArrayDispatch_DataRace(t *testing.T) {
	o := NewArrayDispatch(math.MaxUint32)
	last := o.Add(func(_ interface{}) {})

	for i := 0; i < 2; i++ {
		svc.NewLoop(func() {
			o.Add(func(_ interface{}) {})
			atomic.AddUintptr(&last, 1)
		})
		svc.NewLoop(func() {
			o.Call(uintptr(rand.Intn(int(atomic.LoadUintptr(&last)+1))), nil)
		})
	}
	time.Sleep(time.Second)
}
