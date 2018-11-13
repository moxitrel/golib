package svc

import (
	"math"
	"testing"
	"time"
)

func BenchmarkChan_SendTimer(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		<-c
	})

	// init timer
	delayTimer := NewTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delayTimer.Start(_STOP_DELAY)
		select {
		case c <- nil:
			delayTimer.Stop()
		case <-delayTimer.C:
		}
	}
}

func BenchmarkChan_SendAfter(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		<-c
	})

	// init timer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case c <- nil:
		case <-time.After(_STOP_DELAY):
		}
	}
}

func BenchmarkChan_Send(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		<-c
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c <- nil
	}
}

func BenchmarkChan_Select1(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		<-c
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case c <- nil:
		}
	}
}

func BenchmarkChan_Select2(b *testing.B) {
	c := make(chan interface{})
	c2 := make(chan interface{})
	NewLoop(func() {
		<-c
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case c <- nil:
		case c2 <- nil:
		}
	}
}

func BenchmarkChan_Select3(b *testing.B) {
	c := make(chan interface{})
	c2 := make(chan interface{})
	c3 := make(chan interface{})
	NewLoop(func() {
		<-c
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case c <- nil:
		case c2 <- nil:
		case c3 <- nil:
		}
	}
}

func BenchmarkChan_Select4(b *testing.B) {
	c := make(chan interface{})
	c2 := make(chan interface{})
	c3 := make(chan interface{})
	c4 := make(chan interface{})
	NewLoop(func() {
		<-c
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case c <- nil:
		case c2 <- nil:
		case c3 <- nil:
		case c4 <- nil:
		}
	}
}

func BenchmarkInitTimer_Recv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := time.NewTimer(0)
		<-t.C
	}
}
func BenchmarkInitTimer_Stop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := time.NewTimer(-1)
		if !t.Stop() {
			<-t.C
		}
	}
}

func BenchmarkPoolSubmit_Wait(b *testing.B) {
	o := Pool{
		max:   1,
		arg:   make(chan interface{}, 0),
		fun:   func(interface{}) {},
		delay: time.Second,
	}
	NewLoop(func() {
		<-o.arg
	})
	var arg interface{}
	delayTimer := NewTimer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delayTimer.Start(o.getDelay())
		select {
		case o.arg <- arg:
			delayTimer.Stop()
		case <-delayTimer.C:
		}
	}
}
func BenchmarkPoolSubmit_Select(b *testing.B) {
	o := Pool{
		max: 1,
		arg: make(chan interface{}, 0),
		fun: func(interface{}) {},
	}
	NewLoop(func() {
		<-o.arg
	})
	var arg interface{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case o.arg <- arg:
		default:
			select {
			case o.arg <- arg:
			}
		}
	}
}

func BenchmarkPoolPerf_Pool_0(b *testing.B) {
	o := NewPool(1, 1, _POOL_DELAY, _POOL_TIMEOUT, 0, func(interface{}) {})
	call := o.Submitter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		call(nil)
	}
}
func BenchmarkPoolPerf_Func_0(b *testing.B) {
	o := NewFunc(0, func(interface{}) {})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.Call(nil)
	}
}
func BenchmarkPoolPerf_Loop_0(b *testing.B) {
	f := func(x interface{}) {}
	c := make(chan interface{})
	NewLoop(func() {
		f(<-c)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c <- nil
	}
}
func BenchmarkPoolPerf_Pool_1(b *testing.B) {
	o := NewPool(1, 1, _POOL_DELAY, _POOL_TIMEOUT, 1, func(interface{}) {})
	call := o.Submitter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		call(nil)
	}
}
func BenchmarkPoolPerf_Func_1(b *testing.B) {
	o := NewFunc(1, func(interface{}) {})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.Call(nil)
	}
}
func BenchmarkPoolPerf_Loop_1(b *testing.B) {
	f := func(x interface{}) {}
	c := make(chan interface{}, 1)
	NewLoop(func() {
		f(<-c)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c <- nil
	}
}
func BenchmarkPoolPerf_Pool_8(b *testing.B) {
	o := NewPool(1, 1, _POOL_DELAY, _POOL_TIMEOUT, math.MaxUint8, func(interface{}) {})
	call := o.Submitter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		call(nil)
	}
}
func BenchmarkPoolPerf_Func_8(b *testing.B) {
	o := NewFunc(math.MaxUint8, func(interface{}) {})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.Call(nil)
	}
}
func BenchmarkPoolPerf_Loop_8(b *testing.B) {
	f := func(x interface{}) {}
	c := make(chan interface{}, math.MaxUint8)
	NewLoop(func() {
		f(<-c)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c <- nil
	}
}
func BenchmarkPoolPerf_Pool_16(b *testing.B) {
	o := NewPool(1, 1, _POOL_DELAY, _POOL_TIMEOUT, math.MaxUint16, func(interface{}) {})
	call := o.Submitter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		call(nil)
	}
}
func BenchmarkPoolPerf_Func_16(b *testing.B) {
	o := NewFunc(math.MaxUint16, func(interface{}) {})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.Call(nil)
	}
}
func BenchmarkPoolPerf_Loop_16(b *testing.B) {
	f := func(x interface{}) {}
	c := make(chan interface{}, math.MaxUint16)
	NewLoop(func() {
		f(<-c)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c <- nil
	}
}
