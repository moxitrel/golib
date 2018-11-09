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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c <- nil
	}
}

//func BenchmarkChan_SendSelect1(b *testing.B) {
//	c := make(chan interface{})
//	NewLoop(func() {
//		<-c
//	})
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		select {
//		case c <- nil:
//		}
//	}
//}

func BenchmarkChan_SendSelect2(b *testing.B) {
	c := make(chan interface{})
	c2 := make(chan interface{})
	NewLoop(func() {
		<-c
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case c <- nil:
		case c2 <- nil:
		}
	}
}

//func BenchmarkChan_SendSelect2ForBreak(b *testing.B) {
//	c := make(chan interface{})
//	c2 := make(chan interface{})
//	NewLoop(func() {
//		<-c
//	})
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		for {
//			select {
//			case c <- nil:
//			case c2 <- nil:
//			}
//			break
//		}
//	}
//}

//func BenchmarkChan_SendSelect3(b *testing.B) {
//	c := make(chan interface{})
//	c2 := make(chan interface{})
//	c3 := make(chan interface{})
//	NewLoop(func() {
//		<-c
//	})
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		select {
//		case c <- nil:
//		case c2 <- nil:
//		case c3 <- nil:
//		}
//	}
//}

func BenchmarkChan_SendWithTimerAfter(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		<-c
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case c <- nil:
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func BenchmarkChan_SendWithTimerReset(b *testing.B) {
	c := make(chan interface{})
	NewLoop(func() {
		<-c
	})

	// init timer
	delayTimer := time.NewTimer(0)
	if !delayTimer.Stop() {
		<-delayTimer.C
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delayTimer.Reset(time.Second)
		select {
		case c <- nil:
			if !delayTimer.Stop() {
				<-delayTimer.C
			}
		case <-delayTimer.C:
		}
	}
}

func BenchmarkInitTimer_Recv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := time.NewTimer(0)
		<-t.C
	}
}

//func BenchmarkInitTimer_Stop0(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		t := time.NewTimer(0)
//		if !t.Stop() {
//			<-t.C
//		}
//	}
//}

func BenchmarkInitTimer_Stop(b *testing.B) {
	var t = 100 * time.Millisecond
	for i := 0; i < b.N; i++ {
		t := time.NewTimer(t)
		if !t.Stop() {
			<-t.C
		}
	}
}
