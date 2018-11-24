package svc

import (
	"testing"
	"time"
)

func TestLoop_Example(t *testing.T) {
	var n uint64 = 0
	o := NewLoop(func() {
		n++
	})
	time.Sleep(time.Millisecond)
	o.Stop()
	o.Join()

	if n == 0 {
		t.Errorf("n == %v, want !0", n)
	} else {
		t.Logf("loop count: %v", n)
	}
}

func TestLoop_DataRace(t *testing.T) {
	o := NewLoop(func() {})
	for i := 0; i < 2; i++ {
		NewLoop(func() {
			o.State()
		})
	}
	for i := 0; i < 2; i++ {
		NewLoop(func() {
			o.Stop()
		})
		NewLoop(func() {
			o.Join()
		})
	}
	o.Join()
}
