package golib

import "sync"

type _BytesPool struct {
	sync.Pool
}

var BytesPool = _BytesPool{
	Pool: sync.Pool{
		New: func() interface{} {
			return []byte(nil)
		},
	},
}

func (o *_BytesPool) Get(n uint) (v []byte) {
	v = o.Pool.Get().([]byte)
	if v == nil {
		v = make([]byte, n)
	} else if uint(cap(v)) >= n {
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
