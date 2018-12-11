package benchmarks

import (
	"github.com/moxitrel/golib/svc"
	"reflect"
	"testing"
	"time"
)

func BenchmarkChan_SelectRecv0(b *testing.B) {
	c := make(chan interface{})
	svc.NewLoop(func() {
		c <- struct{}{}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-c
	}
}
func BenchmarkChan_SelectRecv1(b *testing.B) {
	c := make(chan interface{})
	svc.NewLoop(func() {
		c <- struct{}{}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case <-c:
		}
	}
}
func BenchmarkChan_SelectRecv1Default(b *testing.B) {
	c := make(chan interface{})
	svc.NewLoop(func() {
		c <- struct{}{}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case <-c:
		default:
			<-c
		}
	}
}
func BenchmarkChan_SelectRecv2(b *testing.B) {
	c := make(chan interface{})
	c2 := make(chan interface{})
	svc.NewLoop(func() {
		c <- struct{}{}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case <-c:
		case <-c2:
		}
	}
}
func BenchmarkChan_SelectRecv2Timer(b *testing.B) {
	c := make(chan interface{})
	svc.NewLoop(func() {
		c <- struct{}{}
	})

	// init timer
	delayTimer := svc.NewTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delayTimer.Start(time.Second)
		select {
		case <-c:
			delayTimer.Stop()
		case <-delayTimer.C:
		}
	}
}
func BenchmarkChan_SelectRecv3(b *testing.B) {
	c := make(chan interface{})
	c2 := make(chan interface{})
	c3 := make(chan interface{})
	svc.NewLoop(func() {
		c <- struct{}{}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case <-c:
		case <-c2:
		case <-c3:
		}
	}
}
func BenchmarkChan_SelectRecv4(b *testing.B) {
	c := make(chan interface{})
	c2 := make(chan interface{})
	c3 := make(chan interface{})
	c4 := make(chan interface{})
	svc.NewLoop(func() {
		c <- struct{}{}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case <-c:
		case <-c2:
		case <-c3:
		case <-c4:
		}
	}
}
func benchmarkChan_SelectRecv(b *testing.B, n int) {
	c := make(chan interface{})
	svc.NewLoop(func() {
		c <- struct{}{}
	})
	xs := make([]reflect.SelectCase, n)
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
func BenchmarkChan_SelectRecv100(b *testing.B) {
	benchmarkChan_SelectRecv(b, 100)
}
func BenchmarkChan_SelectRecv1000(b *testing.B) {
	benchmarkChan_SelectRecv(b, 1000)
}
func BenchmarkChan_SelectRecv10000(b *testing.B) {
	benchmarkChan_SelectRecv(b, 1000)
}
func BenchmarkChan_SelectRecv10000_0(b *testing.B) {
	c := make(chan interface{})
	svc.NewLoop(func() {
		c <- struct{}{}
	})
	xs := make([]reflect.SelectCase, 10000)
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
