package golib

import "sync"

type _BytesPool struct {
	sync.Pool
}

func (o *_BytesPool) Get(n int) (v []byte) {
	v = o.Pool.Get().([]byte)
	if n <= cap(v) {
		v = v[:n]
	} else {
		o.Put(v)
		v = make([]byte, n)
	}
	return
}

func (o *_BytesPool) Put(x []byte) {
	o.Pool.Put(x)
}

var BytesPool = _BytesPool{
	Pool: sync.Pool{
		New: func() interface{} {
			return make([]byte, 0)
		},
	},
}
