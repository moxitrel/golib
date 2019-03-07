package gosvc

import (
	"fmt"
	"testing"
)

//
// Tests
//
func TestSvc_State(t *testing.T) {
	o := NewSvc(nil, nil, func() {})

	// state = ST_RUNNING when created
	if o.State() != ST_RUNNING {
		t.Errorf("state != ST_RUNNING")
	}

	// state = ST_STOPPED after .Stop()
	o.Stop()
	if o.State() != ST_STOPPED {
		t.Errorf("state != ST_STOPPED")
	}
}

func TestSvc_DataRace(t *testing.T) {
	o := NewSvc(nil, nil, func() {})
	for i := 0; i < 2; i++ {
		go func() {
			for {
				o.State()
			}
		}()
		go func() {
			for {
				o.Stop()
			}
		}()
	}
}

//
// Benchmarks
//
func BenchmarkSvc_SwitchTest(b *testing.B) {
	o := &Svc{
		state: ST_RUNNING,
	}
	do := func() {}

	for i := 0; i < b.N; i++ {
		switch o.State() {
		case ST_RUNNING:
			do()
		case ST_STOPPED:
			return
		default:
			panic(fmt.Errorf("state:%v isn't valid", o.State()))
		}
	}
}

func BenchmarkSvc_NoTest(b *testing.B) {
	do := func() {}

	for i := 0; i < b.N; i++ {
		do()
	}
}
