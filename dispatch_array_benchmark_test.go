package golib

import (
	"math"
	"testing"
)

func BenchmarkArrayDispatch_Call(b *testing.B) {
	o := NewArrayDispatch(math.MaxInt8)
	index := o.Add(func(interface{}) {})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		o.Call(index, nil)
	}
}

func BenchmarkArrayDispatch_UnsafeCall(b *testing.B) {
	o := NewArrayDispatch(math.MaxInt8)
	index := o.Add(func(interface{}) {})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		o.UnsafeCall(index, nil)
	}
}
