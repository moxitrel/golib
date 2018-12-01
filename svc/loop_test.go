package svc

import (
	"sync"
	"testing"
)

func TestLoop_Example(t *testing.T) {
	var n uint64 = 0
	var signalOnce sync.Once
	var startSignal = make(chan struct{})
	o := NewLoop(func() {
		n++
		signalOnce.Do(func() {
			startSignal <- struct{}{}
		})
	})
	<-startSignal
	o.Stop()
	o.Wait()

	if n == 0 {
		t.Errorf("n == %v, want !0", n)
	}
}
