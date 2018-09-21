package golib

import (
	"fmt"
	"testing"
)

func TestDefer(t *testing.T) {
	x := BytesPool.Get(0)
	fmt.Printf("init: %v\n", len(x))

	defer func() {
		fmt.Printf("defer: %v\n", len(x))
	}()

	x = append(x, make([]byte, 1024)...)
	fmt.Printf("update: %v\n", len(x))
}
