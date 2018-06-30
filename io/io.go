package io

import (
	"io"
	"errors"
	"fmt"
	"github.com/moxitrel/golib/log"
)

func WriteAll(writer io.Writer, data []byte) error {
	if writer == nil {
		panic(errors.New(fmt.Sprintf("%v: writer = %v, want non nil!", log.CallerName(0), writer)))
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
