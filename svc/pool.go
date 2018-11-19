/*

func         PoolWrapper (func (interface{})                     ) *Pool
func (*Pool) WithCount(min   uint         , max     uint         ) *Pool
func (*Pool) WithTime (delay time.Duration, timeout time.Duration) *Pool

func (*Pool) Call     (interface{})

*/
package svc

import (
	"github.com/moxitrel/golib"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// Wrap time.Timer
type Timer struct {
	*time.Timer
}

func NewTimer() (o Timer) {
	o = Timer{
		Timer: time.NewTimer(-1),
	}
	o.Stop()
	return
}
func (o Timer) Start(timeout time.Duration) {
	o.Reset(timeout)
}
func (o Timer) Stop() {
	if !o.Timer.Stop() {
		<-o.C
	}
}

const (
	_POOL_MIN     = 0
	_POOL_MAX     = math.MaxUint16
	_POOL_DELAY   = 200 * time.Millisecond
	_POOL_TIMEOUT = time.Minute

	// time to wait for receiving sent args when receive stop signal
	_STOP_DELAY = 200 * time.Millisecond
)

// Tell goroutine to exit
type _StopSignal struct{}

var stopSignal = _StopSignal{}

// Start [min, max] goroutines of <Pool.fun> to process <Pool.arg>
//
// * Example
// f := func(x interface{}) { time.Sleep(time.Second) }
// p := PoolWrapper(f)	    // start 1 goroutines of f
// p.Call("1")			// run f("1") in background and return immediately
// p.Call("2")			// run f("2") in background after block <Pool.delay> ns
// p.Call("3")			// run f("3") in background after block <Pool.delay> ns
//
type Pool struct {
	// the current number of goroutines
	// NOTE: put in head to make <cur> 64-bit aligned
	cur int64
	// at least <min> goroutines will be created and live all the time
	min uint32
	// the max number of goroutines can be created
	max uint32
	// create a new goroutine if <arg> is blocked for <delay> ns
	// a proper value should be >= 0.1s
	delay time.Duration
	// destroy the goroutines which idle for <timeout> ns
	timeout time.Duration

	fun func(interface{})
	arg chan interface{}

	// lock cur when updating
	curLock sync.Locker

	wg sync.WaitGroup
}

func (o *Pool) getCur() int64 {
	return atomic.LoadInt64(&o.cur)
}
func (o *Pool) getMin() int64 {
	return int64(atomic.LoadUint32(&o.min))
}
func (o *Pool) setMin(min uint32) {
	atomic.StoreUint32(&o.min, min)
}
func (o *Pool) getMax() int64 {
	return int64(atomic.LoadUint32(&o.max))
}
func (o *Pool) setMax(max uint32) {
	atomic.StoreUint32(&o.max, max)
}
func (o *Pool) getDelay() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&o.delay)))
}
func (o *Pool) setDelay(timeout time.Duration) {
	atomic.StoreInt64((*int64)(&o.delay), int64(timeout))
}
func (o *Pool) getTimeout() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&o.timeout)))
}
func (o *Pool) setTimeout(timeout time.Duration) {
	atomic.StoreInt64((*int64)(&o.timeout), int64(timeout))
}

func NewPool(min, max uint, delay, timeout time.Duration, bufferSize uint, fun func(interface{})) (o *Pool) {
	if min > max {
		golib.Panic("min:%v > max:%v, want min <= max", min, max)
	}
	if fun == nil {
		golib.Panic("fun == nil, want !nil")
	}

	o = &Pool{
		cur:     0,
		min:     uint32(min),
		max:     uint32(max),
		delay:   delay,
		timeout: timeout,
		fun:     fun,
		arg:     make(chan interface{}, bufferSize),
		curLock: &sync.Mutex{},
		wg:      sync.WaitGroup{},
	}

	// if timeout is 0, new process will exit immediately which decrease the cur.
	for i := uint32(0); i < o.min; i++ {
		o.newProcess()
	}
	return
}

func PoolWrapper(fun func(interface{})) (func() func(interface{}), func()) {
	v := NewPool(_POOL_MIN, _POOL_MAX, _POOL_DELAY, _POOL_TIMEOUT, 0, fun)
	return v.Submitter, v.Stop
}

// Create a new goroutine.
func (o *Pool) newProcess() {
	// NOTE: Additionally use atomic operator for o.cur to avoid data race from .Submitter() which has no lock.
	//       Only lock o.cur when updating.
	//
	// XXX: It's not reliable in theory, but ok in practice.
	//if atomic.AddInt64(&o.cur, 1) > o.getMax() {
	//	atomic.AddInt64(&o.cur, -1)	// no goroutine created
	//	return
	//}
	o.curLock.Lock()
	if o.getCur() >= o.getMax() {
		o.curLock.Unlock()
		return
	}
	atomic.AddInt64(&o.cur, 1)
	o.curLock.Unlock()

	timeoutTimer := NewTimer()
	o.wg.Add(1)
	go func() {
		defer o.wg.Done()
		var arg interface{}
		for {
			select {
			case arg = <-o.arg: // skip creating timer if not blocked
				goto HANDLE_ARG
			default:
				switch timeout := o.getTimeout(); {
				case timeout < 0: // skip creating timer if wait forever
					arg = <-o.arg
					goto HANDLE_ARG
				default:
					timeoutTimer.Start(timeout)
					select {
					case arg = <-o.arg:
						timeoutTimer.Stop()
						goto HANDLE_ARG
					case <-timeoutTimer.C: // timeout, idle too long
						o.curLock.Lock()
						if o.getCur() > o.getMin() {
							atomic.AddInt64(&o.cur, -1)
							o.curLock.Unlock()
							return
						}
						o.curLock.Unlock()
					}
				}
			}
		HANDLE_ARG:
			switch arg {
			case stopSignal: // try to send stop-signal to another goroutine if any alive
				timeoutTimer.Start(o.getTimeout())
				select {
				case o.arg <- stopSignal:
					timeoutTimer.Stop()
				case <-timeoutTimer.C:
				}
				return
			default:
				o.fun(arg)
			}
		}
	}()
}

func (o *Pool) Stop() {
	if o.getMax() > 0 {
		o.setMin(0)
		o.setMax(0)
		o.setTimeout(_STOP_DELAY) // fetch appending args
		o.arg <- stopSignal
	}
}

func (o *Pool) Join() {
	o.wg.Wait()
}

func (o *Pool) Submitter() func(interface{}) {
	delayTimer := NewTimer()
	return func(arg interface{}) {
	RESUBMIT:
		switch max := o.getMax(); {
		case max <= 0: // stopped?
			// return
		case o.getCur() >= max: // wait forever
			// NOTE: use atomic instead of lock for o.cur for performance
			o.arg <- arg
		default:
			select {
			case o.arg <- arg: // skip creating timer if not blocked
				// return
			default:
				switch delay := o.getDelay(); {
				case delay == 0: // skip creating timer if delay = 0
					o.newProcess()
					goto RESUBMIT
				case delay < 0: // wait forever
					o.arg <- arg
				default:
					delayTimer.Start(delay)
					select {
					case o.arg <- arg:
						delayTimer.Stop()
					case <-delayTimer.C:
						// NOTE: If <delay> is too small, select may choose this case even <o.arg> isn't blocked,
						//       which may be caused by gc.  A proper value should be >= 0.1s in my test.
						o.newProcess()
						goto RESUBMIT
					}
				}
			}
		}
	}
}

// Set when to create or kill a goroutine.
// A new goroutine will be created after the argument blocked for <delay> ns.
// A goroutine will be killed after idle for <timeout> ns
func (o *Pool) SetTime(delay time.Duration, timeout time.Duration) {
	o.setDelay(delay)
	o.setTimeout(timeout)
}

// Change how many goroutines the Pool can create.
func (o *Pool) SetCount(min uint, max uint) {
	if min > max {
		golib.Panic("min:%v > max:%v, want min <= max", min, max)
	}

	o.setMin(uint32(min))
	o.setMax(uint32(max))

	n := o.getCur() - o.getMin()
	for i := int64(0); i < n; i++ {
		o.newProcess()
	}
}
