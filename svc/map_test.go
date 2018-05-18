package svc

import (
	"testing"
	"time"
	//"gitlab.com/clogwire/v4/log"
	"fmt"
)
//
//type StringMessager struct {
//	*Map
//}
//func StringMessagerOf(x *Map, f func(string)) (v *StringMessager) {
//	v = &StringMessager{
//		Map: x,
//	}
//	v.Register(*new(string), func(x interface{}) {
//		f(x.(string))
//	})
//	return v
//}
//func (o *StringMessager) AddMessage(x string) {
//	o.Map.AddMessage(x)
//}

func TestMessager(t *testing.T) {
	o := NewMap()
	o.Start()
	defer o.Stop()

	o.Register(*new(string), func(xAny interface{}) {
		x := xAny.(string)
		fmt.Printf("%v\n", x)
	})
	o.Register(nil, func(xAny interface{}) {
		fmt.Printf("nil")
	})

	// send nil
	o.Call(nil)
	// send no handler
	o.Call("a")
	o.Call(9)

	time.Sleep(time.Second * 1)
}
