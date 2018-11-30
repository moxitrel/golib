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
	}
}
