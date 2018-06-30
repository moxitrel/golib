package log

import (
	"runtime"
	"fmt"
	"path/filepath"
)

func CallerName(n int) string {
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
		v += CallerName(i)
		v += " -> "
	}
	v += CallerName(1)
	return v
}
