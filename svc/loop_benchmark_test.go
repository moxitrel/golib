package svc

import (
	"sync/atomic"
	"testing"
)

func BenchmarkLoop_AtomicIf(b *testing.B) {
	thunk := func() {}
	o := &Loop{
		state: RUNNING,
	}
	const PAUSE = RUNNING + 1
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		state := atomic.LoadUintptr(&o.state)
		if state == RUNNING {
			thunk()
		} else if state == PAUSE {
			thunk()
		} else {
			thunk()
		}
	}
}

func BenchmarkLoop_AtomicSwitch(b *testing.B) {
	thunk := func() {}
	o := &Loop{
		state: RUNNING,
	}
	const PAUSE = RUNNING + 1
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch atomic.LoadUintptr(&o.state) {
		case RUNNING:
			thunk()
		case PAUSE:
			thunk()
		default:
			thunk()
		}
	}
}

func BenchmarkLoop_Atomic(b *testing.B) {
	thunk := func() {}
	o := &Loop{
		state: RUNNING,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if atomic.LoadUintptr(&o.state) == RUNNING {
			thunk()
		}
	}
}

func BenchmarkLoop_Raw(b *testing.B) {
	thunk := func() {}
	o := &Loop{
		state: RUNNING,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if o.state == RUNNING {
			thunk()
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
