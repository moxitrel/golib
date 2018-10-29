package golib

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
)

// <pkg>.<func>.<line-no>
func callerPos(n int) string {
	pc, _, line, ok := runtime.Caller(n)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s.%d", path.Base(runtime.FuncForPC(pc).Name()), line)
}

func CallerPos() string {
	var s strings.Builder

	pcs := make([]uintptr, 64)
	n := runtime.Callers(0, pcs)
	for i := n - 1; i > 2; i-- {
		s.WriteString(callerPos(i))
		s.WriteString(" -> ")
	}
	s.WriteString(callerPos(2))
	return s.String()
}

// panic with current function's info
func Panic(format string, args ...interface{}) {
	panic(errors.New(fmt.Sprintf(CallerPos()+": "+format+"\n", args...)))
}

func Warn(format string, args ...interface{}) {
	fmt.Printf("WARN %v: ", CallerPos())
	fmt.Printf(format, args...)
	fmt.Printf("\n")
}
