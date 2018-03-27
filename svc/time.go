/*
func NewTime(accuracy time.Duration) *Time
func (*Time) Add   (cond func() bool, do func()) *Task
func (*Time) Delete(*Task)
func (*Time) Start ()
func (*Time) Stop  ()
func (*Time) At    (time.Time, func()) *Task
*/
package svc

import (
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
	"time"
)

type Task struct {
	cond func() bool //should not be blocked
	do   func()
}

type Time struct {
	*Service
	tasks sets.Set
	apply Apply
}

func NewTime(accuracy time.Duration) (v *Time) {
	v.tasks = hashset.New()
	v.apply = NewApply(TIME_APPLY_POOL_SIZE)
	v.Service = New(func() {
		now := time.Now()
		time.Sleep(now.Truncate(accuracy).Add(accuracy).Sub(now) % accuracy)
		for _, taskAny := range v.tasks.Values() {
			task, _ := taskAny.(*Task)
			if task.cond() == true {
				v.apply.Add(task.do)
			}
		}
	})
	return
}

func (o *Time) Add(cond func() bool, do func()) *Task {
	task := &Task{cond, do}
	if cond != nil && do != nil {
		o.tasks.Add(task)
	}
	return task
}

func (o *Time) Delete(task *Task) {
	o.tasks.Remove(task)
}

func (o *Time) Start() {
	o.apply.Start()
	o.Service.Start()
}

func (o *Time) Stop() {
	o.Service.Stop()
	o.apply.Stop()
}

// Run do() at future.
// If future is before now, run immediately
func (o *Time) At(future time.Time, do func()) *Task {
	task := o.Add(
		func() bool { return time.Now().After(future) },
		do)
	task.do = func() {
		do()
		o.Delete(task)
	}
	return task
}
