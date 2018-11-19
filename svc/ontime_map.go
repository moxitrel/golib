/*
func NewMapOnTime(accuracy time.Duration) *MapOnTime
func (*MapOnTime) Add   	(func())					*Task
func (*MapOnTime) Delete	(*Task)

func (*MapOnTime) At    	(time.MapOnTime,     func()) 	*Task
func (*MapOnTime) Every 	(time.Duration, func()) 	*Task
*/
package svc

import (
	"github.com/moxitrel/golib"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

var taskPool = sync.Pool{
	New: func() interface{} {
		return &Task{}
	},
}

type MapOnTime struct {
	taskLen uint64   // used as the key of tasks
	tasks   sync.Map // map[Id]*Task
	// 1s    : very slow
	// 0.5s  : slow
	// 0.4s  : normal
	// 0.3s  : quick
	// 0.25s : ok
	// 0.2s  : fast
	accuracy time.Duration
	*Loop
	//now time.Time
}

// NOTE: sleep, ticker may have a cpu overhead when idle
// 		 https://zhuanlan.zhihu.com/p/45959147
// 		 https://groups.google.com/forum/#!topic/golang-nuts/_XEuq69kUOY
// 		 https://github.com/golang/go/issues/27707
// 		 https://github.com/golang/go/issues/25471
func NewMapOnTime(accuracy time.Duration) (o *MapOnTime) {
	if accuracy <= 0 {
		golib.Panic("accuracy <= 0, want > 0")
	}
	o = &MapOnTime{
		accuracy: accuracy,
		tasks:    sync.Map{},
		taskLen:  0,
	}
	o.Loop = NewLoop(func() {
		now := time.Now()
		time.Sleep(now.Truncate(accuracy).Add(accuracy).Sub(now) % o.accuracy)

		o.tasks.Range(func(key, value interface{}) bool {
			task := value.(*Task)
			life := atomic.AddInt64(&task.life, -o.accuracy.Nanoseconds())
			switch {
			case life > 0:
				// continue to wait
			case life >= -o.accuracy.Nanoseconds():
				task.do()
			default:
				o.tasks.Delete(key)
				taskPool.Put(value)
			}
			return true
		})
	})
	return
}

func NewTickerOnTime(accuracy time.Duration) (o *MapOnTime) {
	if accuracy <= 0 {
		golib.Panic("accuracy <= 0, want > 0")
	}
	o = &MapOnTime{
		accuracy: accuracy,
		tasks:    sync.Map{},
		taskLen:  0,
	}
	var ticker *time.Ticker
	o.Loop = NewHookedLoop(func() {
		ticker = time.NewTicker(accuracy)
		now := time.Now()
		time.Sleep(now.Truncate(accuracy).Add(accuracy).Sub(now) % o.accuracy)
	}, func() {
		<-ticker.C
		o.tasks.Range(func(key, value interface{}) bool {
			task := value.(*Task)
			life := atomic.AddInt64(&task.life, -o.accuracy.Nanoseconds())
			switch {
			case life > 0:
				// continue to wait
			case life >= -o.accuracy.Nanoseconds():
				task.do()
			default:
				o.tasks.Delete(key)
			}
			return true
		})
	}, func() {
		ticker.Stop()
	})

	return
}

func (o *MapOnTime) _addTask(task *Task) {
	taskLen := atomic.AddUint64(&o.taskLen, 1)
	if taskLen == 0 {
		golib.Panic("taskLen overflow, too many adding")
	}
	o.tasks.Store(taskLen-1, task)
}

func (o *MapOnTime) Delete(task *Task) {
	if task == nil {
		return
	}
	atomic.StoreInt64(&task.life, math.MinInt64)
}

// Run thunk() once at <future>.
// If future is before now, run at next check
func (o *MapOnTime) At(future time.Time, do func()) (v *Task) {
	if do == nil {
		golib.Warn("do == nil, ignored")
		return nil
	}
	v = taskPool.Get().(*Task)
	v.do = do
	v.life = future.Sub(time.Now()).Nanoseconds()
	o._addTask(v)
	return
}

// Run thunk() every <interval> ns
// require interval < (|math.MinInt64| / 2), 146y
func (o *MapOnTime) Every(interval time.Duration, do func()) (v *Task) {
	if do == nil {
		golib.Warn("do == nil, ignored")
		return nil
	}
	if interval < o.accuracy {
		interval = o.accuracy
	}

	now := time.Now()
	v = taskPool.Get().(*Task)
	v.do = func() {
		atomic.AddInt64(&v.life, interval.Nanoseconds())
		do()
	}
	v.life = now.Truncate(interval).Add(interval).Sub(now).Nanoseconds()
	o._addTask(v)
	return
}
