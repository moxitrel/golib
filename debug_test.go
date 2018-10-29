package golib

import (
	"go/build"
	"testing"
)

func TestCaller(t *testing.T) {
	t.Logf("gopath: %v", build.Default.GOPATH)
	t.Logf("callerPath: %v", callerPath(0))
	t.Logf("callName: %v", callerName(0))
	t.Logf("CallTree: %v", CallerTree())
}
