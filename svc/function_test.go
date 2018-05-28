package svc

import (
	"testing"
	"time"
)

func Test_Function(t *testing.T) {
	o := NewFunction(func(x interface{}) {
		t.Logf("%v", x)
	})
	defer o.Stop()
	o.Call(1)
	o.Call(2)
	o.Call(3)
	time.Sleep(time.Millisecond)
}

func Test_FunctionCallAfterStop(t *testing.T) {
	o := NewFunction(func(x interface{}) {
		t.Logf("%v", x)
	})
	o.Stop()

	// no panic
	// no effect
	o.Call(1)
}

func Test_FunctionNewWithNil(t *testing.T) {
	// no panic
	o := NewFunction(nil)
	defer o.Stop()

	// no panic
	// no effect
	o.Call(1)
}

func Test_FunctionStopCallRace(t *testing.T) {
	o := NewFunction(func(x interface{}) {
		time.Sleep(1 * time.Second)
	})
	defer o.Join()
	defer o.Stop()
	go func() {
		for o.state == RUNNING {
			o.Call(0)
		}
	}()
	time.Sleep(1 * time.Second)
}