/*

NewFun (f & ([Any]:)): "finish f(Any ...) added by Call()"
	Call Any ...	 : "run f(Any ...) once"
	Start			 :
	Stop 			 :

*** e.g.

type F struct {
	Fun
}

func NewF() (v *F) {
	v = &F{
		Fun: *NewFun(func(argv []interface{}) {
			x1 := argv[0].(T)	//1. convert type
			f(x)				//2. do the things
		}),
	}
	return
}

func (o *F) Call(x T) {
	o.Fun.Call(x)
}

*/
package svc

type Fun struct {
	Service
	argvs chan []interface{} //a service usually has a buffer
	//stopRead chan struct{}		//signal <-argvs to quit if blocked
}

// f   : 1. convert type; 2. do the things
// argv: passed from Call()
func NewFun(f func(argv []interface{})) (v *Fun) {
	v = &Fun{
		argvs: make(chan []interface{}, 100*10000), //todo: specify buffer size
	}

	v.Service = *New(func() {
		//select {
		//case argv := <-v.argvs:
		//	f(argv)
		//case <-v.stopRead:
		//	//nop
		//}
		argv := <-v.argvs
		f(argv)
	})
	return
}

func (o *Fun) Call(argv ...interface{}) {
	o.argvs <- argv
}

//func (o *Fun) Start() {
//	// clear stopRead
//	select {
//	case <-o.stopRead:
//		//nop
//	default:
//		//nop
//	}
//
//	o.Service.Stop()
//}
//
//func (o *Fun) Stop() {
//	select {
//	case o.stopRead <- struct{}{}:
//		//nop
//	default:
//		//nop
//	}
//
//	o.Service.Stop()
//}
