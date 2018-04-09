package err

import (
	"errors"
	"fmt"
)

func Require(ok bool, format string, args ...interface{}) {
	if !ok {
		panic(errors.New(fmt.Sprintf(format, args...)))
	}
}

var Ensure = Require
