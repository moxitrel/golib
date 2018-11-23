package svc

import "testing"

func BenchmarkFunc_0(b *testing.B) {
	o := NewFunc(0, func(x interface{}) {})
	NewLoop(func() {
		o.Call(struct{}{})
	})
	for i := 0; i < b.N; i++ {

	}
}
