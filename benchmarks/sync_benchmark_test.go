package benchmarks

import (
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

func BenchmarkAtomic_Read(b *testing.B) {
	var o uintptr = 98765
	for i := 0; i < b.N; i++ {
		atomic.LoadUintptr(&o)
	}
}
func BenchmarkAtomic_Write(b *testing.B) {
	var o uintptr = 98765
	for i := 0; i < b.N; i++ {
		atomic.StoreUintptr(&o, 6321)
	}
}
func BenchmarkAtomic_Add(b *testing.B) {
	var o uintptr = 98765
	for i := 0; i < b.N; i++ {
		atomic.AddUintptr(&o, 6321)
	}
}
func BenchmarkAtomic_CAS_F(b *testing.B) {
	var o uintptr = 98765
	for i := 0; i < b.N; i++ {
		atomic.CompareAndSwapUintptr(&o, 0, 9)
	}
}
func BenchmarkAtomic_CAS_S(b *testing.B) {
	var o uintptr = 98765
	for i := 0; i < b.N; i++ {
		atomic.CompareAndSwapUintptr(&o, o, 9)
	}
}
func BenchmarkAtomic_Pointer_Read(b *testing.B) {
	var o unsafe.Pointer
	for i := 0; i < b.N; i++ {
		atomic.LoadPointer(&o)
	}
}
func BenchmarkAtomic_Pointer_Write(b *testing.B) {
	var o unsafe.Pointer
	for i := 0; i < b.N; i++ {
		atomic.StorePointer(&o, unsafe.Pointer(nil))
	}
}
func BenchmarkAtomic_Interface_Read(b *testing.B) {
	var o atomic.Value
	o.Store(98765)
	for i := 0; i < b.N; i++ {
		o.Load()
	}
}
func BenchmarkAtomic_Interface_Write(b *testing.B) {
	var o atomic.Value
	for i := 0; i < b.N; i++ {
		o.Store(98765)
	}
}

//func BenchmarkAtomic_64_1(b *testing.B) {
//	x := int64(0)
//	for i := 0; i < b.N; i++ {
//		atomic.AddInt64(&x, -1)
//	}
//}
//func BenchmarkAtomic_64_N(b *testing.B) {
//	x := int64(0)
//	rand.Seed(time.Now().UnixNano())
//	d := rand.Int63()
//	for i := 0; i < b.N; i++ {
//		atomic.AddInt64(&x, -d)
//	}
//}
//func BenchmarkAtomic_32_1(b *testing.B) {
//	x := int32(0)
//	for i := 0; i < b.N; i++ {
//		atomic.AddInt32(&x, 1)
//	}
//}
//func BenchmarkAtomic_32_N(b *testing.B) {
//	x := int32(0)
//	rand.Seed(time.Now().UnixNano())
//	d := rand.Int31()
//	for i := 0; i < b.N; i++ {
//		atomic.AddInt32(&x, d)
//	}
//}
func BenchmarkLock_Mutex_Read(b *testing.B) {
	o := sync.Mutex{}
	for i := 0; i < b.N; i++ {
		o.Lock()
		o.Unlock()
	}
}
func BenchmarkLock_Mutex_Write(b *testing.B) {
	o := sync.Mutex{}
	x := uintptr(0)
	for i := 0; i < b.N; i++ {
		o.Lock()
		x = 98765
		o.Unlock()
	}
	// avoid x declared and not used
	atomic.LoadUintptr((*uintptr)(&x))
}
func BenchmarkLock_RWMutex_RLock(b *testing.B) {
	o := sync.RWMutex{}
	for i := 0; i < b.N; i++ {
		o.RLock()
		o.RUnlock()
	}
}
func BenchmarkLock_RWMutex_WLock(b *testing.B) {
	o := sync.RWMutex{}
	x := uintptr(0)
	for i := 0; i < b.N; i++ {
		o.Lock()
		x = 98765
		o.Unlock()
	}
	// avoid x declared and not used
	atomic.LoadUintptr((*uintptr)(&x))
}

//func BenchmarkLock_Cond_Mutex(b *testing.B) {
//	o := sync.NewCond(&sync.Mutex{})
//	for i := 0; i < b.N; i++ {
//		o.L.Lock()
//		o.L.Unlock()
//	}
//}
