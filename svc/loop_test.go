package svc

import (
	"sync"
	"testing"
	"time"
)

func Test_NewLoopWithNil(t *testing.T) {
	o := NewLoop(nil)
	o.Stop()
}

func Test_Loop(t *testing.T) {
	i := 0
	o := NewLoop(func() {
		i++
	})
	defer o.Stop()

	time.Sleep(time.Millisecond)
	if i == 0 {
		t.Errorf("i == 0, want !0")
	} else {
		t.Logf("i: %v", i)
	}
}

func Test_LoopJoin(t *testing.T) {
	startSignal := make(chan struct{})
	startOnce := sync.Once{}
	i := 0
	o := NewLoop(func() {
		startOnce.Do(func() {
			startSignal <- struct{}{}
		})
		time.Sleep(100 * time.Millisecond)
		i = 1
	})
	<-startSignal
	t.Logf("wg: %v", o.wg)

	o.Stop()
	o.Join()

	if i != 1 {
		t.Errorf("i = %v, want 1", i)
	}
}
