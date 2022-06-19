package machine

type MachinesCURD interface {
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
}

type ExportPortPair struct {
	HostPort      int
	ContainerPort int
}

type MachineCURD interface {
	Name() string
	// HostDir returns the configuration files used for that particular machine under host
	HostDir() string
	Up(forceCreate bool, withKubeflow bool) error
	Destroy(force bool) error
	Info() (*MachineInfo, error)
	ExportKubeConfig(path string, forceOverwrite bool) error
	Portforward(svc, namespace string, fromPort int) (int, error)
	GetPods(namespace string) error
}

type MachineInfo struct {
	CpuInfo *CpuInfo
	MemInfo *MemInfo
	GpuInfo *GpuInfo
	KubeApi string
	Status  string
}
