package golib

import (
	"math"
	"math/rand"
	"reflect"
	"sync/atomic"
	"testing"
)

func BenchmarkFunCall_Reflect(b *testing.B) {
	var o = reflect.ValueOf(func(interface{}) {})

	for i := 0; i < b.N; i++ {
		o.Call([]reflect.Value{reflect.ValueOf(0)})
	}
}
func BenchmarkFunCall_AtomicValue(b *testing.B) {
	var o atomic.Value
	o.Store(func(interface{}) {})

	for i := 0; i < b.N; i++ {
		f := o.Load().(func(interface{}))
		f(0)
	}
}
func BenchmarkFunCall_Interface(b *testing.B) {
	var o interface{} = func(interface{}) {}

	for i := 0; i < b.N; i++ {
		f := o.(func(interface{}))
		f(0)
	}
}
func BenchmarkFunCall_Direct(b *testing.B) {
	var o = func(interface{}) {}

	for i := 0; i < b.N; i++ {
		o(0)
	}
}
func BenchmarkFunCall_ArrayDispatch(b *testing.B) {
	o := NewArrayDispatch(uintptr(1 + rand.Intn(math.MaxInt8)))
	n := o.Add(func(i interface{}) {})

	for i := 0; i < b.N; i++ {
		o.Call(n, 0)
	}
}
func BenchmarkFunCall_MapDispatch(b *testing.B) {
	o := NewMapDispatch()
	o.Set(0, func(i interface{}) {})

	for i := 0; i < b.N; i++ {
		o.Call(0, 0)
	}
}
