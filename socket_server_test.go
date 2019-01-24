package gosvc

//func TestServeMixin_Serve(t *testing.T) {
//	listener, err := net.Listen("tcp", ":8080")
//	if err != nil {
//		t.Fatalf("%v", err)
//	}
//
//	serv := NewLoop(func() {
//		t.Logf("%v", ServeMixin{
//			Listener: listener,
//		}.Serve(func(conn net.Conn) {
//			//t.Logf("%v", conn.RemoteAddr())
//		}))
//	})
//
//	//NewLoop(func() {
//	//	time.Sleep(100 * time.Millisecond)
//	//	for {
//	//		NewLoop(func() {
//	//			_, err := net.Dial("tcp", "0.0.0.0:8080")
//	//			if err != nil {
//	//				t.Logf("%v: %v", time.Now(), err)
//	//			}
//	//		})
//	//	}
//	//})
//	time.Sleep(time.Second)
//	listener.Close()
//	serv.Stop()
//	serv.Wait()
//}
