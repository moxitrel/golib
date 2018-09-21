package svc

import (
	"math"
	"math/rand"
	"sync"
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

func TestFunc_New(t *testing.T) {
	var x = 0
	signalBegin := make(chan struct{})
	signalEnd := make(chan struct{})
	o := NewFunc(math.MaxUint8, func(arg interface{}) {
		signalBegin <- struct{}{}
		x = arg.(int)
		signalEnd <- struct{}{}
	})
	defer o.Stop()

	o.Call(1)
	o.Call(2)
	o.Call(3)

	for _, v := range []int{1, 2, 3} {
		<-signalBegin
		<-signalEnd
		if x != v {
			t.Fatalf("x == %v, want %v", x, v)
		}
	}
}

func TestFunc_CallAfterStop(t *testing.T) {
	var x = 0
	o := NewFunc(math.MaxUint8, func(arg interface{}) {
		x = arg.(int)
	})
	o.Stop()
	o.Join()

	// no effect after stop
	o.Call(1)
	if x != 0 {
		t.Errorf("x = %v, want 0", x)
	}
}

func TestFunc_StopCallRace(t *testing.T) {
	var startSignal = struct {
		sync.Once
		signal chan struct{}
	}{
		signal: make(chan struct{}),
	}

	n := uint64(0)
	recver := NewFunc(uint(rand.Int()), func(interface{}) {
		startSignal.Do(func() {
			startSignal.signal <- struct{}{}
		})
		n++
	})
	sender := NewLoop(func() {
		recver.Call(0)
	})

	<-startSignal.signal
	time.Sleep(time.Millisecond)

	recver.Stop()
	sender.Stop()
	recver.Join()
	sender.Join()
	if len(recver.args) != 0 {
		t.Errorf("args.len = %v, want 0", len(recver.args))
	} else {
		t.Logf("process count: %v", n)
	}
}
