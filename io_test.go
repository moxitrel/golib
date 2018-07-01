package golib

import (
	"testing"
	"time"
)

func TestTcpOnce(t *testing.T) {
	TcpOnce(
		"baidu.com:80",
		[]byte("GET /index.htm HTTP/1.1\r\nContent-Length: 0\r\n\r\n"),
		time.Minute,
		func(rsp []byte) bool {
			t.Logf("%s", rsp)
			return true
		})
}
