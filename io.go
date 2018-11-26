package golib

import (
	"fmt"
	"io"
	"net"
)

// Send all bytes in data to writer.
func WriteAll(writer io.Writer, data []byte) error {
	if writer == nil {
		Panic("writer == nil, want !nil")
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

// Read bytes from reader until ok() returns true.
func ReadUnitl(reader io.Reader, buffer []byte, ok func([]byte) bool) (err error) {
	if reader == nil {
		Panic("reader == nil, want !nil")
	}
	if len(buffer) < 1 {
		Panic("buffer.Size == 0, want > 0")
	}
	if ok == nil {
		Panic("ok == nil, want !nil")
	}
	// receive response unitl callback success or timeout
	var n int
	for i := 0; ; {
		n, err = reader.Read(buffer[i:])
		i += n
		if n > 0 && ok(buffer[:i]) {
			// success
			err = nil
			break
		}
		if err != nil {
			break
		}
		if i == len(buffer) {
			err = fmt.Errorf("buffer is full")
			break
		}
	}
	return
}

// remoteAddr: e.g. "192.168.0.1:8080"
func WithTcp(remoteAddr string, cb func(net.Conn) error) error {
	if cb == nil {
		return nil
	}

	// connect
	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	return cb(conn)
}

func WithUdp(localPort uint16, cb func(net.PacketConn) error) error {
	if cb == nil {
		return nil
	}

	// connect
	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%v", localPort))
	if err != nil {
		return err
	}
	defer conn.Close()

	return cb(conn)
}
