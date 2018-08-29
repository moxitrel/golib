package svc

import (
	"testing"
)

func TestHandler_Register_Dup(t *testing.T) {
	h := NewHandler()
	flag := 0
	f1 := func(_ interface{}) { flag = 1 }
	f2 := func(_ interface{}) { flag = 2 }

	h.Register(1, f1)
	h.Register(1, f2)
	h.Handle(1, struct{}{})
	if flag != 2 {
		t.Errorf("flag = %v, want 2", flag)
	}
}
