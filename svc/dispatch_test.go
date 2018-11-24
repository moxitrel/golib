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

type SetMsg struct {
	unhashable []byte
}

func (o SetMsg) DispatchKey() interface{} {
	return nil
}
func TestDispatch_Set(t *testing.T) {
	o := NewDispatch(8, 1)
	defer func() {
		o.Stop()
		o.Join()
	}()
	o.Set(SetMsg{nil}, func(arg DispatchMsg) {})
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
