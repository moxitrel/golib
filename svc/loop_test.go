package svc

import (
	"math"
	"sync"
	"testing"
	"time"
)

func TestLoop_Example(t *testing.T) {
	var i uint64 = 0
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

		if i < math.MaxUint64 {
			i++
		}
	})
	defer func() {
		o.Stop()
		o.Join()
	}()

	<-loopStartSignal.signal
	time.Sleep(time.Microsecond)

	if i == 0 {
		t.Errorf("i == %v, want !0", i)
	} else {
		t.Logf("i == %v", i)
	}
}
