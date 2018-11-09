/*

func         PoolWrap (func (interface{})                        ) *Pool
func (*Pool) WithCount(min   uint         , max     uint         ) *Pool
func (*Pool) WithTime (delay time.Duration, timeout time.Duration) *Pool

func (*Pool) Call    (interface{})

*/
package svc

import (
	"github.com/moxitrel/golib"
	"math"
	"sync/atomic"
	"time"
)

const (
	_POOL_MIN     = 0
	_POOL_MAX     = math.MaxUint16
	_POOL_DELAY   = 100 * time.Millisecond
	_POOL_TIMEOUT = time.Minute
)

// Start [min, max] goroutines of <Pool.fun> to process <Pool.arg>
//
// * Example
// f := func(x interface{}) { time.Sleep(time.Second) }
// p := PoolWrap(f)	    // start 1 goroutines of f
// p.Call("1")			// run f("1") in background and return immediately
// p.Call("2")			// run f("2") in background after block <Pool.delay> ns
// p.Call("3")			// run f("3") in background after block <Pool.delay> ns
//
type Pool struct {
	// the current number of coroutines
	// put in head to make <cur> 64-bit aligned
	cur int64
	// at least <min> coroutines will be created and live all the time
	min uint32
	// the max number of coroutines can be created
	max uint32
	// create a new coroutine if <arg> is blocked for <delay> ns
	// a proper value should be >= 0.1s
	delay time.Duration
	// destroy the coroutines which idle for <timeout> ns
	timeout time.Duration

	fun func(interface{})
	arg chan interface{}
}

func NewPool(min, max uint, delay, timeout time.Duration, fun func(interface{})) (v *Pool) {
	if fun == nil {
		golib.Panic("fun == nil, want !nil")
	}

	v = &Pool{
		fun:     fun,
		arg:     make(chan interface{}),
		min:     uint32(min),
		cur:     0,
		max:     uint32(max),
		delay:   delay,
		timeout: timeout,
	}
	for v.cur < int64(v.min) {
		v.newProcess()
	}
	return
}

func PoolWrap(fun func(interface{})) (func(interface{}) /* stop */, func()) {
	v := NewPool(_POOL_MIN, _POOL_MAX, _POOL_DELAY, _POOL_TIMEOUT, fun)
	return v.Call,
		func() {
			v.WithCount(0, 0)
			v.WithTime(v.delay, 0)
		}
}

// The created coroutine won't quit unless time out. Set min to 0 if want to quit all.
func (o *Pool) newProcess() {
	// XXX: o.cur may overflow if parallel call exceed math.MaxInt64
	if atomic.AddInt64(&o.cur, 1) > int64(o.max) {
		// no coroutine created, restore the value
		atomic.AddInt64(&o.cur, -1)
		return
	}

	var loop *Loop
	loop = NewLoop(func() {
		select {
		case arg := <-o.arg:
			// the outer select: skip creating timer when busy
			o.fun(arg)
		default:
			select {
			case arg := <-o.arg:
				o.fun(arg)
			case <-time.After(o.timeout): // quit if idle for <timeout> ns
				if atomic.AddInt64(&o.cur, -1) >= int64(o.min) {
					loop.Stop()
				} else {
					// coroutine isn't killed, restore the value
					// XXX: o.cur may > o.max if o.cur increased by newProcess() at the same time
					atomic.AddInt64(&o.cur, 1)
				}
			}
		}
	})
}

func (o *Pool) Call(arg interface{}) {
	if o.cur >= int64(o.max) {
		// skip creating timer, when busy and no more goroutine can be created
		o.arg <- arg
	} else {
		select {
		case o.arg <- arg:
		case <-time.After(o.delay):
			// If <delay> is too small, select may choose this case even <o.arg> isn't blocked.
			o.newProcess()
			o.Call(arg)
		}
	}
}

// Set when to create or kill a goroutine.
// A new goroutine will be created after the argument blocked for ^delay ns.
// A goroutine will be killed after idle for ^timeout ns
func (o *Pool) WithTime(delay time.Duration, timeout time.Duration) *Pool {
	o.delay = delay
	o.timeout = timeout
	return o
}

// Change how many goroutines the Pool can create, ^min <= count <= ^max.
func (o *Pool) WithCount(min uint, max uint) *Pool {
	if min > max {
		golib.Panic("min:%v > max:%v !", min, max)
	}

	o.min = uint32(min)
	o.max = uint32(max)
	for o.cur < int64(o.min) {
		o.newProcess()
	}
	return o
}
