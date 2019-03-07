package gosvc

import (
	"net"
	"testing"
	"time"
)

func TestServeMixin_Serve(t *testing.T) {
	t.Skipf("affect other test\n")

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
		return
	}

	srv := ServeMixin{
		Listener: listener,
	}

	NewLoop(func() {
		for i := 0; i < 5000; i++ {
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

	go func () {
		time.Sleep(100 * time.Millisecond)
		listener.Close()
	}()
	srv.Serve(func(bytes []byte, conn net.Conn) int {
		// handle recv bytes
		return 0
	})
}
