package svc

import (
	"testing"
	"time"
)

func BenchmarkChan_Send(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		<-c
	})
	for i := 0; i < b.N; i++ {
		c <- nil
	}
}
func BenchmarkChan_SendWithTimeout(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		<-c
	})
	for i := 0; i < b.N; i++ {
		select {
		case c <- nil:
		case <-time.After(100 * time.Millisecond):
		}
	}
}
