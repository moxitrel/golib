/*
func NewTime(accuracy time.Duration) *Time
func (*Time) Add   (cond func() bool, do func()) *Task
func (*Time) Delete(*Task)
func (*Time) Start ()
func (*Time) Stop  ()

func (*Time) At    (time.Time    , func()) *Task
func (*Time) Loop  (time.Duration, func()) *Task
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
	Service
	accuracy time.Duration
	tasks    sets.Set
	thunk    Thunk
}

func NewTime(accuracy time.Duration) (v *Time) {
	v = &Time{
		accuracy: accuracy,
		tasks:    hashset.New(),
		thunk:    *NewThunk(),
	}
	v.Service = *New(func() {
		now := time.Now()
		time.Sleep(now.Truncate(accuracy).Add(accuracy).Sub(now) % accuracy)
		for _, taskAny := range v.tasks.Values() {
			task := taskAny.(*Task)
			v.thunk.Call(func() {
				if task.cond() == true {
					task.do()
				}
			})
		}
	})
	return
}

func (o *Time) Add(cond func() bool, do func()) (v *Task) {
	v = &Task{cond, do}
	if cond != nil && do != nil {
		o.tasks.Add(v)
	}
	return
}

func (o *Time) Delete(task *Task) {
	o.tasks.Remove(task)
}

func (o *Time) Start() {
	o.thunk.Start()
	o.Service.Start()
}

func (o *Time) Stop() {
	o.Service.Stop()
	o.thunk.Stop()
}

// Run do() at future.
// If future is before now, run immediately
func (o *Time) At(future time.Time, do func()) (v *Task) {
	v = o.Add(
		func() bool { return !time.Now().Before(future) },
		do)
	v.do = func() {
		do()
		o.Delete(v)
	}
	return
}

// Run do() every interval ns
func (o *Time) Loop(interval time.Duration, do func()) *Task {
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
