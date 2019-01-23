package gosvc

import (
	"sync"
	"testing"
)

func TestLoop_Example(t *testing.T) {
	var n uint64 = 0
	var signalOnce sync.Once
	var startSignal = make(chan struct{})
	o := NewLoop(func() {
		n++
		signalOnce.Do(func() {
			startSignal <- struct{}{}
		})
	})
	<-startSignal
	o.Stop()
	o.Wait()

	if n == 0 {
		t.Errorf("n == %v, want !0", n)
	}
}

func TestLoop_Wait(t *testing.T) {
	o := NewLoop(func() {})
	o.Stop()
	o.Wait()
}

func TestLoop_DataRace(t *testing.T) {
	o := NewLoop(func() {})

	for i := 0; i < 2; i++ {
		go func() {
			for {
				o.State()
			}
		}()
	}
	for i := 0; i < 2; i++ {
		go func() {
			for {
				o.Stop()
			}
		}()
	}
	for i := 0; i < 2; i++ {
		go func() {
			for {
				o.Wait()
			}
		}()
	}
	o.Wait()
}
