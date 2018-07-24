package golib

import (
	"io"
	"net"
	"time"
)

// Send all bytes in data to writer.
func WriteAll(writer io.Writer, data []byte) error {
	if writer == nil {
		Panic("writer = nil, want !nil")
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

// Receive bytes from reader until ok() return true.
func ReadUnitl(reader io.Reader, ok func([]byte) bool) (buffer []byte, err error) {
	if reader == nil {
		Panic("reader = nil, want !nil")
	}
	if ok == nil {
		Panic("ok = nil, want !nil")
	}
	buffer = BytesPool.Get()
	// receive response unitl callback success or timeout
	var n int
	for i := 0; ; {
		n, err = reader.Read(buffer[i:])
		i += n
		if n > 0 && ok(buffer[:i]) {
			// success
			buffer, err = buffer[:i], nil
			break
		}
		if err != nil {
			break
		}
		if i == len(buffer) { // buffer is full
			buffer = append(buffer, make([]byte, len(buffer))...)
		}
	}
	return
}

// Send req, and receive response until cb() return true.
func WriteCb(rw io.ReadWriter, req []byte, cb func(io.ReadWriter, []byte) bool) (err error) {
	if rw == nil {
		Panic("rw = nil, want !nil")
	}
	if len(req) == 0 {
		return nil
	}

	// send request
	err = WriteAll(rw, req)
	if err != nil {
		return
	}

	// handle response
	if cb == nil {
		return
	}
	data, err := ReadUnitl(rw, func(buffer []byte) bool {
		return cb(rw, buffer)
	})

	BytesPool.Put(data)
	return
}

// Send one request, and receive the response on TCP
// remoteAddr: e.g. "192.168.0.1:8080"
// cb: handle response; if return false, continue receiving response data; if return true, quit
func WithTcpWrite(remoteAddr string, sentData []byte, cb func(net.Conn, []byte) bool, timeout time.Duration) (err error) {
	if len(sentData) == 0 {
		return nil
	}

	// connect
	conn, err := net.DialTimeout("tcp", remoteAddr, timeout)
	if err != nil {
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(timeout))

	err = WriteCb(conn, sentData, func(_ io.ReadWriter, bytes []byte) bool {
		return cb(conn, bytes)
	})

	return
}

func WithTcp(remoteAddr string, cb func(net.Conn)) (err error) {
	if cb == nil {
		return
	}

	// connect
	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	cb(conn)

	return
}
