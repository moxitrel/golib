package svc

import (
	"math"
	"runtime"
	"sync"
	"testing"
	"time"
)

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
func Test_LimitWrapMin(t *testing.T) {
	ngo1 := runtime.NumGoroutine()

	f := func(x interface{}) {}
	var min uint16 = 3
	var max uint16 = 100
	var delay = time.Duration(0)
	var timeout = 100 * time.Millisecond
	f = LimitWrap(f, &min, &max, &delay, &timeout)
	fs := NewFunction(100, f)
	ngo1 += 1 // master coroutine created by NewFunction()

	defer fs.Join()
	defer fs.Stop()
	defer time.Sleep(time.Millisecond)

	ngo2 := runtime.NumGoroutine()
	t.Logf("Goroutine.Count: %v", ngo2)
	if ngo1+int(min) != ngo2 {
		t.Errorf("Goroutine.Count: %v, want %v", ngo2, ngo1+int(min))
	}

	// 2.
	min = 0
	time.Sleep(2 * timeout)
	ngo2 = runtime.NumGoroutine()
	t.Logf("Goroutine.Count: %v", ngo2)
	if ngo1 != ngo2 {
		t.Errorf("Goroutine.Count: %v, want %v", ngo2, ngo1+int(min))
	}
}

// 2. all created coroutine should quit if set min = 0
func Test_LimitWrapTimeout(t *testing.T) {
	f := func(x interface{}) {
		t.Logf("%v", time.Now())
		time.Sleep(100 * time.Millisecond)
	}

	var min uint16 = 1
	var max uint16 = 100
	var delay = 100 * time.Millisecond
	var timeout = delay
	f = LimitWrap(f, &min, &max, &delay, &timeout)
	fs := NewFunction(math.MaxUint16, f)
	defer fs.Join()
	defer fs.Stop()
	defer time.Sleep(time.Millisecond)

	for i := 0; i < 100; i++ {
		fs.Call(nil)
	}
}

func Test_Select(t *testing.T) {
	n := 10000 * 10000
	delay := time.Duration(0) //1 * time.Millisecond
	c := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		c <- struct{}{}
	}
	for i := 0; i < n; i++ {
		select {
		case <-c:
			// The case here is to ensure <c> is blocked
			//
			// Don't it seems doing the same thing as the case in default clause?
			// No, if <delay> is a small value, it would be interfered by gc.
		default:
			select {
			case <-c:
			case <-time.After(delay):
				t.Fatalf("%v: %v+%v: select time.After(), want <-c", delay, i, len(c))
			}
		}
	}
}
