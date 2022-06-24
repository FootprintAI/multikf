package vagrant

import (
	machine "github.com/footprintai/multikf/pkg/machine"
)

const vagrant machine.Provisioner = "vagrant"

func init() {
	machine.RegisterProvisioner(vagrant, NewVagrantMachines)
}
