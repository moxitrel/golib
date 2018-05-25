package svc

import (
	"testing"
	"time"
	//"gitlab.com/clogwire/v4/log"
	"fmt"
)

//
//type StringMessager struct {
//	*Dispatch
//}
//func StringMessagerOf(x *Dispatch, f func(string)) (v *StringMessager) {
//	v = &StringMessager{
//		Dispatch: x,
//	}
//	v.Register(*new(string), func(x interface{}) {
//		f(x.(string))
//	})
//	return v
//}
//func (o *StringMessager) AddMessage(x string) {
//	o.Dispatch.AddMessage(x)
//}

func TestMessager(t *testing.T) {
	o := NewDispatch()
	defer o.Stop()

	o.Register(*new(string), func(xAny interface{}) {
		x := xAny.(string)
		fmt.Printf("%v\n", x)
	})
	o.Register(nil, func(xAny interface{}) {
		fmt.Printf("nil\n")
	})

	// send nil
	o.Handle(nil)
	// send no handler
	o.Handle("a")
	o.Handle(9)

	time.Sleep(time.Second * 1)
}
