package svc

import (
	"testing"
	"time"
)

func TestService(t *testing.T) {
	i := 0
	defer func() {
		t.Logf("i: %v", i)
	}()
	defer New(func() {
		i ++
	}).Stop()
	time.Sleep(time.Millisecond)
}
