package svc

import (
	"testing"
)

type Msg struct {
	value int
}

var p = new(int)

func (o Msg) DispatchKey() interface{} {
	return p
}

func TestDispatch_Example(t *testing.T) {
	//v := ""
	//
	//o := NewDispatch(8, 1)
	//defer func() {
	//	o.Stop()
	//	o.Join()
	//}()
	//o.Set(Msg{34}, func(arg interface{}) {
	//	v = fmt.Sprintf("%v", arg.(Msg).value)
	//})
	//
	//arg := "11:56"
	//o.Call(Msg{78})
	//time.Sleep(o.delay + 100*time.Millisecond)
	//if v != arg {
	//	t.Errorf("v = %v, want %v", v, arg)
	//}
}

func BenchmarkDispatch_Call(b *testing.B) {
	o := NewDispatch(0, 1)
	o.Set(p, func(interface{}) {})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.Call(Msg{1})
	}
}
func BenchmarkDispatch_Chan(b *testing.B) {
	c1 := make(chan interface{})
	c2 := make(chan interface{})
	NewLoop(func() {
		<-c1
	})
	NewLoop(func() {
		c1 <- (<-c2)
	})
	for i := 0; i < b.N; i++ {
		c2 <- Msg{1}
	}
}
