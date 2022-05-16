package template

import (
	"io"

	"github.com/footprintai/multikf/pkg/machine"
)

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

type KubeAPIIPGetter interface {
	GetKubeAPIIP() string
}

type SSHPortGetter interface {
	GetSSHPort() int
}

type CpuMemoryGetter interface {
	GetCPUs() int
	GetMemory() int
}

type GpuGetter interface {
	GetGPUs() int
}

type ExportPortsGetter interface {
	GetExportPorts() []machine.ExportPortPair
}
