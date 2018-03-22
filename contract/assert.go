package contract

import "fmt"

func Assert(ok bool, format string, args ...interface{}) {
	if !ok {
		panic(fmt.Sprintf(format, args...))
	}
}