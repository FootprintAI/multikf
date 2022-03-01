package machine

type MachinesCURD interface {
	NewMachine(string, MachineConfiger) (MachineCURD, error)
	ListMachines() ([]MachineCURD, error)
}

type MachineConfiger interface {
	GetCPUs() int
	GetMemory() int // in M bytes
}

type MachineCURD interface {
	Name() string
	// HostDir returns the configuration files used for that particular machine under host
	HostDir() string
	Up(forceCreate bool) error
	Destroy(force bool) error
	Info() (*MachineInfo, error)
	ExportKubeConfig(path string, forceOverwrite bool) error
	Portforward(svc, namespace string, fromPort int) (int, error)
}

type MachineInfo struct {
	CpuInfo *CpuInfo
	MemInfo *MemInfo
	GpuInfo *GpuInfo
	KubeApi string
	Status  string
}
