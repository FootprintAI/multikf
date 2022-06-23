package docker

import (
	"encoding/json"

	machinecmd "github.com/footprintai/multikf/pkg/machine/cmd"
	"github.com/footprintai/multikf/pkg/machine/ioutil"
	"github.com/go-cmd/cmd"
	"sigs.k8s.io/kind/pkg/log"
)

func NewDockerCli(logger log.Logger, verbose bool) (*DockerCli, error) {
	return &DockerCli{logger: logger, verbose: verbose}, nil
}

type DockerCli struct {
	logger  log.Logger
	verbose bool
}

type dockerState struct {
	Status string `json:"status"`
}

func (cli *DockerCli) GetClusterStatus(containername ContainerName) (string, error) {
	cmdAndArgs := []string{
		"docker",
		"inspect",
		containername.Name(),
		"--format='{{json .State}}'",
	}
	sr, _, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return "", err
	}
	d := dockerState{}
	blob, _ := ioutil.ReadAll(sr)
	stripped := blob[1 : len(blob)-2] // remove ' xxx '\n
	if err := json.Unmarshal(stripped, &d); err != nil {
		return "", err
	}
	return d.Status, nil
}

func (cli *DockerCli) RemoteExec(containername ContainerName, cmd string) (resp string, err error) {
	cmdAndArgs := []string{
		"docker",
		"exec",
		containername.Name(),
		"sh",
		"-c",
		cmd,
	}
	sr, _, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return "", err
	}
	all, _ := ioutil.ReadAll(sr)
	return string(all), nil
}

func (cli *DockerCli) runCmd(cmdAndArgs []string) (ioutil.StreamReader, <-chan cmd.Status, error) {
	return machinecmd.NewCmd(cli.logger, cli.verbose).Run(cmdAndArgs...)
}
