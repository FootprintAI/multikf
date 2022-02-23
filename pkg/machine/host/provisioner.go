package host

import (
	machine "github.com/footprintai/multikind/pkg/machine"
)

var docker machine.Provisioner = "docker"

func init() {
	machine.RegisterProvisioner(docker, NewHostMachines)
}
