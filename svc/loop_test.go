package svc

import (
	"math"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func TestLoop_Example(t *testing.T) {
	t.Logf("Loop.size: %v", unsafe.Sizeof(*NewLoop(func() {})))

	var n uint64 = 0
	var loopStartSignal = struct {
		sync.Once
		signal chan struct{}
	}{
		signal: make(chan struct{}),
	}

	o := NewLoop(func() {
		loopStartSignal.Do(func() {
			loopStartSignal.signal <- struct{}{}
		})

		if n < math.MaxUint64 {
			n++
		}
	})

	<-loopStartSignal.signal
	time.Sleep(time.Microsecond)
	o.Stop()
	o.Join()

	if n == 0 {
		t.Errorf("n == %v, want !0", n)
	} else {
		t.Logf("process count: %v", n)
	}
}

func TestLoop_DataRace(t *testing.T) {
	o := NewLoop(func() {})
	for i := 0; i < 3; i++ {
		NewLoop(func() {
			o.State()
		})
		NewLoop(func() {
			o.Stop()
		})
		NewLoop(func() {
			o.Join()
		})
	}
	o.Join()
}
