package golib

import (
	"testing"
)

func TestCallerName(t *testing.T) {
	callerName := callerPos(1)
	rightCallerName := "golib.TestCallerName.8"
	if callerName != rightCallerName {
		t.Errorf("callerPos(0): %v /= %v", callerName, rightCallerName)
	}
}

//func TestPanic(t *testing.T) {
//	msg := "msg"
//	defer func() {
//		err := recover()
//		rightErrorMsg := fmt.Sprintf("%v%v\n", "golib.TestPanic.26: ", msg)
//		if err.(error).Error() != rightErrorMsg {
//			t.Errorf("Panic(%v): %v /= %v", msg, err, rightErrorMsg)
//		}
//	}()
//
//	Panic(msg)
//}

func TestCallerPos(t *testing.T) {
	t.Logf("%v", CallerPos())
}
