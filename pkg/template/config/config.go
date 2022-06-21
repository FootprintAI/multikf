package config

import "github.com/footprintai/multikf/pkg/machine"

type DefaultTemplateConfig struct {
	name        string
	cpus        int // number of cpus allocated
	memory      int // number of bytes memory allocated
	sshPort     int
	kubeApiPort int
	kubeApiIP   string
	gpus        int
	exportPorts []machine.ExportPortPair
}

func NewDefaultTemplateConfig(name string, cpus int, memory int, sshport int, kubeApiPort int, kubeApiIP string, gpus int, exportPorts []machine.ExportPortPair) *DefaultTemplateConfig {
	return &DefaultTemplateConfig{
		name:        name,
		cpus:        cpus,
		memory:      memory,
		sshPort:     sshport,
		kubeApiPort: kubeApiPort,
		kubeApiIP:   kubeApiIP,
		gpus:        gpus,
		exportPorts: exportPorts,
	}
}

func (t *DefaultTemplateConfig) GetName() string {
	return t.name
}

func (t *DefaultTemplateConfig) GetMemory() int {
	return t.memory
}

func (t *DefaultTemplateConfig) GetCPUs() int {
	return t.cpus
}

func (t *DefaultTemplateConfig) GetKubeAPIPort() int {
	return t.kubeApiPort
}

func (t *DefaultTemplateConfig) GetKubeAPIIP() string {
	return t.kubeApiIP
}

func (t *DefaultTemplateConfig) GetGPUs() int {
	return t.gpus
}

func (t *DefaultTemplateConfig) GetSSHPort() int {
	return t.sshPort
}

func (t *DefaultTemplateConfig) GetExportPorts() []machine.ExportPortPair {
	return t.exportPorts
}
