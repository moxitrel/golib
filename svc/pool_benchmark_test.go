package svc

import (
	"testing"
	"time"
)

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
func benchmarkPoolPerf_Loop(b *testing.B, n int) {
	f := func(x interface{}) {}
	c := make(chan interface{}, n)
	NewLoop(func() {
		f(<-c)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c <- nil
	}
}
func benchmarkPoolPerf_Func(b *testing.B, n int) {
	o := NewFunc(uint(n), func(interface{}) {})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.Call(nil)
	}
}
func benchmarkPoolPerf_Pool(b *testing.B, n int) {
	o := NewPool(1, 1, 0, time.Minute, uint(n), func(interface{}) {})
	call := o.Call

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		call(nil)
	}
}

func BenchmarkPoolPerf_Pool_0(b *testing.B) {
	benchmarkPoolPerf_Pool(b, 0)
}
func BenchmarkPoolPerf_Func_0(b *testing.B) {
	benchmarkPoolPerf_Func(b, 0)
}
func BenchmarkPoolPerf_Loop_0(b *testing.B) {
	benchmarkPoolPerf_Loop(b, 0)
}
func BenchmarkPoolPerf_Pool_1(b *testing.B) {
	benchmarkPoolPerf_Pool(b, 1)
}
func BenchmarkPoolPerf_Func_1(b *testing.B) {
	benchmarkPoolPerf_Func(b, 1)

}
func BenchmarkPoolPerf_Loop_1(b *testing.B) {
	benchmarkPoolPerf_Loop(b, 1)
}
func BenchmarkPoolPerf_Pool_256(b *testing.B) {
	benchmarkPoolPerf_Pool(b, 256)
}
func BenchmarkPoolPerf_Func_256(b *testing.B) {
	benchmarkPoolPerf_Func(b, 256)

}
func BenchmarkPoolPerf_Loop_256(b *testing.B) {
	benchmarkPoolPerf_Loop(b, 256)
}
func BenchmarkPoolPerf_Pool_1024(b *testing.B) {
	benchmarkPoolPerf_Pool(b, 1024)
}
func BenchmarkPoolPerf_Func_1024(b *testing.B) {
	benchmarkPoolPerf_Func(b, 1024)

}
func BenchmarkPoolPerf_Loop_1024(b *testing.B) {
	benchmarkPoolPerf_Loop(b, 1024)
}
func BenchmarkPoolPerf_Pool_65536(b *testing.B) {
	benchmarkPoolPerf_Pool(b, 65536)
}
func BenchmarkPoolPerf_Func_65536(b *testing.B) {
	benchmarkPoolPerf_Func(b, 65536)
}
func BenchmarkPoolPerf_Loop_65536(b *testing.B) {
	benchmarkPoolPerf_Loop(b, 65536)

}
