package template

import "github.com/footprintai/multikf/pkg/machine"

type TemplateFileConfig struct {
	// for host machine, we won't able to configure cpu/memory used, as kubelet inside a container can still access its host.
	Name        string
	SSHPort     int
	KubeApiPort int
	KubeApiIP   string
	GPUs        int
	ExportPorts []machine.ExportPortPair

	DefaultPassword string
}

func (t *TemplateFileConfig) GetName() string {
	return t.Name
}

func (t *TemplateFileConfig) GetKubeAPIPort() int {
	return t.KubeApiPort
}

func (t *TemplateFileConfig) GetKubeAPIIP() string {
	return t.KubeApiIP
}

func (t *TemplateFileConfig) GetGPUs() int {
	return t.GPUs
}

func (t *TemplateFileConfig) GetSSHPort() int {
	return t.SSHPort
}

func (t *TemplateFileConfig) GetExportPorts() []machine.ExportPortPair {
	return t.ExportPorts
}

func (t *TemplateFileConfig) GetDefaultPassword() string {
	return t.DefaultPassword
}