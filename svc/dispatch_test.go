package svc

import (
	"testing"
)

type Msg struct {
	value int
}

func (Msg) DispatchKey() interface{} {
	return nil
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
	o.Set(Msg{}, func(DispatchMsg) {})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.Call(Msg{})
	}
}
func BenchmarkDispatch_Chan(b *testing.B) {
	c1 := make(chan interface{})
	NewLoop(func() {
		<-c1
	})
	for i := 0; i < b.N; i++ {
		c1 <- nil
	}
}
