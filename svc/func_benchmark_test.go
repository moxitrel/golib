package svc

import (
	"math"
	"testing"
)

func BenchmarkFunc_SelectTest(b *testing.B) {
	o := Func{
		fun:  func(interface{}) {},
		args: make(chan interface{}, math.MaxInt8),
	}
	o.args <- nil
	var arg interface{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		o.args <- nil

		select {
		case arg = <-o.args:
		default:
		}

		if arg != (_StopSignal{}) {
			o.fun(arg)
		}
	}
}

//func BenchmarkFunc_Select(b *testing.B) {
//	o := Func{
//		fun:  func(interface{}) {},
//		args: make(chan interface{}, math.MaxInt8),
//	}
//	o.args <- nil
//	var arg interface{}
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		o.args <- nil
//
//		select {
//		case arg = <-o.args:
//		default:
//		}
//
//		o.fun(arg)
//	}
//}
//
//func BenchmarkFunc_Test(b *testing.B) {
//	o := Func{
//		fun:  func(interface{}) {},
//		args: make(chan interface{}, math.MaxInt8),
//	}
//	o.args <- nil
//	var arg interface{}
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		o.args <- nil
//
//		arg = <-o.args
//
//		if arg != (_StopSignal{}) {
//			o.fun(arg)
//		}
//	}
//}

func BenchmarkFunc_Raw(b *testing.B) {
	o := Func{
		fun:  func(interface{}) {},
		args: make(chan interface{}, math.MaxInt8),
	}
	o.args <- nil
	var arg interface{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		o.args <- nil

		arg = <-o.args

		o.fun(arg)
	}
}

func BenchmarkFunc_Direct(b *testing.B) {
	args := make(chan interface{}, math.MaxInt8)
	args <- nil
	var arg interface{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		args <- nil

		arg = <-args

		func(interface{}) {}(arg)
	}
}
