package golib

import (
	"testing"
	"time"
)

func TestWriteAll(t *testing.T) {
	// panic if writer = nil
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("err = nil; want !nil")
		} else {
			t.Logf("err: %v", err)
		}
	}()
	WriteAll(nil, nil)
}

func TestTcpOnce(t *testing.T) {
	TcpOnce(
		"baidu.com:80",
		[]byte("GET /index.htm HTTP/1.1\r\nContent-Length: 0\r\n\r\n"),
		func(rsp []byte) bool {
			t.Logf("%s", rsp)
			return true
		},
		time.Minute)
}
