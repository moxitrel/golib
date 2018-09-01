/*
func NewTimeService(accuracy time.Duration) *TimeService
func (*TimeService) Add   	(func(*Task))					*Task
func (*TimeService) Delete	(*Task)

func (*TimeService) At    	(time.TimeService, func()) 		*Task
func (*TimeService) Every   (time.Duration	 , func()) 		*Task
*/
package svc

import (
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/moxitrel/golib"
	"time"
)

type Task struct{ do func() }

type TimeService struct {
	*LoopService
	accuracy time.Duration
	tasks    sets.Set
}

func NewTimeService(accuracy time.Duration) (v *TimeService) {
	v = &TimeService{
		accuracy: accuracy,
		tasks:    hashset.New(),
	}
	v.LoopService = NewLoopService(func() {
		now := time.Now()
		time.Sleep(now.Truncate(v.accuracy).Add(v.accuracy).Sub(now) % v.accuracy)

		for _, value := range v.tasks.Values() {
			task := value.(*Task)
			task.do()
		}
	})
	return
}

func (o *TimeService) Add(do func()) (v *Task) {
	v = &Task{
		do: do,
	}
	if do == nil {
		golib.Warn("^do shouldn't be nil!\n")
		return
	}
	o.tasks.Add(v)
	return
}

func (o *TimeService) Delete(task *Task) {
	o.tasks.Remove(task)
}

// Run thunk() once at <future>.
// If future is before now, run at next check
func (o *TimeService) At(future time.Time, thunk func()) (v *Task) {
	v = o.Add(func() {
		if !time.Now().Before(future) {
			thunk()
			o.Delete(v)
		}
	})
	return
}

// Run thunk() every <interval> ns
func (o *TimeService) Every(interval time.Duration, thunk func()) (v *Task) {
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
