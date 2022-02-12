package machine

type MachineCURD interface {
	Name() string
	// HostDir returns the configuration files used for that particular machine under host
	HostDir() string
	Up(forceCreate bool) error
	Destroy(force bool) error
	Info() (*MachineInfo, error)
	ExportKubeConfig(path string, forceOverwrite bool) error
}

type MachineInfo struct {
	CpuInfo *CpuInfo
	MemInfo *MemInfo
	Status  string
}
