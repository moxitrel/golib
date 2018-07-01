package golib

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
)

// <pkg>.<func>.<line-no>
func Caller(n int) string {
	pc, _, line, ok := runtime.Caller(n + 1)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%s.%d", filepath.Base(runtime.FuncForPC(pc).Name()), line)
}

func CallTree(n int) string {
	n++
	var v = ""
	for i := n; i > 1; i-- {
		v += Caller(i)
		v += " -> "
	}
	v += Caller(1)
	return v
}

// panic with current function's info
func Panic(xs ...interface{}) {
	panic(errors.New(Caller(1) + ": " + fmt.Sprint(xs...)))
}
