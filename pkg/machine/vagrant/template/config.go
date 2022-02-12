package template

// TemplateFileConfig is a union template file config
type TemplateFileConfig struct {
	Name string
	// NOTE: only support virtualbox now

	CPUs   int // number of cpus allocated
	Memory int // number of bytes memory allocated

	// NOTE: GPUs are not supported now
	// GPUs string

	SSHPort     int
	KubeApiPort int
}

func (t *TemplateFileConfig) GetName() string {
	return t.Name
}

func (t *TemplateFileConfig) GetKubeAPIPort() int {
	return t.KubeApiPort
}

func (t *TemplateFileConfig) GetSSHPort() int {
	return t.SSHPort
}

func (t *TemplateFileConfig) GetMemory() int {
	return t.Memory
}

func (t *TemplateFileConfig) GetCPUs() int {
	return t.CPUs
}
