package machine

import (
	"errors"
	"fmt"

	"sigs.k8s.io/kind/pkg/log"
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

func MustParseProvisioner(s string) Provisioner {
	//fmt.Printf("provisionerstr:%s\n", s)
	p, err := ParseProvisioner(s)
	if err != nil {
		panic(err)
	}
	return p
}

const (
	Unknwon Provisioner = "unknown"
)

func NewMachineFactory(provisioner Provisioner, logger log.Logger, dir string, verbose bool) (MachineCURDFactory, error) {
	logger.V(1).Infof("allocate machine with provisioner:%+v\n", provisioner)
	fac, found := provisionerRegister[provisioner]
	if !found {
		return nil, fmt.Errorf("provisioner:%s is not found\n", provisioner)
	}
	return fac(logger, dir, verbose), nil
}

type FactoryFunc func(logger log.Logger, dir string, verbose bool) MachineCURDFactory

var provisionerRegister = map[Provisioner]FactoryFunc{}

func RegisterProvisioner(p Provisioner, fac FactoryFunc) error {
	//fmt.Printf("register provisoner:%s\n", p.String())
	if _, found := provisionerRegister[p]; found {
		return fmt.Errorf("duplicated register proviosner:%s\n", p)
	}
	provisionerRegister[p] = fac
	return nil
}

func ForEachProvisioner(iter func(p Provisioner)) {
	for p := range provisionerRegister {
		iter(p)
	}
}
