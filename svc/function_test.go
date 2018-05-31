package svc

import (
	"sync"
	"testing"
	"time"
	"runtime"
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
	o := NewFunction(FunctionBufferSize, func(x interface{}) {
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

	timeout := 100 * time.Millisecond
	f := func(x interface{}) {}
	var min uint16 = 3
	f = LimitWrap(f, &min, 100, 0, timeout)
	fs := NewFunction(100, f)
	ngo1 += 1	// master coroutine created by NewFunction()

	defer fs.Join()
	defer fs.Stop()
	defer time.Sleep(time.Millisecond)

	ngo2 := runtime.NumGoroutine()
	t.Logf("Goroutine.Count: %v", ngo2)
	if ngo1 + int(min) != ngo2 {
		t.Errorf("Goroutine.Count: %v, want %v", ngo2, ngo1 + int(min))
	}

	// 2.
	min = 0
	time.Sleep(2 * timeout)
	ngo2 = runtime.NumGoroutine()
	t.Logf("Goroutine.Count: %v", ngo2)
	if ngo1 != ngo2 {
		t.Errorf("Goroutine.Count: %v, want %v", ngo2, ngo1 + int(min))
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
		case <-time.After(delay):
			t.Fatalf("%v: %v+%v: select time.After(), want <-c", delay, i, len(c))
		}
	}
}
