package machine

import (
	"github.com/footprintai/multikf/pkg/k8s"
	"github.com/footprintai/multikf/pkg/machine/cmd/kubectl"
	"github.com/footprintai/multikf/pkg/mirror"
)

type MachineCURDFactory interface {
	EnsureRuntime() error
	NewMachine(string, MachineConfiger) (MachineCURD, error)
	ListMachines() ([]MachineCURD, error)
}

type MachineConfiger interface {
	GetCPUs() int
	GetMemory() int // in M bytes
	GetGPUs() int
	GetKubeAPIIP() string
	GetExportPorts() []ExportPortPair
	GetForceOverwriteConfig() bool
	AuditEnabled() bool
	GetWorkers() int
	GetNodeLabels() []NodeLabel
	GetLocalPath() string
	GetNodeVersion() k8s.KindK8sVersion
	mirror.Getter // Embed the mirror.Getter interface
	// Info displays all configurations
	Info() string
}

type ExportPortPair struct {
	HostPort      int
	ContainerPort int
}

type NodeLabel struct {
	Key   string
	Value string
}

type MachineType string

func (m MachineType) String() string {
	return string(m)
}

const (
	MachineTypeDocker  MachineType = "docker"
	MachineTypeVagrant MachineType = "vagrant"
)

type MachineCURD interface {
	Name() string
	Type() MachineType
	// HostDir returns the configuration files used for that particular machine under host
	GetKubeCli() *kubectl.CLI
	GetKubeConfig() string
	HostDir() string
	Up() error
	Destroy() error
	Info() (*MachineInfo, error)
	ExportKubeConfig(path string, forceOverwrite bool) error
}

type MachineInfo struct {
	CpuInfo         *CpuInfo
	MemInfo         *MemInfo
	GpuInfo         *GpuInfo
	KubeApi         string
	Status          string
	RegistryMirrors []mirror.Registry // Using the Registry type from mirror package
}
