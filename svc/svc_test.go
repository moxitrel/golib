package svc

import (
	"testing"
)

func TestSvc_DataRace(t *testing.T) {
	o := NewSvc(nil, nil, func() {})
	for i := 0; i < 2; i++ {
		NewLoop(func() {
			o.State()
		})
	}
	for i := 0; i < 2; i++ {
		NewLoop(func() {
			o.Stop()
		})
	}
}
