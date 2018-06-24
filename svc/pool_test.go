package svc

import (
	"runtime"
	"testing"
	"time"
)

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
	time.Sleep(30 * time.Second + PoolTimeOut)
	ngoTimeout := runtime.NumGoroutine()
	if ngoTimeout != ngoNewPool {
		t.Errorf("Goroutine.Count: %v, want %v", ngoTimeout, ngoNewPool)
	}
}
