package svc

import (
	"sync"
	"testing"
)

func TestSvc_Stop(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	o := NewSvc(nil, func() {

	}, func() {
		wg.Done()
	})
	o.Stop()
	wg.Wait()
}
