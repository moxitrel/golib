package golib

import (
	"io"
	"net"
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

// Receive bytes from reader until ok() returns true.
func ReadUnitl(reader io.Reader, ok func([]byte) bool) (buffer []byte, err error) {
	if reader == nil {
		Panic("reader = nil, want !nil")
	}
	if ok == nil {
		Panic("ok = nil, want !nil")
	}
	buffer = BytesPool.Get(1024)
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

// remoteAddr: e.g. "192.168.0.1:8080"
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
