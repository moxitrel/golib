package golib

import (
	"errors"
	"fmt"
	"path"
	"runtime"
)

// <pkg>.<func>.<line-no>
func Caller(n int) string {
	pc, _, line, ok := runtime.Caller(n + 1)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s.%d", path.Base(runtime.FuncForPC(pc).Name()), line)
}

func CallTree(n int) (v string) {
	for i := n + 1; i > 1; i-- {
		v += Caller(i)
		v += " : "
	}
	v += Caller(1)
	return v
}

// panic with current function's info
func Panic(format string, args ...interface{}) {
	panic(errors.New(fmt.Sprintf(Caller(1)+": "+format+"\n", args...)))
}

func Warn(format string, args ...interface{}) {
	fmt.Printf("WARN "+Caller(1)+": "+format+"\n", args...)
}
