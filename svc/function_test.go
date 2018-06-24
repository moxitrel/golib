package svc

import (
	"sync"
	"testing"
	"time"
)

func Test_StopSignal_Uniqueness(t *testing.T) {
	type MockStopSignal struct{}
	mockStopSignal := MockStopSignal{}
	structStopSignal := struct{}{}

	t.Logf("_STOP_SIGNAL: %#v", _STOP_SIGNAL)
	t.Logf("mockStopSignal: %#v", mockStopSignal)
	t.Logf("structStopSignal: %#v", structStopSignal)

	if interface{}(mockStopSignal) == interface{}(_STOP_SIGNAL) ||
		interface{}(structStopSignal) == interface{}(_STOP_SIGNAL) {
		t.Fatalf("_STOP_SIGNAL isn't unique.")
	}
}

func Test_Function(t *testing.T) {
	o := NewFunction(100, func(x interface{}) {
		t.Logf("%v", x)
	})
	defer o.Stop()
	o.Call(1)
	o.Call(2)
	o.Call(3)
	time.Sleep(time.Millisecond)
}

func Test_FunctionCallAfterStop(t *testing.T) {
	o := NewFunction(100, func(x interface{}) {
		t.Logf("%v", x)
	})
	o.Stop()

	// no panic
	// no effect
	o.Call(1)
}

func Test_FunctionNewWithNil(t *testing.T) {
	// no panic
	o := NewFunction(100, nil)
	defer o.Stop()

	// no panic
	// no effect
	o.Call(1)
}

func Test_FunctionStopCallRace(t *testing.T) {
	startSignal := make(chan struct{})
	startOnce := sync.Once{}
	o := NewFunction(DefaultBufferSize, func(x interface{}) {
		startOnce.Do(func() {
			startSignal <- struct{}{}
		})
		//t.Logf("%v", time.Now())
	})
	o.Call(nil)
	<-startSignal //ensure o is started

	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		for o.state == RUNNING {
			o.Call(0)
		}
		wg.Done()
	}()
	time.Sleep(10 * time.Millisecond)
	o.Stop()
	o.Join()
	wg.Wait()
	if len(o.args) != 0 {
		t.Errorf("args.len = %v, want 0", len(o.args))
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
//	fs := NewFunction(100, f)
//	ngo1 += 1 // coroutine created by NewFunction()
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
//	fs := NewFunction(math.MaxUint16, f)
//	defer fs.Join()
//	defer fs.Stop()
//	defer time.Sleep(time.Millisecond)
//
//	for i := 0; i < 100; i++ {
//		fs.Call(nil)
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
