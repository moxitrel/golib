package svc

import (
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
	o := NewFunction(FunctionBufferSize, func(x interface{}) {
		startOnce.Do(func() {
			startSignal <- struct{}{}
		})
		//t.Logf("%v", time.Now())
	})
	o.Call(nil)
	<-startSignal

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
