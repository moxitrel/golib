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
func BenchmarkFunCall_MapDispatch(b *testing.B) {
	o := NewMapDispatch()
	key := o.NewKey()
	o.Set(key, func(interface{}) {})

	for i := 0; i < b.N; i++ {
		o.Call(key, 0)
	}
}
func BenchmarkFunCall_Map(b *testing.B) {
	n := 0
	o := map[interface{}]interface{}{
		n: func() {},
	}

	for i := 0; i < b.N; i++ {
		o[n].(func())()
	}
}
func BenchmarkFunCall_ArrayDispatch(b *testing.B) {
	o := NewArrayDispatch(uintptr(1 + rand.Intn(math.MaxInt8)))
	n := o.NewKey()
	o.Set(n, func(i interface{}) {})

	for i := 0; i < b.N; i++ {
		o.Call(n, 0)
	}
}
func BenchmarkFunCall_AtomicValue(b *testing.B) {
	o := make([]atomic.Value, 1)
	o[0].Store(func(interface{}) {})

	for i := 0; i < b.N; i++ {
		f := o[0].Load().(func(interface{}))
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
