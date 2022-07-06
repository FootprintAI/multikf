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

func (c *Cmd) Run(cmdAndArgs ...string) (*ioutil.CmdOutputStream, <-chan cmd.Status, error) {
	c.logger.V(1).Infof("cmd->%s\n", cmdAndArgs)
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	runcmd := cmd.NewCmdOptions(cmdOptions, cmdAndArgs[0], cmdAndArgs[1:]...)
	statusChan1 := make(chan cmd.Status, 1)
	statusChan2 := make(chan cmd.Status, 1)
	go newChanForwarder(runcmd.Start(), statusChan1, statusChan2)

	cinfo := &ioutil.CommandInfo{
		CommandStatus: statusChan1,
		Command:       runcmd,
	}

	return ioutil.NewCmdOutputStream(c.logger, cinfo), statusChan2, nil
}

func newChanForwarder(src <-chan cmd.Status, dests ...chan<- cmd.Status) {
	status := <-src
	for _, dest := range dests {
		dest <- status
	}
	return
}
