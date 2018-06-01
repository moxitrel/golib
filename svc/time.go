/*
func NewTime(accuracy time.Duration) *Time
func (*Time) Add   (func()) *Task
func (*Time) Delete(*Task)

func (*Time) RunOnce(func() bool) 					*Task
func (*Time) At    	(time.Time    	 , func()) 		*Task
func (*Time) Every  (time.Duration	 , func()) 		*Task
*/
package svc

import (
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
	"time"
)

type Task struct{ thunk func() }

type Time struct {
	*Loop
	accuracy time.Duration
	tasks    sets.Set
}

func NewTime(accuracy time.Duration) (v *Time) {
	v = &Time{
		accuracy: accuracy,
		tasks:    hashset.New(),
	}
	v.Loop = NewLoop(func() {
		now := time.Now()
		time.Sleep(now.Truncate(v.accuracy).Add(v.accuracy).Sub(now) % v.accuracy)

		for _, value := range v.tasks.Values() {
			task := value.(*Task)
			task.thunk()
		}
	})
	return
}

func (o *Time) Add(thunk func()) (v *Task) {
	v = &Task{thunk: thunk}
	if thunk != nil {
		o.tasks.Add(v)
	}
	return
}

func (o *Time) Delete(task *Task) {
	o.tasks.Remove(task)
}

// Remove the task if do() return true
func (o *Time) RunOnce(do func() bool) (v *Task) {
	v = o.Add(func() {}) //make a placeholder
	v.thunk = func() {
		if do() == true {
			o.Delete(v)
		}
	}
	return
}

// Run thunk() once at <future>.
// If future is before now, run at next check
func (o *Time) At(future time.Time, thunk func()) (v *Task) {
	v = o.RunOnce(func() bool {
		if time.Now().Before(future) {
			return false
		} else {
			thunk()
			return true
		}
	})
	return
}

// Run thunk() every <interval> ns
func (o *Time) Every(interval time.Duration, thunk func()) (v *Task) {
	tnext := time.Now().Truncate(interval).Add(interval)
	v = o.Add(func() {
		now := time.Now()
		if !now.Before(tnext) {
			tnext = now.Truncate(interval).Add(interval)
			thunk()
		}
	})
	return
}
