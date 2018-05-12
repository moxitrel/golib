package svc

import (
	"testing"
	"gitlab.com/clogwire/v4/log"
	"time"
)

func printBytes(xs []byte) {
	for _, c := range xs {
		log.Info("%c", c)
	}
}

type PrintBytes struct {
	Fun
}

func NewParser() (v *PrintBytes) {
	v = &PrintBytes{
		Fun: *NewFun(func(x interface{}) {
			printBytes(x.([]byte))
		}),
	}
	return
}
func (o *PrintBytes) Call(x []byte) {
	o.Fun.Call(x)
}

func TestFun(t *testing.T) {
	o := NewParser()
	o.Start()
	o.Call([]byte("abcdefg"))
	time.Sleep(time.Second)
}
