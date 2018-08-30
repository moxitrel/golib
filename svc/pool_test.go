package svc

import (
	"runtime"
	"testing"
	"time"
)

func Test_Select(t *testing.T) {
	t.Skipf("skip test select")
	n := 10000 * 10000
	delay := 100 * time.Millisecond
	c := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		c <- struct{}{}
	}
	for i := 0; i < n; i++ {
		select {
		case <-c:
		case <-time.After(delay):
			// If <delay> is too small, select may choose this case even <o.arg> isn't blocked.
			// May be interfered by the delay caused by gc.
			//
			// 100ms is a proper value in my test if channel buffer large data
			t.Fatalf("%v: %v+%v: select time.After(), want <-c", delay, i, len(c))
		}
	}
}

func TestPool_NumGoroutine(t *testing.T) {
	PoolTimeOut = time.Second
	ngoBegin := runtime.NumGoroutine()

	// f generates 2 goroutines
	f := NewPool(func(x interface{}) {
		time.Sleep(30 * time.Second)
	})
	time.Sleep(time.Millisecond) //wait goroutine started
	ngoNewPool := runtime.NumGoroutine()
	if ngoNewPool != ngoBegin+2 {
		t.Errorf("Goroutine.Count: %v, want %v", ngoNewPool, ngoBegin+2)
	}

	// f has 90 goroutines, 2 old, 88 new
	nCall := 90
	for i := 0; i < nCall; i++ {
		f.Call(nil)
	}
	time.Sleep(time.Millisecond)
	ngoCall := runtime.NumGoroutine()
	if ngoCall != ngoBegin+nCall {
		t.Errorf("Goroutine.Count: %v, want %v", ngoCall, ngoBegin+nCall)
	}

	// f remains 2 goroutines after timeout
	time.Sleep(30*time.Second + PoolTimeOut)
	ngoTimeout := runtime.NumGoroutine()
	if ngoTimeout != ngoNewPool {
		t.Errorf("Goroutine.Count: %v, want %v", ngoTimeout, ngoNewPool)
	}
}

// 2. all created coroutine should quit if set min = 0
//func Test_LimitWrapMin(t *testing.T) {
//	ngo1 := runtime.NumGoroutine()
//
//	f := func(x interface{}) {}
//	var min uint = 7
//	var max uint = 100
//	var delay = time.Duration(0)
//	var timeout = 100 * time.Millisecond
//	f = PoolOf(f, &min, &max, &delay, &timeout)
//	fs := NewFuncService(100, f)
//	ngo1 += 1 // coroutine created by NewFuncService()
//
//	defer fs.Join()
//	defer fs.Stop()
//	defer time.Sleep(time.Millisecond)
//
//	ngo2 := runtime.NumGoroutine()
//	if ngo1+int(min) != ngo2 {
//		t.Errorf("Goroutine.Count: %v, want %v", ngo2, ngo1+int(min))
//	}
//
//	// 2.
//	min = 0
//	time.Sleep(2 * timeout)
//	ngo2 = runtime.NumGoroutine()
//	if ngo1 != ngo2 {
//		t.Errorf("Goroutine.Count: %v, want %v", ngo2, ngo1+int(min))
//	}
//}

// 2. all created coroutine should quit if set min = 0
//func Test_LimitWrapTimeout(t *testing.T) {
//	f := func(x interface{}) {
//		t.Logf("%v", time.Now())
//		time.Sleep(100 * time.Millisecond)
//	}
//
//	var min uint = 1
//	var max uint = 100
//	var delay = 100 * time.Millisecond
//	var timeout = delay
//	f = PoolOf(f, &min, &max, &delay, &timeout)
//	fs := NewFuncService(math.MaxUint16, f)
//	defer fs.Join()
//	defer fs.Stop()
//	defer time.Sleep(time.Millisecond)
//
//	for i := 0; i < 100; i++ {
//		fs.Apply(nil)
//	}
//}

//func TestPool(t *testing.T) {
//	f := func(x interface{}) {
//		if x == nil {
//			return
//		}
//		y := x.(func())
//		y()
//	}
//	var min uint = 1
//	var max uint = 1
//	var delay = 500 * time.Millisecond
//	var timeout = time.Minute
//	o := PoolOf(f, &min, &max, &delay, &timeout)
//
//	o(nil)
//	o(func() {
//		t1 := time.Now()
//		t2 := time.Now()
//		t.Logf("%v", t2.Sub(t1))
//	})
//
//	time.Sleep(time.Millisecond)
//}
