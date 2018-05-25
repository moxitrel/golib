package svc

import (
	"testing"
	"time"
)

func Test_NewWithNil(t *testing.T) {
	o := New(nil)
	o.Stop()
}

func TestService_New(t *testing.T) {
	i := 0
	o := New(func() {
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

func TestService_Stop(t *testing.T) {
	o := New(func() {
		time.Sleep(time.Millisecond)
	})
	o.Stop()

	if o.state != STOPPED {
		t.Errorf("o.state != STOPPED, want STOPPED")
	}
}
