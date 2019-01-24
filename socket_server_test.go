package gosvc

import (
	"net"
	"testing"
	"time"
)

func TestServeMixin_Serve(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
		return
	}

	srv := ServeMixin{
		Listener: listener,
	}
	go srv.Serve(func(bytes []byte, conn net.Conn) int {
		// handle recv bytes
		return 0
	})

	NewLoop(func() {
		for {
			NewLoop(func() {
				conn, err := net.Dial("tcp", listener.Addr().String())
				if err != nil {
					t.Logf("%v: %v", time.Now(), err)
					return
				}
				logIfError(WriteAll(conn, []byte("hi")))
			})
		}
	})
	time.Sleep(100 * time.Millisecond)
	listener.Close()
}
