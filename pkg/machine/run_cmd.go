package machine

import (
	"fmt"
	"io"

	"github.com/go-cmd/cmd"
	"sigs.k8s.io/kind/pkg/log"
)

func NewCmd(logger log.Logger, verbose bool) *Cmd {
	return &Cmd{
		logger:  logger,
		verbose: verbose,
	}
}

type Cmd struct {
	logger  log.Logger
	verbose bool
}

func (c *Cmd) Run(cmdAndArgs ...string) (StreamReader, <-chan cmd.Status, error) {
	if c.verbose {
		c.logger.V(1).Infof("cmd->%s\n", cmdAndArgs)
	}
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	runcmd := cmd.NewCmdOptions(cmdOptions, cmdAndArgs[0], cmdAndArgs[1:]...)
	status := runcmd.Start()
	// run and output stderr
	for stderrline := range runcmd.Stderr {
		c.logger.V(1).Infof("%s\n", stderrline)
	}
	//stat := <-runStatus

	return newOutputStream(c.logger, runcmd.Stdout), status, nil
}

type StreamReader interface {
	Stdout() error              // stream outputs to stdout
	Read(b []byte) (int, error) // read stream to buffer with Reader
}

// outputStream implements io.Reader by wrapping the line channel
type outputStream struct {
	logger log.Logger
	ch     chan string
}

func newOutputStream(logger log.Logger, ch chan string) *outputStream {
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
		o.logger.V(0).Infof("%s\n", lineLog)
	}
	return nil
}
