package golib

import (
	"errors"
	"fmt"
	"go/build"
	"path"
	"runtime"
	"strings"
)

// <dir>/<file>/<func>/<line>
func callerPath(n int) string {
	pc, file, line, ok := runtime.Caller(n)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v/%v()/%v",
		strings.TrimPrefix(file, build.Default.GOPATH),
		strings.Split(path.Base(runtime.FuncForPC(pc).Name()), ".")[1],
		line)
}

// <pkg>.<func>.<line-no>
func callerName(n int) string {
	pc, _, line, ok := runtime.Caller(n)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v:%v", path.Base(runtime.FuncForPC(pc).Name()), line)
}

func CallerTree() string {
	var s strings.Builder

	pcs := make([]uintptr, 64)
	n := runtime.Callers(0, pcs)
	s.WriteString(callerName(2))
	for i := 3; i < n; i++ {
		s.WriteString(" <- ")
		s.WriteString(callerName(i))
	}

	return s.String()
}

// panic with current function's info
func Panic(format string, args ...interface{}) {
	panic(errors.New(fmt.Sprintf(callerPath(2)+": "+format+"\n", args...)))
}

func Warn(format string, args ...interface{}) {
	fmt.Printf("WARN %v: ", callerPath(2))
	fmt.Printf(format, args...)
	fmt.Printf("\n")
}
