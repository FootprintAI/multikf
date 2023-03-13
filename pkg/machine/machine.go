package machine

import "github.com/footprintai/multikf/pkg/machine/kubectl"

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
	CpuInfo *CpuInfo
	MemInfo *MemInfo
	GpuInfo *GpuInfo
	KubeApi string
	Status  string
}
