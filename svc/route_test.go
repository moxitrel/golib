package svc

import (
	"testing"
	"time"
	//"gitlab.com/clowire/gateway/v4/log"
	"reflect"
)

//
//type StringMessager struct {
//	*Route
//}
//func StringMessagerOf(x *Route, f func(string)) (v *StringMessager) {
//	v = &StringMessager{
//		Route: x,
//	}
//	v.Register(*new(string), func(x interface{}) {
//		f(x.(string))
//	})
//	return v
//}
//func (o *StringMessager) AddMessage(x string) {
//	o.Route.AddMessage(x)
//}

func TestRoute_RegisterAndCall(t *testing.T) {
	o := NewRoute()
	defer o.Stop()

	oT := reflect.Invalid
	o.Register("", func(xAny interface{}) {
		oT = reflect.String
	})
	o.Register(int(0), func(xAny interface{}) {
		oT = reflect.Int
	})

	o.Call("a")
	time.Sleep(time.Millisecond)
	if oT != reflect.String {
		t.Errorf("oT = %s, want %s", reflect.String, oT)
	}
	o.Call(5)
	time.Sleep(time.Millisecond)
	if oT != reflect.Int {
		t.Errorf("oT = %s, want %s", reflect.Int, oT)
	}
	// no handler
	o.Call(nil)
	o.Call(9)
}
