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
//
func TestRoute_RegisterAndCall(t *testing.T) {
	o := NewRoute(3)
	defer o.Stop()

	oT := reflect.Invalid
	o.Register("", func(xAny interface{}) {
		oT = reflect.String
	})
	o.Register(int(0), func(xAny interface{}) {
		oT = reflect.Int
	})

	o.Apply("a")
	time.Sleep(time.Millisecond)
	if oT != reflect.String {
		t.Errorf("oT = %s, want %s", reflect.String, oT)
	}
	o.Apply(5)
	time.Sleep(time.Millisecond)
	if oT != reflect.Int {
		t.Errorf("oT = %s, want %s", reflect.Int, oT)
	}
	// no handler
	o.Apply(nil)
	o.Apply(9)
}

func TestHandler(t *testing.T) {
	handler := NewHandler()
	handler.Register(1, func(_ interface{}){})
	handler.Register(1, func(_ interface{}){
		t.Log(1)
	})
	//handler.Register(func(){}, func(_ interface{}){
	//	t.Log("func")
	//})
	handler.Register(reflect.TypeOf(int(0)), nil)
	handler.Apply(1)
	handler.Apply(2)
	handler.Apply([]byte{1,2,3})
}

func TestMap(t *testing.T) {
	m := make(map[interface{}]interface{})
	delete(m, 0)
}