package host

import (
	machine "github.com/footprintai/multikind/pkg/machine"
)

func NewHostMachines() *HostMachines {
	return &HostMachines{}
}

type HostMachines struct{}

func (hm *HostMachines) NewMachine(name string) *HostMachine {
	return &HostMachine{
		name: name,
	}
}

func (hm *HostMachines) ListMachines() ([]machine.MachineCURD, error) {

	var machines []machine.MachineCURD
	return machines, nil
}

type HostMachine struct {
	name string
}

func (h *HostMachine) Name() string {
	return h.name
}

func (h *HostMachine) HostDir() string {
	return ""
}

func (h *HostMachine) Up(forceDeletedIfNecessary bool) error {
	return nil
}

func (h *HostMachine) ExportKubeConfig(path string, force bool) error {
	return nil
}

func (h *HostMachine) Destroy(force bool) error {
	return nil
}

func (h *HostMachine) Info() (*machine.MachineInfo, error) {
	return nil, nil
}
