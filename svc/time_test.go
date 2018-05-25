package svc

import (
	"testing"
	"time"
)

func TestLoop(t *testing.T) {
	var accuracy = 100 * time.Millisecond

	o := NewTime(accuracy)
	defer o.Stop()

	o.Every(accuracy*2, func() {
		t.Logf("%v\n", time.Now())
	})
	time.Sleep(accuracy * 2 * 5)
}

func TestAtLoop(t *testing.T) {
	var accuracy = 100 * time.Millisecond

	o := NewTime(accuracy)
	defer o.Stop()

	o.At(time.Now(), func() {
		t.Logf("%s\n", time.Now())
	})
	time.Sleep(accuracy * 10)
}
