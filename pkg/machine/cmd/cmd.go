package cmd

import (
	"github.com/go-cmd/cmd"
	"sigs.k8s.io/kind/pkg/log"

	"github.com/footprintai/multikf/pkg/machine/ioutil"
)

func NewCmd(logger log.Logger) *Cmd {
	return &Cmd{
		logger: logger,
	}
}

type Cmd struct {
	logger log.Logger
}

func (c *Cmd) Run(cmdAndArgs ...string) (ioutil.StreamReader, <-chan cmd.Status, error) {
	c.logger.V(1).Infof("cmd->%s\n", cmdAndArgs)
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

	return ioutil.NewOutputStream(c.logger, runcmd.Stdout), status, nil
}
