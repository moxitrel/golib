package golib

import (
	"fmt"
	"testing"
)

func TestCallerName(t *testing.T) {
	callerName := Caller(0)
	rightCallerName := "golib.TestCallerName.9"
	if callerName != rightCallerName {
		t.Errorf("Caller(0): %v /= %v", callerName, rightCallerName)
	}
}

func TestPanic(t *testing.T) {
	msg := "msg"
	defer func() {
		err := recover()
		rightErrorMsg := fmt.Sprintf("%v%v\n", "golib.TestPanic.26: ", msg)
		if err.(error).Error() != rightErrorMsg {
			t.Errorf("Panic(%v): %v /= %v", msg, err, rightErrorMsg)
		}
	}()

	Panic(msg)
}
