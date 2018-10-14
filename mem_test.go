package golib

import (
	"fmt"
	"math/rand"
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

func Test_BytesPool(t *testing.T) {
	x := BytesPool.Get(uint(rand.Int31()))
	t.Logf("x.len: %v", len(x))
	t.Logf("x.cap: %v", cap(x))
	BytesPool.Put(x)
}