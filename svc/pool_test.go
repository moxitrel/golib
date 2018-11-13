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

//func TestPool_NumGoroutine(t *testing.T) {
//	ngoBegin := runtime.NumGoroutine()
//
//	delay := time.Millisecond
//	timeout := time.Second
//	min := rand.Intn(_POOL_MAX) + _POOL_MIN
//	if min > _POOL_MAX {
//		min = _POOL_MAX
//	}
//	f := PoolWrapper(func(_ interface{}) {
//		time.Sleep(time.Second)
//	}).
//		WithTime(delay, timeout).
//		WithCount(uint(min), _POOL_MAX)
//	time.Sleep(time.Millisecond) // wait goroutines to start completely
//
//	ngoNewPool := runtime.NumGoroutine()
//	if ngoNewPool != ngoBegin+min {
//		t.Fatalf("Goroutine.Count: %v, want %v", ngoNewPool, ngoBegin+min)
//	}
//
//	nCall := int(rand.Intn(_POOL_MAX))
//	for i := 0; i < nCall; i++ {
//		f.Call(nil)
//	}
//	ngoCall := runtime.NumGoroutine()
//	wantNgo := ngoNewPool
//	if nCall > ngoNewPool {
//		wantNgo = nCall
//	}
//	if wantNgo > int(f.max) {
//		wantNgo = int(f.max)
//	}
//
//	if ngoCall != wantNgo {
//		t.Fatalf("Goroutine.Count: %v, want %v", ngoCall, wantNgo)
//	}
//
//	for f.cur > int32(f.min) {
//		time.Sleep(f.timeout)
//	}
//	ngoTimeout := runtime.NumGoroutine()
//	if ngoTimeout != ngoNewPool {
//		t.Fatalf("Goroutine.Count: %v, want %v", ngoTimeout, ngoNewPool)
//	}
//}

func TestPool_Example(t *testing.T) {
	ts := make([]time.Time, 0, 100)
	delay := 10 * time.Millisecond
	timeout := (delay + 5*time.Millisecond) * time.Duration(cap(ts))

	f := NewPool(1, _POOL_MAX, delay, timeout, 0, func(x interface{}) {
		ts = append(ts, time.Now())
		time.Sleep(timeout)
	})
	time.Sleep(timeout)

	for i := 0; i < cap(ts); i++ {
		f.Submitter()(nil)
	}

	for i := 0; i < len(ts)-1; i++ {
		dt := ts[i+1].Sub(ts[i])
		if dt < delay || dt > delay+100*time.Millisecond {
			t.Errorf("%v: dt = %v, want [%v, %v]", i, dt, delay, delay+10*time.Millisecond)
		}
	}
	f.Stop()
	for f.cur > 0 {
		time.Sleep(f.timeout / 2)
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
		}(o.Submitter())
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
	timeout := time.Second + time.Duration(rand.Int31())
	o := NewPool(uint(rand.Intn(math.MaxInt8)), math.MaxInt8, -1, timeout, 0, func(i interface{}) {})
	t1 := time.Now()
	o.Stop()
	o.Join()
	t2 := time.Now()
	t.Logf("t1: %v", t1)
	t.Logf("timeout: %v", timeout)
	t.Logf("t2: %v", t2)
	if t2.Sub(t1.Add(timeout)) > time.Second {
		t.Errorf("pool should be stop in %v", t2.Add(time.Second).Sub(t1.Add(timeout)))
	}
}
