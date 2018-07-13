package golib

import "sync"

type _BytesPool struct {
	sync.Pool
}
func (o *_BytesPool) Get() []byte {
	return o.Pool.Get().([]byte)
}
func (o *_BytesPool) Put(x []byte) {
	o.Pool.Put(x)
}

var BytesPool = _BytesPool{
	Pool: sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	},
}
