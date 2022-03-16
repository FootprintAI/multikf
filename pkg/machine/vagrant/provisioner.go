package vagrant

import (
	machine "github.com/footprintai/multikf/pkg/machine"
)

var vagrant machine.Provisioner = "vagrant"

func init() {
	machine.RegisterProvisioner(vagrant, NewVagrantMachines)
}
