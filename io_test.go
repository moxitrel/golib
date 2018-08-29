package golib

import (
	"testing"
)

func TestWriteAll(t *testing.T) {
	// panic if writer = nil
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("err = nil; want !nil")
		} else {
			t.Logf("err: %v", err)
		}
	}()
	WriteAll(nil, nil)
}
