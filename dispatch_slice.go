package golib

// not thread-safe, no deletion
type SliceDispatch struct {
	pool []func(interface{})
}

func MakeSliceDispatch(size uint) SliceDispatch {
	return SliceDispatch{
		pool: make([]func(interface{}), 0, size),
	}
}

func (o SliceDispatch) Add(handler func(interface{})) (index int) {
	if handler == nil {
		Panic("handler == nil, want !nil")
	}
	index = len(o.pool)
	o.pool = append(o.pool, handler)
	return
}

func (o SliceDispatch) Set(index int, handler func(interface{})) {
	if index >= cap(o.pool) {
		Panic("index:%v is out of range:%v", index, len(o.pool))
	}
	if handler == nil {
		Panic("handler == nil, want !nil")
	}
	if index > len(o.pool) {
		o.pool = o.pool[:index+1]
	}
	o.pool[index] = handler
	return
}

func (o SliceDispatch) Call(index int, arg interface{}) {
	if index >= len(o.pool) {
		Panic("index:%v is out of range:%v", index, len(o.pool))
	}
	handler := o.pool[index]
	if handler == nil {
		Panic("index:%v, handler doesn't exist", index)
	}
	handler(arg)
}
