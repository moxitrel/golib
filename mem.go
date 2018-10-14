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
	if v != nil && uint(cap(v)) >= n {
		v = v[:n]
	} else {
		// throw the item (no Put(v)),
		// or you may always get a []byte whose cap() is always < n
		v = make([]byte, n)
	}
	return
}

func (o *_BytesPool) Put(x []byte) {
	o.Pool.Put(x)
}
