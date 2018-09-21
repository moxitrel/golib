package svc

import (
	"math"
	"sync"
	"testing"
	"time"
)

func TestLoop_Example(t *testing.T) {
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
	defer func() {
		o.Stop()
		o.Join()
	}()

	<-loopStartSignal.signal
	time.Sleep(time.Microsecond)

	if n == 0 {
		t.Errorf("n == %v, want !0", n)
	} else {
		t.Logf("process count: %v", n)
	}
}
