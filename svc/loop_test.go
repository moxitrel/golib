package svc

import (
	"testing"
	"time"
)

func TestLoop_Example(t *testing.T) {
	i := 0
	o := NewLoop(func() {
		i++
	})
	defer func() {
		o.Stop()
		o.Join()
	}()

	time.Sleep(time.Millisecond)
	if i == 0 {
		t.Errorf("i == 0, want !0")
	} else {
		t.Logf("i: %v", i)
	}
}
