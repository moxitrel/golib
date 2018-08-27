package svc

import (
	"testing"
	"time"
	//"gitlab.com/clowire/gateway/v4/log"
	"reflect"
)

//
//type StringMessager struct {
//	*HandlerService
//}
//func StringMessagerOf(x *HandlerService, f func(string)) (v *StringMessager) {
//	v = &StringMessager{
//		HandlerService: x,
//	}
//	v.Register(*new(string), func(x interface{}) {
//		f(x.(string))
//	})
//	return v
//}
//func (o *StringMessager) AddMessage(x string) {
//	o.HandlerService.AddMessage(x)
//}
//
func TestRoute_RegisterAndCall(t *testing.T) {
	o := NewHandlerService(3)
	defer o.Stop()

	oT := reflect.Invalid
	o.Register("", func(xAny interface{}) {
		oT = reflect.String
	})
	o.Register(int(0), func(xAny interface{}) {
		oT = reflect.Int
	})

	o.Handle("a")
	time.Sleep(time.Millisecond)
	if oT != reflect.String {
		t.Errorf("oT = %s, want %s", reflect.String, oT)
	}
	o.Handle(5)
	time.Sleep(time.Millisecond)
	if oT != reflect.Int {
		t.Errorf("oT = %s, want %s", reflect.Int, oT)
	}
	// no handler
	o.Handle(nil)
	o.Handle(9)
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
	handler.Handle(1)
	handler.Handle(2)
	handler.Handle(nil)
	handler.Handle([]byte{1,2,3})
}

func TestMap(t *testing.T) {
	m := make(map[interface{}]interface{})
	delete(m, 0)
}

func TestHandleNil(t *testing.T) {
	t.Logf("%t", reflect.TypeOf(nil))
}