package gosvc

import (
	"container/ring"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkNow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		time.Now()
	}
}

func BenchmarkTime_Sleep_0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(0)
	}
}
func BenchmarkTime_OnTime_0(b *testing.B) {
	o := &MapOnTime{
		accuracy: 0,
		tasks:    sync.Map{},
		taskLen:  0,
	}
	var task *Task
	task = o.Every(0, func() {
		atomic.StoreInt64(&task.life, 1)
	})
	for i := 0; i < b.N; i++ {
		time.Sleep(0)

		o.tasks.Range(func(key, value interface{}) bool {
			task := value.(*Task)
			switch taskLife := atomic.LoadInt64(&task.life); {
			case taskLife > 0:
				// do nothing, keep this case
			case taskLife == 0:
				task.do()
			case taskLife < 0:
				o.tasks.Delete(key)
			}
			atomic.AddInt64(&task.life, -1)
			return true
		})
	}
}
func BenchmarkTime_Sleep_1ns(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(time.Nanosecond)
	}
}
func BenchmarkTime_OnTime_1ns(b *testing.B) {
	o := &MapOnTime{
		accuracy: time.Nanosecond,
		tasks:    sync.Map{},
		taskLen:  0,
	}
	var task *Task
	task = o.Every(0, func() {
		atomic.StoreInt64(&task.life, 1)
	})
	for i := 0; i < b.N; i++ {
		now := time.Now()
		time.Sleep(now.Truncate(o.accuracy).Add(o.accuracy).Sub(now) /*% o.accuracy*/)

		o.tasks.Range(func(key, value interface{}) bool {
			task := value.(*Task)
			switch taskLife := atomic.LoadInt64(&task.life); {
			case taskLife > 0:
				// do nothing, keep this case
			case taskLife == 0:
				task.do()
			case taskLife < 0:
				o.tasks.Delete(key)
			}
			atomic.AddInt64(&task.life, -1)
			return true
		})
	}
}
func BenchmarkTime_Ticker_1ns(b *testing.B) {
	o := time.NewTicker(time.Nanosecond)
	for i := 0; i < b.N; i++ {
		<-o.C
	}
}
func BenchmarkTime_Sleep_1us(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(time.Microsecond)
	}
}
func BenchmarkTime_OnTime_1us(b *testing.B) {
	o := &MapOnTime{
		accuracy: time.Microsecond,
		tasks:    sync.Map{},
		taskLen:  0,
	}
	var task *Task
	task = o.Every(0, func() {
		atomic.StoreInt64(&task.life, 1)
	})
	for i := 0; i < b.N; i++ {
		now := time.Now()
		time.Sleep(now.Truncate(o.accuracy).Add(o.accuracy).Sub(now) /*% o.accuracy*/)

		o.tasks.Range(func(key, value interface{}) bool {
			task := value.(*Task)
			switch taskLife := atomic.LoadInt64(&task.life); {
			case taskLife > 0:
				// do nothing, keep this case
			case taskLife == 0:
				task.do()
			case taskLife < 0:
				o.tasks.Delete(key)
			}
			atomic.AddInt64(&task.life, -1)
			return true
		})
	}
}
func BenchmarkTime_Ticker_1us(b *testing.B) {
	o := time.NewTicker(time.Microsecond)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-o.C
	}
}
func BenchmarkTime_Sleep_10us(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(10 * time.Microsecond)
	}
}
func BenchmarkTime_OnTime_10us(b *testing.B) {
	o := &MapOnTime{
		accuracy: 10 * time.Microsecond,
		tasks:    sync.Map{},
		taskLen:  0,
	}
	var task *Task
	task = o.Every(0, func() {
		atomic.StoreInt64(&task.life, 1)
	})
	for i := 0; i < b.N; i++ {
		now := time.Now()
		time.Sleep(now.Truncate(o.accuracy).Add(o.accuracy).Sub(now) % o.accuracy)

		o.tasks.Range(func(key, value interface{}) bool {
			task := value.(*Task)
			switch taskLife := atomic.LoadInt64(&task.life); {
			case taskLife > 0:
				// do nothing, keep this case
			case taskLife == 0:
				task.do()
			case taskLife < 0:
				o.tasks.Delete(key)
			}
			atomic.AddInt64(&task.life, -1)
			return true
		})
	}
}
func BenchmarkTime_Ticker_10us(b *testing.B) {
	o := time.NewTicker(10 * time.Microsecond)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-o.C
	}
}
func BenchmarkTime_Sleep_100us(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(100 * time.Microsecond)
	}
}
func BenchmarkTime_OnTime_100us(b *testing.B) {
	o := &MapOnTime{
		accuracy: 100 * time.Microsecond,
		tasks:    sync.Map{},
		taskLen:  0,
	}
	var task *Task
	task = o.Every(0, func() {
		atomic.StoreInt64(&task.life, 1)
	})
	for i := 0; i < b.N; i++ {
		now := time.Now()
		time.Sleep(now.Truncate(o.accuracy).Add(o.accuracy).Sub(now) % o.accuracy)

		o.tasks.Range(func(key, value interface{}) bool {
			task := value.(*Task)
			switch taskLife := atomic.LoadInt64(&task.life); {
			case taskLife > 0:
				// do nothing, keep this case
			case taskLife == 0:
				task.do()
			case taskLife < 0:
				o.tasks.Delete(key)
			}
			atomic.AddInt64(&task.life, -1)
			return true
		})
	}
}
func BenchmarkTime_Ticker_100us(b *testing.B) {
	o := time.NewTicker(100 * time.Microsecond)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-o.C
	}
}
func BenchmarkTime_Sleep_1ms(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(time.Millisecond)
	}
}
func BenchmarkTime_OnTime_1ms(b *testing.B) {
	o := &MapOnTime{
		accuracy: time.Millisecond,
		tasks:    sync.Map{},
		taskLen:  0,
	}
	var task *Task
	task = o.Every(0, func() {
		atomic.StoreInt64(&task.life, 1)
	})
	for i := 0; i < b.N; i++ {
		now := time.Now()
		time.Sleep(now.Truncate(o.accuracy).Add(o.accuracy).Sub(now) % o.accuracy)

		o.tasks.Range(func(key, value interface{}) bool {
			task := value.(*Task)
			switch taskLife := atomic.LoadInt64(&task.life); {
			case taskLife > 0:
				// do nothing, keep this case
			case taskLife == 0:
				task.do()
			case taskLife < 0:
				o.tasks.Delete(key)
			}
			atomic.AddInt64(&task.life, -1)
			return true
		})
	}
}
func BenchmarkTime_Ticker_1ms(b *testing.B) {
	o := time.NewTicker(time.Millisecond)
	for i := 0; i < b.N; i++ {
		<-o.C
	}
}

func BenchmarkAddTask_OnTimeAdd(b *testing.B) {
	o := NewOnTime(time.Minute)
	now := time.Now()
	for i := 0; i < b.N; i++ {
		o.At(now.Add(time.Minute), func() {})
	}
}
func BenchmarkAddTask_NewTimer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.NewTimer(-1)
	}
}

func BenchmarkFor_Array8(b *testing.B) {
	o := make([]interface{}, math.MaxUint8)
	for i := 0; i < len(o); i++ {
		o[i] = i
	}
	f := func(interface{}) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(o); i++ {
			f(o[i])
		}
	}
}
func BenchmarkFor_List8(b *testing.B) {
	o := ring.New(math.MaxUint8)
	for i := 0; i < o.Len(); i++ {
		o.Value = i
		o = o.Next()
	}
	lock := (&sync.RWMutex{}).RLocker()
	f := func(interface{}) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.Lock()
		p := o.Next()
		lock.Unlock()
		for ; p != o; p = p.Next() {
			f(p.Value.(interface{}))
		}
	}
}
func BenchmarkFor_SyncMap8(b *testing.B) {
	o := sync.Map{}
	for i := 0; i < math.MaxUint8; i++ {
		o.Store(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.Range(func(_, _ interface{}) bool {
			return true
		})
	}
}
func BenchmarkFor_Array1000(b *testing.B) {
	o := make([]interface{}, 1000)
	for i := 0; i < len(o); i++ {
		o[i] = i
	}
	f := func(interface{}) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(o); i++ {
			f(o[i])
		}
	}
}
func BenchmarkFor_List1000(b *testing.B) {
	o := ring.New(1000)
	for i := 0; i < o.Len(); i++ {
		o.Value = i
		o = o.Next()
	}
	lock := (&sync.RWMutex{}).RLocker()
	f := func(interface{}) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.Lock()
		p := o.Next()
		lock.Unlock()
		for ; p != o; p = p.Next() {
			f(p.Value.(interface{}))
		}
	}
}
func BenchmarkFor_SyncMap1000(b *testing.B) {
	o := sync.Map{}
	for i := 0; i < 1000; i++ {
		o.Store(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.Range(func(_, _ interface{}) bool {
			return true
		})
	}
}
func BenchmarkFor_Array10000(b *testing.B) {
	o := make([]interface{}, 10000)
	for i := 0; i < len(o); i++ {
		o[i] = i
	}
	f := func(interface{}) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(o); i++ {
			f(o[i])
		}
	}
}
func BenchmarkFor_List10000(b *testing.B) {
	o := ring.New(10000)
	for i := 0; i < o.Len(); i++ {
		o.Value = i
		o = o.Next()
	}
	lock := (&sync.RWMutex{}).RLocker()
	f := func(interface{}) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.Lock()
		p := o.Next()
		lock.Unlock()
		for ; p != o; p = p.Next() {
			f(p.Value.(interface{}))
		}
	}
}
func BenchmarkFor_SyncMap10000(b *testing.B) {
	o := sync.Map{}
	for i := 0; i < 10000; i++ {
		o.Store(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.Range(func(_, _ interface{}) bool {
			return true
		})
	}
}
func BenchmarkFor_Array15(b *testing.B) {
	o := make([]interface{}, math.MaxInt16)
	for i := 0; i < len(o); i++ {
		o[i] = i
	}
	f := func(interface{}) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(o); i++ {
			f(o[i])
		}
	}
}
func BenchmarkFor_List15(b *testing.B) {
	o := ring.New(math.MaxInt16)
	for i := 0; i < o.Len(); i++ {
		o.Value = i
		o = o.Next()
	}
	lock := (&sync.RWMutex{}).RLocker()
	f := func(interface{}) {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.Lock()
		p := o.Next()
		lock.Unlock()
		for ; p != o; p = p.Next() {
			f(p.Value.(interface{}))
		}
	}
}
func BenchmarkFor_SyncMap15(b *testing.B) {
	o := sync.Map{}
	for i := 0; i < math.MaxInt16; i++ {
		o.Store(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.Range(func(_, _ interface{}) bool {
			return true
		})
	}
}
