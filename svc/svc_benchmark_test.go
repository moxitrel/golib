package svc

import (
	"sync/atomic"
	"testing"
)

func BenchmarkSvc_FuncTest(b *testing.B) {
	o := &Svc{
		state: RUNNING,
	}
	do := func() {}

	for i := 0; i < b.N; i++ {
		switch o.State() {
		case RUNNING:
			do()
		default:
			goto BREAK_FOR
		}
	BREAK_FOR:
	}
}
func BenchmarkSvc_AtomicTest(b *testing.B) {
	o := &Svc{
		state: RUNNING,
	}
	do := func() {}

	for i := 0; i < b.N; i++ {
		if atomic.LoadInt32(&o.state) == RUNNING {
			do()
		}
	}
}
func BenchmarkSvc_RawTest(b *testing.B) {
	o := &Svc{
		state: RUNNING,
	}
	do := func() {}

	for i := 0; i < b.N; i++ {
		if o.state == RUNNING {
			do()
		}
	}
}
func BenchmarkSvc_NoTest(b *testing.B) {
	//o := &Svc{
	//	state: RUNNING,
	//}
	do := func() {}
	for i := 0; i < b.N; i++ {
		do()
	}
}
