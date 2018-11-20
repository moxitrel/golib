/*
func NewListOnTime(accuracy time.Duration) *OnTime
func (*OnTime) Add   	(func())					*Task
func (*OnTime) Delete	(*Task)

func (*OnTime) At    	(time.OnTime,     func()) 	*Task
func (*OnTime) Every 	(time.Duration, func()) 	*Task
*/
package svc

import (
	"container/ring"
	"github.com/moxitrel/golib"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type Task struct {
	// >0: wait
	life int64
	do   func()
}

type OnTime struct {
	*Loop
	// 1s    : very slow
	// 0.5s  : slow
	// 0.4s  : normal
	// 0.3s  : quick
	// 0.25s : ok
	// 0.2s  : fast
	accuracy  time.Duration
	tasks     *ring.Ring
	tasksLock sync.Locker
}

// NOTE: sleep, ticker may have a cpu overhead when idle
// 		 https://zhuanlan.zhihu.com/p/45959147
// 		 https://groups.google.com/forum/#!topic/golang-nuts/_XEuq69kUOY
// 		 https://github.com/golang/go/issues/27707
// 		 https://github.com/golang/go/issues/25471
func NewOnTime(accuracy time.Duration) (o *OnTime) {
	if accuracy <= 0 {
		golib.Panic("accuracy <= 0, want > 0")
	}
	o = &OnTime{
		accuracy:  accuracy,
		tasks:     ring.New(1), // the first node is preserved, used as mark
		tasksLock: new(sync.Mutex),
	}
	o.Loop = NewLoop(func() {
		now := time.Now()
		time.Sleep(now.Truncate(accuracy).Add(accuracy).Sub(now) % o.accuracy)

		o.tasksLock.Lock()
		firstNode := o.tasks.Next()
		o.tasksLock.Unlock()
		if firstNode != o.tasks { // not empty?
			// start from 2nd node, leave firstNode at last, no lock needed when deleting
			for node := firstNode.Next(); node != o.tasks; node = node.Next() {
				task := node.Value.(*Task)
				life := atomic.AddInt64(&task.life, -o.accuracy.Nanoseconds())
				switch {
				case life > 0:
					// continue to wait
				case life >= -o.accuracy.Nanoseconds():
					task.do()
				default:
					// delete the task
					node.Prev().Unlink(1)
				}
			}
			// access the firstNode, need lock when deleting
			node := firstNode
			task := node.Value.(*Task)
			life := atomic.AddInt64(&task.life, -o.accuracy.Nanoseconds())
			switch {
			case life > 0:
			case life >= -o.accuracy.Nanoseconds():
				task.do()
			default:
				o.tasksLock.Lock()
				node.Prev().Unlink(1)
				o.tasksLock.Unlock()
			}
		}
	})
	return
}

func (o *OnTime) _addTask(task *Task) {
	// add after head
	o.tasksLock.Lock()
	o.tasks.Link(&ring.Ring{Value: task})
	o.tasksLock.Unlock()
}

func (o *OnTime) Delete(task *Task) {
	if task == nil {
		return
	}
	atomic.StoreInt64(&task.life, math.MinInt64)
}

// Run thunk() once at <future>.
// If future is before now, run at next check
func (o *OnTime) At(future time.Time, do func()) (v *Task) {
	if do == nil {
		golib.Warn("do == nil, ignored")
		return nil
	}
	v = &Task{
		do:   do,
		life: future.Sub(time.Now()).Nanoseconds(),
	}
	o._addTask(v)
	return
}

// Run thunk() every <interval> ns
// require interval < (|math.MinInt64| / 2), 146y
func (o *OnTime) Every(interval time.Duration, do func()) (v *Task) {
	if do == nil {
		golib.Warn("do == nil, ignored")
		return nil
	}
	if interval < o.accuracy {
		interval = o.accuracy
	}

	now := time.Now()
	v = &Task{
		do: func() {
			atomic.AddInt64(&v.life, interval.Nanoseconds())
			do()
		},
		life: now.Truncate(interval).Add(interval).Sub(now).Nanoseconds(),
	}
	o._addTask(v)
	return
}
