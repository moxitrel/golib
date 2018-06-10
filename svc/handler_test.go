package svc

import (
	"testing"
	"time"
	//"gitlab.com/clowire/gateway/v4/log"
	"fmt"
)

//
//type StringMessager struct {
//	*Handler
//}
//func StringMessagerOf(x *Handler, f func(string)) (v *StringMessager) {
//	v = &StringMessager{
//		Handler: x,
//	}
//	v.Register(*new(string), func(x interface{}) {
//		f(x.(string))
//	})
//	return v
//}
//func (o *StringMessager) AddMessage(x string) {
//	o.Handler.AddMessage(x)
//}

func TestMessager(t *testing.T) {
	o := NewHandler()
	defer o.Stop()

	o.Register("", func(xAny interface{}) {
		x := xAny.(string)
		fmt.Printf("%v\n", x)
	})
	o.Register(nil, func(xAny interface{}) {
		fmt.Printf("nil\n")
	})

	// send nil
	o.Do(nil)
	// send no handler
	o.Do("a")
	o.Do(9)

	time.Sleep(time.Second * 1)
}
