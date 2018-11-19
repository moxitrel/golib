package svc

import (
	"sync/atomic"
	"testing"
)

func BenchmarkLoop_FuncTest(b *testing.B) {
	o := &Loop{
		state: RUNNING,
	}
	State := func() uintptr {
		return atomic.LoadUintptr(&o.state)
	}
	do := func() {}
	for i := 0; i < b.N; i++ {
		if State() == RUNNING {
			do()
		}
	}
}
func BenchmarkLoop_AtomicTest(b *testing.B) {
	o := &Loop{
		state: RUNNING,
	}
	do := func() {}
	for i := 0; i < b.N; i++ {
		if atomic.LoadUintptr(&o.state) == RUNNING {
			do()
		}
	}
}
func BenchmarkLoop_RawTest(b *testing.B) {
	o := &Loop{
		state: RUNNING,
	}
	do := func() {}
	for i := 0; i < b.N; i++ {
		if o.state == RUNNING {
			do()
		}
	}
}
func BenchmarkLoop_NoTest(b *testing.B) {
	do := func() {}
	for i := 0; i < b.N; i++ {
		do()
	}
}
