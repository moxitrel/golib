package golib

import (
	"fmt"
	"io"
	"net"
	"time"
)

func WriteAll(writer io.Writer, data []byte) error {
	if writer == nil {
		Panic(fmt.Sprintf("writer = nil, want !nil", ))
	}
	if len(data) == 0 {
		return nil
	}

	n, err := writer.Write(data)
	if err != nil {
		return err
	}
	if n < len(data) {
		return WriteAll(writer, data[n:])
	}

	return nil
}

// Send one request, and receive the response on TCP
// remoteAddr: e.g. "192.168.0.1:8080"
// cb: handle response; if return false, continue receiving response data; if return true, quit
func TcpOnce(remoteAddr string, sentData []byte, cb func([]byte) bool, timeout time.Duration) error {
	if len(sentData) == 0 || cb == nil {
		return nil // NOTE: may use panic() instead
	}

	// connect
	conn, err := net.DialTimeout("tcp", remoteAddr, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(timeout))

	// send request
	err = WriteAll(conn, sentData)
	if err != nil {
		return err
	}

	// receive response unitl callback success or timeout
	buffer := BytesPool.Get()
	defer func() {
		BytesPool.Put(buffer)
	}()
	for i := 0;; {
		n, err := conn.Read(buffer[i:])
		i += n
		if n > 0 && cb(buffer[:i]) {
			// success
			err = nil
			break
		}
		if err != nil {
			// err: io.EOF | net.OpError.Timeout() | ...
			break
		}
		if i+1 == len(buffer) { // buffer is full
			buffer = append(buffer, make([]byte, 1024)...)
		}
	}
	return err
}
