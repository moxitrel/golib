package golib

import (
	"testing"
	"unsafe"
)

//import (
//	"github.com/moxitrel/golib/svc"
//	"math"
//	"math/rand"
//	"sync/atomic"
//	"testing"
//	"time"
//)

//func TestArrayDispatch_DataRace(t *testing.T) {
//	o := NewArrayDispatch(math.MaxUint32)
//	last := o.Add(func(_ interface{}) {})
//
//	for i := 0; i < 2; i++ {
//		svc.NewLoop(func() {
//			atomic.StoreUintptr(&last, o.Add(func(_ interface{}) {}))
//		})
//		svc.NewLoop(func() {
//			o.Call(uintptr(rand.Intn(int(atomic.LoadUintptr(&last)+1))), nil)
//		})
//	}
//	time.Sleep(time.Second)
//}

func TestArrayDispatch(t *testing.T) {
	t.Logf("sizeof ArrayDispatchKey: %v", unsafe.Sizeof(DispatchKey{}))
	t.Logf("sizeof uintptr: %v", unsafe.Sizeof(uintptr(0)))
}
