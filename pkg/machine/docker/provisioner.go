package docker

import (
	machine "github.com/footprintai/multikf/pkg/machine"
)

var docker machine.Provisioner = "docker"

func init() {
	machine.RegisterProvisioner(docker, NewHostMachines)
}
