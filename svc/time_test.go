package svc

import (
	"testing"
	"time"
)

func Test_TimeEvery(t *testing.T) {
	var accuracy = 100 * time.Millisecond

	o := NewTimeService(accuracy)
	defer o.Join()
	defer o.Stop()

	intvl := 2 * accuracy
	o.Every(intvl, func() {
		t.Logf("%v\n", time.Now())
	})
	time.Sleep(5 * intvl)
}
