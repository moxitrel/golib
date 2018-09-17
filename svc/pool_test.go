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
	ngoBegin := runtime.NumGoroutine()

	min := 2
	delay := time.Millisecond
	timeout := time.Second
	f := NewPool(func(_ interface{}) {
		time.Sleep(time.Second)
	})
	f.SetTime(delay, timeout)
	f.SetCount(uint(min), POOL_MAX)
	time.Sleep(timeout + 10*time.Millisecond)
	ngoNewPool := runtime.NumGoroutine()
	if ngoNewPool != ngoBegin+min {
		t.Errorf("Goroutine.Count: %v, want %v", ngoNewPool, ngoBegin+2)
	}

	// f has 90 goroutines, 2 old, 88 new
	nCall := 90
	for i := 0; i < nCall; i++ {
		f.Call(nil)
	}
	ngoCall := runtime.NumGoroutine()
	if ngoCall != ngoBegin+nCall {
		t.Errorf("Goroutine.Count: %v, want %v", ngoCall, ngoBegin+nCall)
	}

	for f.cur > int32(f.min) {
		time.Sleep(f.timeout)
	}
	ngoTimeout := runtime.NumGoroutine()
	if ngoTimeout != ngoNewPool {
		t.Errorf("Goroutine.Count: %v, want %v", ngoTimeout, ngoNewPool)
	}
}

func TestPool_Example(t *testing.T) {
	ts := make([]time.Time, 0, 100)
	delay := 10 * time.Millisecond
	timeout := (delay + 5*time.Millisecond) * time.Duration(cap(ts))

	f := NewPool(func(x interface{}) {
		ts = append(ts, time.Now())
		time.Sleep(timeout)
	})
	f.SetTime(delay, timeout)
	f.SetCount(1, POOL_MAX)
	time.Sleep(timeout)

	for i := 0; i < cap(ts); i++ {
		f.Call(nil)
	}

	for i := 0; i < len(ts)-1; i++ {
		dt := ts[i+1].Sub(ts[i])
		if dt < delay || dt > delay+100*time.Millisecond {
			t.Errorf("%v: dt = %v, want [%v, %v]", i, dt, delay, delay+10*time.Millisecond)
		}
	}
	f.SetCount(0, uint(f.max))
	for f.cur > 0 {
		time.Sleep(f.timeout / 2)
	}
}
