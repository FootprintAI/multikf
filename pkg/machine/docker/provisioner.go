package docker

import (
	machine "github.com/footprintai/multikf/pkg/machine"
)

const docker machine.Provisioner = "docker"

func init() {
	machine.RegisterProvisioner(docker, NewHostMachines)
}
