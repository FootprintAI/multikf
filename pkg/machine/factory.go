package machine

import (
	"errors"
	"fmt"
)

type Provisioner string

func (p Provisioner) String() string {
	return string(p)
}

func ParseProvisioner(s string) (Provisioner, error) {
	for p, _ := range provisionerRegister {
		if s == p.String() {
			return p, nil
		}
	}
	return Unknwon, errors.New("unknown provisioner")
}

const (
	Unknwon Provisioner = "unknown"
)

func NewMachineFactory(provisioner Provisioner, dir string, verbose bool) (MachinesCURD, error) {
	fac, found := provisionerRegister[provisioner]
	if !found {
		return nil, fmt.Errorf("provisioner:%s is not found\n", provisioner)
	}
	return fac(dir, verbose), nil
}

type FactoryFunc func(dir string, verbose bool) MachinesCURD

var provisionerRegister = map[Provisioner]FactoryFunc{}

func RegisterProvisioner(p Provisioner, fac FactoryFunc) error {
	if _, found := provisionerRegister[p]; found {
		return fmt.Errorf("duplicated register proviosner:%s\n", p)
	}
	provisionerRegister[p] = fac
	return nil
}