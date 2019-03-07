package gosvc

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

func Test_NestedSelect(t *testing.T) {
	signal := make(chan struct{})
	go func() {
		signal <- struct{}{}
	}()
	time.Sleep(time.Millisecond)

	n := 0
	timeout := time.Second
	timeAfter := func() <-chan time.Time {
		n += 1
		return time.After(timeout)
	}

	flag := 0
	select {
	case <-signal:
		flag = 1
	default:
		select {
		case <-signal:
			flag = 2
		case <-timeAfter(): // quit if idle for <idle> ns
			flag = 3
		}
	}

	if flag != 1 {
		t.Errorf("flag = %v, want 1", flag)
	}
	if n > 0 {
		t.Errorf("n = %v, want 0", n)
	}
}

func TestPool_Stop(t *testing.T) {
	// goroutine number at start
	ngo0 := runtime.NumGoroutine()

	// start pool
	rand.Seed(time.Now().UnixNano())
	min := rand.Intn(8000)
	timeout := time.Second + time.Duration(rand.Int31())
	o := NewWorkerPool(uint(min), uint(min), timeout, func(interface{}) {})
	t.Logf("ngo / min: %v / %v", runtime.NumGoroutine()-ngo0, min)

	// stop pool
	t1 := time.Now()
	o.Stop()
	o.Wait()
	t2 := time.Now()
	t.Logf("join / idle : %v / %v", t2.Sub(t1), timeout)
	//time.Sleep(_POOL_STOP_DELAY)

	// check the number of goroutine
	if d := runtime.NumGoroutine() - ngo0; d > 0 {
		t.Errorf("%v goroutines left after .Wait(), want 0", d)
	}
	// check the cost time to stop
	if d := t2.Sub(t1); d > timeout+time.Duration(min)*4*time.Millisecond {
		t.Errorf("pool should be stop in %v", t2.Add(time.Second).Sub(t1.Add(timeout)))
	}
}

func TestPool_Wait(t *testing.T) {
	o := NewWorkerPool(1, 1<<20, time.Second, func(interface{}) {})
	c := make(chan struct{})
	NewSvc(func() {
		c <- struct{}{}
	}, nil, func() {
		o.Submit(nil)
	})
	<-c
	time.Sleep(100 * time.Millisecond)

	o.Stop()
	o.Wait()
}
