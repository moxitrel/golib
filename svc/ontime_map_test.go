package svc

import (
	"math"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func TestTime_Every(t *testing.T) {
	var accuracy = 250 * time.Millisecond
	o := NewMapOnTime(accuracy)
	defer o.Wait()
	defer o.Stop()

	intvl := 1 * time.Second
	o.Every(intvl, func() {
		t.Logf("%v\n", time.Now())
	})
	time.Sleep(20 * intvl)
}

func TestSelect_DataRace(t *testing.T) {
	t.Logf("sizeof Timer: %v", unsafe.Sizeof(*time.NewTimer(0)))

	c := make(chan interface{})
	NewLoop(func() {
		c <- nil
	})

	xs := make([]reflect.SelectCase, math.MaxInt8)
	for i, _ := range xs {
		xs[i] = reflect.SelectCase{
			Dir: reflect.SelectRecv,
		}
	}
	stop := make(chan interface{})
	xs[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(stop),
	}
	go func() {
		time.Sleep(time.Second)
		stop <- nil
		xs[1] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(c),
		}
		reflect.Select(xs)
	}()

	reflect.Select(xs)
}
