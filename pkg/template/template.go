package template

import "io"

type TemplateExecutor interface {
	Filename() string
	Execute(io.Writer) error
	Populate(interface{}) error
}

type NameGetter interface {
	GetName() string
}

type KubeAPIPortGetter interface {
	GetKubeAPIPort() int
}

type SSHPortGetter interface {
	GetSSHPort() int
}

type CpuMemoryGetter interface {
	GetCPUs() int
	GetMemory() int
}
