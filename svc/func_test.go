package svc

import (
	"math"
	"testing"
	"time"
)

func Test_StopSignal(t *testing.T) {
	type MockStopSignal struct{}
	mockStopSignal := MockStopSignal{}
	structStopSignal := struct{}{}

	t.Logf("_STOP_SIGNAL: %#v", _StopSignal{})
	t.Logf("mockStopSignal: %#v", mockStopSignal)
	t.Logf("structStopSignal: %#v", structStopSignal)

	if interface{}(mockStopSignal) == interface{}((_StopSignal{})) ||
		interface{}(structStopSignal) == interface{}((_StopSignal{})) {
		t.Fatalf("_STOP_SIGNAL isn't unique.")
	}
}

func TestFunc_NewWithNil(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("NewFunc(_, nil) should be panic.")
		}
	}()

	NewFunc(math.MaxUint8, nil)
}

func TestFunc_New(t *testing.T) {
	o := NewFunc(math.MaxUint8, func(x interface{}) {
		t.Logf("%v", x)
	})
	defer o.Stop()
	o.Call(1)
	o.Call(2)
	o.Call(3)
	time.Sleep(time.Millisecond)
}

func TestFunc_CallAfterStop(t *testing.T) {
	x := 0
	o := NewFunc(math.MaxUint8, func(arg interface{}) {
		x = arg.(int)
	})
	o.Stop()
	time.Sleep(time.Millisecond)

	// no effect after stop
	o.Call(1)
	if x != 0 {
		t.Errorf("x = %v, want %v", x, 0)
	}
}

func TestFunc_StopCallRace(t *testing.T) {
	o := NewFunc(math.MaxUint16, func(interface{}) {})
	time.Sleep(time.Millisecond)

	oCall := NewLoop(func() {
		o.Call(0)
	})
	time.Sleep(10 * time.Millisecond)
	o.Stop()
	oCall.Stop()
	o.Join()
	oCall.Join()
	if len(o.args) != 0 {
		t.Errorf("args.len = %v, want 0", len(o.args))
	}
}
