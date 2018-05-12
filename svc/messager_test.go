package svc

import (
	"fmt"
	"testing"
	"time"
)

func TestMessager(t *testing.T) {
	o := NewMessager()
	o.Start()

	o.Register("0", func(string) { fmt.Printf("hi") })

	// send nil
	o.AddMessage(nil)
	// send no handler
	o.AddMessage("a")
	o.AddMessage(9)

	time.Sleep(time.Second * 1)
}
