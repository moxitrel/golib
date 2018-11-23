package svc

import (
	"math"
	"reflect"
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

func BenchmarkNow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		time.Now()
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
		case c <- struct{}{}:
		default:
			c <- struct{}{}
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
		case c <- struct{}{}:
		case c2 <- struct{}{}:
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
		case c <- struct{}{}:
		case c2 <- struct{}{}:
		case c3 <- struct{}{}:
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
		case c <- struct{}{}:
		case c2 <- struct{}{}:
		case c3 <- struct{}{}:
		case c4 <- struct{}{}:
		}
	}
}
func BenchmarkChan_SelectRecv127(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		c <- nil
	})
	xs := make([]reflect.SelectCase, math.MaxInt8)
	for i, _ := range xs {
		xs[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(make(chan interface{})),
		}
	}
	xs[len(xs)-1] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(c),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reflect.Select(xs)
	}
}
func BenchmarkChan_SelectRecv255(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		c <- nil
	})
	xs := make([]reflect.SelectCase, math.MaxUint8)
	for i, _ := range xs {
		xs[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(make(chan interface{})),
		}
	}
	xs[len(xs)-1] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(c),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reflect.Select(xs)
	}
}
func BenchmarkChan_SelectRecv1000(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		c <- nil
	})
	xs := make([]reflect.SelectCase, 1000)
	for i, _ := range xs {
		xs[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(make(chan interface{})),
		}
	}
	xs[len(xs)-1] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(c),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reflect.Select(xs)
	}
}
func BenchmarkChan_SelectRecv10000(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		c <- nil
	})
	xs := make([]reflect.SelectCase, 10000)
	for i, _ := range xs {
		xs[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(make(chan interface{})),
		}
	}
	xs[len(xs)-1] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(c),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reflect.Select(xs)
	}
}
func BenchmarkChan_SelectRecv32768(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		c <- nil
	})
	xs := make([]reflect.SelectCase, math.MaxInt16)
	for i, _ := range xs {
		xs[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(make(chan interface{})),
		}
	}
	xs[len(xs)-1] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(c),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reflect.Select(xs)
	}
}
func BenchmarkChan_SelectRecv65535(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		c <- nil
	})
	xs := make([]reflect.SelectCase, math.MaxUint16)
	for i, _ := range xs {
		xs[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(make(chan interface{})),
		}
	}
	xs[len(xs)-1] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(c),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reflect.Select(xs)
	}
}
func BenchmarkChan_SelectRecv65535_0(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		c <- nil
	})
	xs := make([]reflect.SelectCase, math.MaxUint16)
	for i, _ := range xs {
		xs[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.Zero(reflect.TypeOf(make(chan interface{}))),
		}
	}
	xs[len(xs)-1] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(c),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reflect.Select(xs)
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
	o := NewPool(1, 1, _POOL_WORKER_DELAY, _POOL_WORKER_TIMEOUT, 0, func(interface{}) {})
	call := o.Call

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
	o := NewPool(1, 1, _POOL_WORKER_DELAY, _POOL_WORKER_TIMEOUT, 1, func(interface{}) {})
	call := o.Call

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
	o := NewPool(1, 1, _POOL_WORKER_DELAY, _POOL_WORKER_TIMEOUT, math.MaxUint8, func(interface{}) {})
	call := o.Call

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
	o := NewPool(1, 1, _POOL_WORKER_DELAY, _POOL_WORKER_TIMEOUT, math.MaxUint16, func(interface{}) {})
	call := o.Call

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
