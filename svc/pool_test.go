package svc

import (
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func Test_Select(t *testing.T) {
	t.Skipf("skip test select")

	n := rand.Int()
	delay := 100 * time.Millisecond
	c := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		c <- struct{}{}
	}

	for i := 0; i < rand.Int(); i++ {
		select {
		case <-c:
		case <-time.After(delay):
			// If <delay> is too small, select may choose this case even <o.arg> isn't blocked.
			// May be interfered by the delay caused by gc.
			//
			// 100ms is a proper value in my test if channel buffer large data
			t.Errorf("%v: %v+%v: select time.After(), want <-c", delay, i, len(c))
		}
	}
}

func TestPool_NumGoroutine(t *testing.T) {
	ngoBegin := runtime.NumGoroutine()

	delay := time.Millisecond
	timeout := time.Second
	min := rand.Intn(POOL_MAX) + POOL_MIN
	if min > POOL_MAX {
		min = POOL_MAX
	}
	f := NewPool(func(_ interface{}) {
		time.Sleep(time.Second)
	}).
		SetTime(delay, timeout).
		SetCount(uint(min), POOL_MAX)
	time.Sleep(time.Millisecond) // wait goroutines to start completely

	ngoNewPool := runtime.NumGoroutine()
	if ngoNewPool != ngoBegin+min {
		t.Fatalf("Goroutine.Count: %v, want %v", ngoNewPool, ngoBegin+min)
	}

	nCall := int(rand.Intn(POOL_MAX))
	for i := 0; i < nCall; i++ {
		f.Call(nil)
	}
	ngoCall := runtime.NumGoroutine()
	wantNgo := ngoNewPool
	switch {
	case nCall > int(f.max):
		wantNgo = int(f.max)
	case nCall > ngoNewPool:
		wantNgo = nCall
	default:
		wantNgo = ngoNewPool
	}

	if ngoCall != wantNgo {
		t.Fatalf("Goroutine.Count: %v, want %v", ngoCall, wantNgo)
	}

	for f.cur > int32(f.min) {
		time.Sleep(f.timeout)
	}
	ngoTimeout := runtime.NumGoroutine()
	if ngoTimeout != ngoNewPool {
		t.Fatalf("Goroutine.Count: %v, want %v", ngoTimeout, ngoNewPool)
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
