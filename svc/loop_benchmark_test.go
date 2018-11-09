package svc

import (
	"sync/atomic"
	"testing"
)

func BenchmarkLoop_AtomicIf(b *testing.B) {
	o := &Loop{
		thunk: func() {},
		state: RUNNING,
	}
	const PAUSE = RUNNING + 1
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		state := atomic.LoadUintptr(&o.state)
		if state == RUNNING {
			o.thunk()
		} else if state == PAUSE {
			o.thunk()
		} else {
			o.thunk()
		}
	}
}

func BenchmarkLoop_AtomicSwitch(b *testing.B) {
	o := &Loop{
		thunk: func() {},
		state: RUNNING,
	}
	const PAUSE = RUNNING + 1
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch atomic.LoadUintptr(&o.state) {
		case RUNNING:
			o.thunk()
		case PAUSE:
			o.thunk()
		default:
			o.thunk()
		}
	}
}

func BenchmarkLoop_Atomic(b *testing.B) {
	o := &Loop{
		thunk: func() {},
		state: RUNNING,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if atomic.LoadUintptr(&o.state) == RUNNING {
			o.thunk()
		}
	}
}

func BenchmarkLoop_Raw(b *testing.B) {
	o := &Loop{
		thunk: func() {},
		state: RUNNING,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if o.state == RUNNING {
			o.thunk()
		}
	}
}

func BenchmarkLoop_Direct(b *testing.B) {
	state := RUNNING
	for i := 0; i < b.N; i++ {
		if state == RUNNING {
			func() {}()
		}
	}
}

func BenchmarkLoop_NoTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {}()
	}
}
