/*

func(   ): no  arg, use Service
func(...): has arg, use Fun

*** e.g.

func f(x T) {
	...
}

type F struct {
	Fun
}

func NewF() (v *F) {
	v = &F{
		Fun: *NewFun(func(x interface{}) {
			f(x.(T))
		}),
	}
	return
}

func (o *PrintBytes) Call(x T) {
	o.Fun.Call(x)
}

*/
package svc

type Fun struct {
	Service
	args chan interface{}
}

func NewFun(f func(interface{})) (v *Fun) {
	v = &Fun{
		args: make(chan interface{}, 1024), //todo: specify 1024 buffer size
	}

	v.Service = *New(func() {
		arg := <-v.args
		f(arg)
	})
	return
}

func (o *Fun) Call(arg interface{}) {
	o.args <- arg
}
