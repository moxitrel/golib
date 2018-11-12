package svc

import (
	"math"
	"math/rand"
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
	NewLoop(func() {
		o.State()
	})
	for i := 0; i < rand.Intn(math.MaxUint8); i++ {
		NewLoop(func() {
			o.Stop()
		})
	}
	o.Join()
}
