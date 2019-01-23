package gosvc

import "time"

// Wrap time.Timer to simplify the usage
//
// * Example
// timer := NewTimer()
//
// timer.Start(x)		// call .Start() at beginning
// select {
// case c <- nil :
// 		// not timeout
// case <-c      :
// 		// not timeout
// case <-timer.C:
// 		// timeout
// }
// timer.Stop()			// call .Stop() at end
//
type Timer time.Timer

func NewTimer() (o *Timer) {
	o = (*Timer)(time.NewTimer(time.Second))
	o.Stop()
	return
}

func (o *Timer) Start(timeout time.Duration) {
	(*time.Timer)(o).Reset(timeout)
}

func (o *Timer) Stop() {
	if !(*time.Timer)(o).Stop() {
		<-o.C
	}
}
