package gosvc

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type ServeMixin struct {
	net.Listener

	Timeout    time.Duration // disconnect when timeout, default 1m
	BufferSize int

	bufferPool sync.Pool
	handler    *Pool
}

const _DEFAULT_BUFFER_SIZE = 1 << 16

var gBufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, _DEFAULT_BUFFER_SIZE)
	},
}

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (o *ServeMixin) Serve(callback func([]byte, net.Conn) int) error {
	if o.Listener == nil {
		panic(fmt.Errorf("Listener == nil, want !nil"))
	}
	if o.Timeout == 0 {
		o.Timeout = time.Minute
	}
	if o.BufferSize <= 0 {
		o.BufferSize = _DEFAULT_BUFFER_SIZE
	}

	if o.BufferSize == _DEFAULT_BUFFER_SIZE {
		o.bufferPool = gBufferPool
	} else {
		o.bufferPool.New = func() interface{} {
			return make([]byte, o.BufferSize)
		}
	}

	o.handler = NewPool(3, 1<<23, 3*time.Minute, func(arg interface{}) {
		conn := arg.(net.Conn)
		defer func() {
			logError(conn.Close())
		}()

		//
		// handle connection
		//
		buffer := o.bufferPool.Get().([]byte)
		defer o.bufferPool.Put(buffer)

		for i := 0; 0 <= i && i <= o.BufferSize; {
			logError(conn.SetReadDeadline(time.Now().Add(o.Timeout)))
			n, err := conn.Read(buffer[i:])
			logError(conn.SetReadDeadline(time.Time{}))
			if n > 0 {
				i = callback(buffer[:i+n], conn)
			}
			if err != nil {
				logError(err)
				i = -1
			}
		}
	})

	for {
		conn, err := o.Listener.Accept()
		if err != nil {
			o.handler.Stop()
			return err
		}
		o.handler.Call(conn)
	}
}

// Send all bytes in data to writer.
func WriteAll(writer io.Writer, data []byte) error {
	if writer == nil {
		panic(fmt.Errorf("writer == nil, want !nil"))
	}
	for len(data) > 0 {
		n, err := writer.Write(data)
		if err != nil {
			return err
		}
		data = data[n:]
	}
	return nil
}
