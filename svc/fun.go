/*

NewFun f		: "finish f(Any) added by Call()"
	Call arg	: "schedule f(arg)"
	Stop 		:

*** e.g.

// 1. define a new type derive Fun
type T struct {
	Fun
}

// 2. define construction
func NewF() *T {
	// 2.1. define a function with signature func(interface{})
	f := func (arg interface{}) {
		x := arg.(ArgT)		//recover type first
		...					//do the things
	}

	// 2.2. return
	return &F{ *NewFun(f) }
}

// 3. override Call() with desired argument type
func (o *T) Call(x ArgT) {
	o.Fun.Call(x)
}

*/
package svc

import (
	"sync"
	"fmt"
	"errors"
)

type Fun struct {
	fun       func(interface{})
	args      chan interface{}
	startOnce sync.Once
	stopOnce  sync.Once
}

// Return a started fun-service
// fun: apply with arg passed from Call()
func NewFun(fun func(arg interface{})) (v *Fun) {
	if fun == nil {
		panic(errors.New(fmt.Sprintf("fun = nil, want non nil")))
	}
	v = &Fun{
		fun: fun,
	}
	v.Start()
	return
}

func (o *Fun) Call(arg interface{}) {
	o.args <- arg
}

func (o *Fun) Start() {
	o.startOnce.Do(func(){
		o.args = make(chan interface{}, FUN_BUFFER_SIZE)
		o.stopOnce = sync.Once{}
		go func(){
			for arg := range o.args {
				o.fun(arg)
			}
			o.startOnce = sync.Once{}
		}()
	})
}

func (o *Fun) Stop() {
	o.stopOnce.Do(func() {
		close(o.args)
	})
}
