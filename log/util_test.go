package log

import "testing"

func TestCallerName(t *testing.T) {
	callerName := CallerName(0)
	if callerName != "log.TestCallerName.6" {
		t.Errorf("CallerName(0): %v; want log.TestCallerName.6", callerName)
	}
}