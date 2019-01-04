package golib

import (
	"math"
	"math/rand"
	"reflect"
	"sync/atomic"
	"testing"
)

type DispatchKeyTypedef struct{}

func (DispatchKeyTypedef) DispatchKey() interface{} {
	type T struct{}
	return T{}
}

type DispatchKeyReflect struct{}

func (o DispatchKeyReflect) DispatchKey() interface{} {
	return reflect.TypeOf(o)
}
func BenchmarkDispatchKey_Typedef(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DispatchKeyTypedef{}.DispatchKey()
	}
}
func BenchmarkDispatchKey_ReflectTypeOf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DispatchKeyReflect{}.DispatchKey()
	}
}

func BenchmarkFunCall_Reflect(b *testing.B) {
	var o = reflect.ValueOf(func(interface{}) {})

	for i := 0; i < b.N; i++ {
		o.Call([]reflect.Value{reflect.ValueOf(0)})
	}
}
func BenchmarkFunCall_MapDispatcher(b *testing.B) {
	o := new(MapDispatcher)
	o.Set(9, func(interface{}) {})

	for i := 0; i < b.N; i++ {
		f := o.Get(9)
		f(0)
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
func BenchmarkFunCall_ArrayDispatcher(b *testing.B) {
	o := NewArrayDispatcher(uintptr(1 + rand.Intn(math.MaxInt8)))
	n := o.Add(func(i interface{}) {})

	for i := 0; i < b.N; i++ {
		f := o.Get(n)
		f(0)
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
