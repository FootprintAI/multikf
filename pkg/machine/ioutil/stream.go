package ioutil

import (
	"bytes"
	"fmt"
	"io"

	"github.com/go-cmd/cmd"
	"sigs.k8s.io/kind/pkg/log"
)

type StreamReader interface {
	Read(b []byte) (int, error) // read stream to buffer with Reader
}

// CmdOutputStream implements io.Reader by wrapping the line channel
type CmdOutputStream struct {
	logger log.Logger
	cInfo  *CommandInfo
}

type CommandInfo struct {
	CommandStatus <-chan cmd.Status
	Command       *cmd.Cmd
}

func NewCmdOutputStream(logger log.Logger, cInfo *CommandInfo) *CmdOutputStream {
	return &CmdOutputStream{cInfo: cInfo, logger: logger}
}

func (o *CmdOutputStream) Read(b []byte) (int, error) {
	out, more := <-o.cInfo.Command.Stdout
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

func StderrOnError(o *CmdOutputStream) error {

	go func() {
		status := <-o.cInfo.CommandStatus
		if status.Exit != 0 {
			// process exit abnormally, display stderr
			for lineLog := range o.cInfo.Command.Stderr {
				o.logger.V(0).Infof("%s\n", lineLog)
			}
		}
	}()

	for lineLog := range o.cInfo.Command.Stdout {
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
