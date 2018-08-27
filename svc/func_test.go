package svc

import (
	"testing"
	"time"
	"math"
)

func Test_StopSignal_Uniqueness(t *testing.T) {
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

func TestNewFuncWithNil(t *testing.T) {
	// no panic
	o := NewFuncService(math.MaxUint8, nil)
	defer o.Stop()

	// no panic
	// no effect
	o.Call(1)
	o.Call(nil)
	o.Call(struct{}{})
}

func TestNewFunc(t *testing.T) {
	o := NewFuncService(math.MaxUint8, func(x interface{}) {
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
	o := NewFuncService(math.MaxUint8, func(arg interface{}) {
		x = arg.(int)
	})
	o.Stop()

	// no effect after stop
	o.Call(1)
	time.Sleep(time.Millisecond)
	if x != 0 {
		t.Errorf("x = %v, want %v", x, 0)
	}
}

func TestFunc_StopCallRace(t *testing.T) {
	o := NewFuncService(math.MaxUint8, func(interface{}) {})
	time.Sleep(time.Millisecond)

	oCall := NewLoopService(func() {
		o.Call(0)
	})
	time.Sleep(10 * time.Millisecond)
	o.Stop()
	o.Join()
	oCall.Stop()
	oCall.Join()
	if len(o.args) != 0 {
		t.Errorf("args.len = %v, want 0", len(o.args))
	}
}
