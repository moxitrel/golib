package gosvc

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

/* Example

listener, err := net.Listen("tcp", ":0")
if err != nil {
	log.Fatalln(err)
	return
}

srv := ServeMixin{
	Listener: listener,
}
go srv.Serve(func(bytes []byte, conn net.Conn) int {
	// handle recv bytes
	return 0
})

*/
type ServeMixin struct {
	net.Listener

	Timeout    time.Duration // disconnect client when timeout
	BufferSize int           // set recv buffer size

	bufferPool *sync.Pool // recv buffer
	handler    *Pool
}

const _DEFAULT_BUFFER_SIZE = 1 << 16

var gBufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, _DEFAULT_BUFFER_SIZE)
	},
}

func logIfError(err error) {
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
		o.bufferPool = &gBufferPool
	} else {
		o.bufferPool = &sync.Pool{
			New: func() interface{} {
				return make([]byte, o.BufferSize)
			},
		}
	}

	o.handler = NewPool(3, 1<<20, 3*time.Minute, func(arg interface{}) {
		conn := arg.(net.Conn)
		defer func() {
			logIfError(conn.Close())
		}()

		//
		// handle connection
		//
		buffer := o.bufferPool.Get().([]byte)
		defer o.bufferPool.Put(buffer)

		for i := 0; 0 <= i && i <= o.BufferSize; {
			logIfError(conn.SetReadDeadline(time.Now().Add(o.Timeout)))
			n, err := conn.Read(buffer[i:])
			logIfError(conn.SetReadDeadline(time.Time{}))
			if n > 0 {
				i = callback(buffer[:i+n], conn)
			}

			//
			// handle read error
			//
			if err == io.EOF { // connection closed
				i = -1
			} else if err, ok := err.(net.Error); ok && err.Timeout() { // timeout
				i = -1
			} else {
				i = -1
				logIfError(err)
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
