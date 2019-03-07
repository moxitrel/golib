package gosvc

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

type T1 struct{}

func (T1) Type() interface{} {
	type T struct{}
	return T{}
}

type T2 struct{}

func (T2) Type() interface{} {
	type T struct{}
	return T{}
}

func TestFunc_StopSignal(t *testing.T) {
	funcStopSignal := _FuncStopSignal{}

	type MockStopSignal struct{}
	if (MockStopSignal{}) == interface{}(funcStopSignal) {
		t.Errorf("funcStopSignal isn't unique: type MockStopSignal struct{}")
	}
	if (struct{}{}) == interface{}(funcStopSignal) {
		t.Errorf("funcStopSignal isn't unique: struct{}{}")
	}
	if (T1{}.Type()) == funcStopSignal {
		t.Errorf("funcStopSignal isn't unique: T1{}.Type()")
	}
	if (T1{}.Type()) == (T2{}.Type()) {
		t.Errorf("(T1{}.Type()) == (T2{}.Type()), want !=")
	}
}

func TestFunc_Call(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var oX int32 = 0
	signal1 := make(chan struct{})
	signal2 := make(chan struct{})
	o := NewFunc(math.MaxUint16, func(arg interface{}) {
		signal1 <- struct{}{}
		oX = arg.(int32)
		signal2 <- struct{}{}
	})

	max := rand.Int31n(math.MaxInt16)
	xs := make(map[int32]struct{})
	for x := int32(0); x < max; x++ {
		xs[x] = struct{}{}
		o.Call(x)
	}
	for x := int32(0); x < max; x++ {
		<-signal1
		<-signal2
		if _, ok := xs[oX]; !ok {
			t.Fatalf("oX == %v not in xs", oX)
		}
		delete(xs, oX)
	}
}

func TestFunc_CallAfterStop(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	oX := 0
	o := NewFunc(uint(rand.Intn(math.MaxInt16)), func(arg interface{}) {
		oX = arg.(int)
	})
	o.Stop()
	o.Wait()

	// no effect after stop
	o.Call(2)
	time.Sleep(100 * time.Millisecond)
	if oX == 2 {
		t.Errorf("oX == %v, want != 2", oX)
	}
}

func TestFunc_DataRace(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	o := NewFunc(uint(rand.Intn(math.MaxUint16)), func(i interface{}) {})
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
			o.Wait()
		})
	}
	o.Wait()
}

//
// Benchmarks
//
func benchmarkFunc_Func(b *testing.B, n uint) {
	o := NewFunc(n, func(interface{}) {})

	for i := 0; i < b.N; i++ {
		o.Call(i)
	}
	o.Stop()
	o.Wait()
}

func benchmarkFunc_Select(b *testing.B, n uint) {
	f := func(interface{}){}
	args := make(chan int, n)
	timeout := make(chan struct{})

	NewLoop(func() {
		select {
		case arg := <- args:
			f(arg)
		case <-timeout:
		}
	})
	for i := 0; i < b.N; i++ {
		args <- i
	}
	for len(args) > 0 {
		time.Sleep(100 * time.Millisecond)
	}
}

func BenchmarkFunc_Func0(b *testing.B) {
	benchmarkFunc_Func(b, 0)
}
//func BenchmarkFunc_Func10(b *testing.B) {
//	benchmarkFunc_Func(b, 10)
//}
//func BenchmarkFunc_Func100(b *testing.B) {
//	benchmarkFunc_Func(b, 100)
//}
func BenchmarkFunc_Func1000(b *testing.B) {
	benchmarkFunc_Func(b, 1000)
}
func BenchmarkFunc_Func10000(b *testing.B) {
	benchmarkFunc_Func(b, 10000)
}
func BenchmarkFunc_Select0(b *testing.B) {
	benchmarkFunc_Select(b, 0)
}
//func BenchmarkFunc_Select10(b *testing.B) {
//	benchmarkFunc_Select(b, 10)
//}
//func BenchmarkFunc_Select100(b *testing.B) {
//	benchmarkFunc_Select(b, 100)
//}
func BenchmarkFunc_Select1000(b *testing.B) {
	benchmarkFunc_Select(b, 1000)
}
func BenchmarkFunc_Select10000(b *testing.B) {
	benchmarkFunc_Select(b, 10000)
}
