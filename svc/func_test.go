package svc

import (
	"math"
	"testing"
)

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
	x := 0
	o := NewFunc(math.MaxUint16, func(arg interface{}) {
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

func TestFunc_Join(t *testing.T) {
	o := NewFunc(0, func(i interface{}) {})
	o.Stop()
	o.Join()
}

func TestFunc_DataRace(t *testing.T) {
	o := NewFunc(0, func(i interface{}) {})
	for i := 0; i < 3; i++ {
		NewLoop(func() {
			o.Call(nil)
		})
		NewLoop(func() {
			o.Stop()
		})
		NewLoop(func() {
			o.State()
		})
		NewLoop(func() {
			o.Join()
		})
	}
	o.Join()
}
