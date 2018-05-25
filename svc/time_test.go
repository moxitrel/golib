package svc

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestLoop(t *testing.T) {
	var accuracy = time.Duration(rand.Intn(1000)) * time.Millisecond
	var intvl = accuracy
	var loopMax = 5

	o := NewTime(accuracy)
	defer o.Stop()

	i := 0
	o.Loop(accuracy, func() {
		i++
		fmt.Printf("%s\n", time.Now())
	})
	time.Sleep(time.Duration(loopMax) * intvl)
	if i != loopMax {
		t.Errorf("i = %d, want %d", i, loopMax)
	}
}

//func TestAt(t *testing.T) {
//	intvl := time.Second
//	o := NewTime(intvl)
//	o.Start()
//	defer o.Stop()
//
//	now1 := time.Now()
//	o.At(now1.Truncate(time.Second), func() {
//		now2 := time.Now()
//		if now2.Sub(now1) > intvl {
//			t.Errorf("intvl ",)
//		}
//	})
//	time.Sleep(3*intvl)
//}

func TestAtLoop(t *testing.T) {
	var accuracy = 100 * time.Millisecond
	var intvl = accuracy
	var loopMax = 5

	o := NewTime(accuracy)
	defer o.Stop()

	i := 0
	o.At(time.Now().Truncate(10*intvl).Add(5*intvl), func() {
		o.Loop(accuracy, func() {
			i++
			fmt.Printf("%s\n", time.Now())
		})
	})
	time.Sleep(time.Duration(loopMax) * intvl)
}
