/*
func NewTime(accuracy time.Duration) *Time
func (*Time) Add   (cond func() bool, do func()) *Task
func (*Time) Delete(*Task)
func (*Time) Stop  ()

func (*Time) RunOnce(cond func() bool, do func()) 	*Task
func (*Time) At    	(time.Time    	 , func()) 		*Task
func (*Time) Every  (time.Duration	 , func()) 		*Task
*/
package svc

import (
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
	"time"
)

type Task struct {
	cond func() bool
	do   func()
}

type Time struct {
	Loop
	accuracy     time.Duration
	tasks        sets.Set
	thunkService Thunk
}

func NewTime(accuracy time.Duration) (v *Time) {
	v = &Time{
		accuracy:     accuracy,
		tasks:        hashset.New(),
		thunkService: *NewThunk(),
	}
	v.Loop = NewLoop(func() {
		now := time.Now()
		time.Sleep(now.Truncate(v.accuracy).Add(v.accuracy).Sub(now) % v.accuracy)

		for _, taskAny := range v.tasks.Values() {
			task := taskAny.(*Task)
			v.thunkService.Do(func() {
				if task.cond() == true {
					task.do()
				}
			})
		}

		// "send on closed channel" error if stop thunkService together with Loop.Stop(), for tasks is still being traversed
		if v.Loop.state != RUNNING {
			v.thunkService.Stop()
		}
	})
	return
}

func (o *Time) Add(cond func() bool, do func()) (v *Task) {
	v = &Task{cond, do}
	if cond != nil && do != nil {
		o.tasks.Add(v)
	} else {
		// todo: issue warning or panic
	}
	return
}

func (o *Time) Delete(task *Task) {
	o.tasks.Remove(task)
}

func (o *Time) RunOnce(cond func() bool, do func()) (v *Task) {
	v = o.Add(func() bool { return false }, func() {})		// make placeholder task which never run
	v.do = func() {
		do()
		o.Delete(v)
	}
	v.cond = cond	// add task finish
	return
}

// Run do() at future.
// If future is before now, run immediately
func (o *Time) At(future time.Time, do func()) (v *Task) {
	v = o.RunOnce(
		func() bool { return !time.Now().Before(future) },
		do)
	return
}

// Run do() every interval ns
func (o *Time) Every(interval time.Duration, do func()) *Task {
	tnext := time.Now().Truncate(interval).Add(interval)
	return o.Add(func() bool {
		now := time.Now()
		if now.Before(tnext) {
			return false
		} else {
			tnext = now.Truncate(interval).Add(interval)
			return true
		}
	}, do)
}
