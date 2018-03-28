package err

import "fmt"

func Require(ok bool, format string, args ...interface{}) {
	if !ok {
		panic(fmt.Sprintf(format, args...))
	}
}

var Ensure = Require
