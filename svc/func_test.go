package svc

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestFunc_StopSignal(t *testing.T) {
	type MockStopSignal struct{}
	mockStopSignal := MockStopSignal{}
	structStopSignal := struct{}{}

	t.Logf("stopSignal: %#v", stopSignal)
	t.Logf("mockStopSignal: %#v", mockStopSignal)
	t.Logf("structStopSignal: %#v", structStopSignal)

	if mockStopSignal == interface{}(stopSignal) ||
		structStopSignal == interface{}(stopSignal) {
		t.Errorf("stopSignal isn't unique.")
	}
}

func TestFunc_New(t *testing.T) {
	var x = 0
	signalBegin := make(chan struct{})
	signalEnd := make(chan struct{})
	o := NewFunc(math.MaxUint16, func(arg interface{}) {
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
	rand.Seed(time.Now().UnixNano())
	x := 0
	o := NewFunc(uint(rand.Intn(math.MaxInt16)), func(arg interface{}) {
		x = arg.(int)
	})
	o.Stop()
	o.Join()

	// no effect after stop
	o.Call(2)
	if x == 2 {
		t.Errorf("x == %v, want != 2", x)
	}
}

func TestFunc_Join(t *testing.T) {
	o := NewFunc(0, func(i interface{}) {})
	o.Stop()
	o.Join()
}

func TestFunc_DataRace(t *testing.T) {
	o := NewFunc(0, func(i interface{}) {})
	for i := 0; i < 2; i++ {
		NewLoop(func() {
			o.Call(nil)
		})
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
