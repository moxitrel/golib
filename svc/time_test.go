package svc

import (
	"testing"
	"time"
)

func Test_TimeEvery(t *testing.T) {
	var accuracy = 100 * time.Millisecond

	o := NewTime(accuracy)
	defer o.Join()
	defer o.Stop()

	intvl := 2 * accuracy
	o.Every(intvl, func() {
		t.Logf("%v\n", time.Now())
	})
	time.Sleep(5 * intvl)
}

func TestTime(t *testing.T) {
	i := 0
	timed := NewTime(time.Millisecond)
	timed.add(func() {
		i++
		t.Logf("%v: %v", i, time.Now())
	})
	time.Sleep(time.Second)
	timed.Stop()
	timed.Join()
}
