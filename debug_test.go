package golib

import "testing"

func TestCallerName(t *testing.T) {
	callerName := CallerName(0)
	rightCallerName := "golib.TestCallerName.6"
	if callerName != rightCallerName {
		t.Errorf("CallerName(0): %v /= %v", callerName, rightCallerName)
	}
}