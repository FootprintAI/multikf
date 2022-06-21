package ioutil

import (
	"bytes"
	"fmt"
	"io"

	"sigs.k8s.io/kind/pkg/log"
)

type StreamReader interface {
	Stdout() error              // stream outputs to stdout
	Read(b []byte) (int, error) // read stream to buffer with Reader
}

// outputStream implements io.Reader by wrapping the line channel
type outputStream struct {
	logger log.Logger
	ch     chan string
}

func NewOutputStream(logger log.Logger, ch chan string) *outputStream {
	return &outputStream{ch: ch, logger: logger}
}

func (o *outputStream) Read(b []byte) (int, error) {
	out, more := <-o.ch
	if !more {
		return 0, io.EOF
	}
	if len(out) > len(b) {
		panic(fmt.Sprintf("insufficient buffer size(buf:%d, data:%d), data could be lost", len(b), len(out)))
	}
	n := copy(b[:len(b)-1], []byte(out))
	b[n] = '\n'
	return n + 1, nil
}

func (o *outputStream) Stdout() error {
	for lineLog := range o.ch {
		o.logger.V(1).Infof("%s\n", lineLog)
	}
	return nil
}

func ReadAll(r io.Reader) ([]byte, error) {
	buf := &bytes.Buffer{}
	b := make([]byte, 1024*1024*10 /*10M buffer*/)
	for {
		n, err := r.Read(b)
		if err == nil {
			buf.Write(b[:n])
		} else {
			if err == io.EOF {
				return buf.Bytes(), nil
			} else {
				return nil, err
			}
		}
	}
}
