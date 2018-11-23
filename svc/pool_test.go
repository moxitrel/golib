package svc

import (
	"math"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func Test_StopSignal(t *testing.T) {
	type MockStopSignal struct{}
	mockStopSignal := MockStopSignal{}
	structStopSignal := struct{}{}

	t.Logf("stopSignal: %#v", stopSignal)
	t.Logf("mockStopSignal: %#v", mockStopSignal)
	t.Logf("structStopSignal: %#v", structStopSignal)

	if mockStopSignal == interface{}(stopSignal) ||
		structStopSignal == interface{}(stopSignal) {
		t.Errorf("stopSignal isn't unique.")
	}
}

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
		case <-timeAfter(): // quit if idle for <timeout> ns
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

func TestPool_Example(t *testing.T) {
	ts := make([]time.Time, 0, 100)
	delay := 10 * time.Millisecond
	timeout := (delay + 5*time.Millisecond) * time.Duration(cap(ts))

	f := NewPool(1, _POOL_WORKER_MAX, delay, timeout, 0, func(x interface{}) {
		ts = append(ts, time.Now())
		time.Sleep(timeout)
	})
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
	f.Stop()
	for f.cur > 0 {
		time.Sleep(timeout / 2)
	}
}

func TestPool_DataRace(t *testing.T) {
	rand.Seed(int64(time.Now().Nanosecond()))

	min := rand.Intn(math.MaxInt8)
	max := min + rand.Intn(math.MaxInt8)
	delay := time.Duration(int64(_STOP_DELAY) + rand.Int63n(math.MaxInt32))
	timeout := time.Duration(int64(_STOP_DELAY) + rand.Int63n(math.MaxInt32))
	t.Logf("min: %v", min)
	t.Logf("max: %v", max)
	t.Logf("delay: %v", delay)
	t.Logf("timeout: %v", timeout)

	nBegin := runtime.NumGoroutine()

	o := NewPool(uint(min), uint(max), time.Duration(delay), time.Duration(timeout), 0, func(interface{}) {})
	time.Sleep(time.Millisecond)
	delta := runtime.NumGoroutine() - nBegin
	if delta != min {
		t.Errorf("ngo = %v, want %v", delta, min)
	}
	t.Logf("ngo.begin:%v = min:%v", delta, min)

	nCall := rand.Intn(math.MaxInt8)
	for i := 0; i < nCall; i++ {
		func(call func(interface{})) {
			NewLoop(func() {
				call(nil)
			})
		}(o.Call)
	}
	time.Sleep(delay + _STOP_DELAY)
	delta = runtime.NumGoroutine() - nBegin
	if delta != min+nCall {
		t.Errorf("ngo = %v, want %v", delta, min+nCall)
	}
	t.Logf("ngo.afterDelay:%v = min:%v + nCall:%v", delta, min, nCall)

	nStop := rand.Intn(math.MaxInt8)
	for i := 0; i < nStop; i++ {
		NewLoop(func() {
			o.Stop()
		})
	}
	time.Sleep(timeout + _STOP_DELAY)
	delta = runtime.NumGoroutine() - (nBegin + nCall + nStop)
	if delta != 0 {
		t.Errorf("ngo = %v, want %v", delta, 0)
	}
	t.Logf("ngo.stop: %v", delta)
}

func TestPool_Join(t *testing.T) {
	// goroutine number at start
	ngo0 := runtime.NumGoroutine()

	// start pool
	rand.Seed(time.Now().UnixNano())
	poolMin := rand.Intn(8000)
	poolTimeout := time.Duration(rand.Int31())
	if poolTimeout < time.Second {
		poolTimeout += time.Second
	}
	o := NewPool(uint(poolMin), uint(poolMin), _POOL_CALL_DELAY, poolTimeout, 0, func(interface{}) {})
	t.Logf("ngo: %v", runtime.NumGoroutine()-ngo0)

	// stop pool
	t1 := time.Now()
	o.Stop()
	o.Join()
	t2 := time.Now()
	t.Logf("join / timeout : %v / %v", t2.Sub(t1), poolTimeout)
	time.Sleep(_STOP_DELAY)

	// check
	if d := runtime.NumGoroutine() - ngo0; d > 0 {
		t.Errorf("%v goroutines left after .Join(), want 0", d)
	}
	if d := t2.Sub(t1); d > poolTimeout+time.Duration(poolMin)*time.Millisecond {
		t.Errorf("pool should be stop in %v", t2.Add(time.Second).Sub(t1.Add(poolTimeout)))
	}
}

func TestPool_Call_NoDelay(t *testing.T) {
	t.Skipf("need fix: .Call() may start more than 1 worker")
	// goroutine number at start
	ngo0 := runtime.NumGoroutine()

	o := NewPool(0, math.MaxUint8, 0, time.Hour, 0, func(interface{}) {})
	o.Call(nil)

	// check
	if d := runtime.NumGoroutine() - ngo0; d > 1 {
		t.Errorf("%v workers started, want 1", d)
	}
}
