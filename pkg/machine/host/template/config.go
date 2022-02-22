package template

type TemplateFileConfig struct {
	// for host machine, we won't able to configure cpu/memory used, as kubelet inside a container can still access its host.
	Name        string
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
