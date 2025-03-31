package template

import (
	"fmt"
	"io"

	"github.com/footprintai/multikf/pkg/k8s"
	"github.com/footprintai/multikf/pkg/machine"
	"github.com/footprintai/multikf/pkg/mirror"
)

type TemplateExecutor interface {
	Filename() string
	Execute(io.Writer) error
	Populate(interface{}) error
}

type NameGetter interface {
	GetName() string
}

type NodeVersionGetter interface {
	GetNodeVersion() k8s.KindK8sVersion
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

type DefaultPasswordGetter interface {
	GetDefaultPassword() string
}

type AuditEnabler interface {
	AuditEnabled() bool
	AuditFileAbsolutePath() string
}

type WorkersGetter interface {
	GetWorkers() []Worker
}

// RegistryGetter defines the interface to get registry mirror configuration
type RegistryGetter interface {
	GetRegistryMirrors() []mirror.Registry
}

type Worker struct {
	Id          string
	UseGPU      bool
	LocalPath   string
	NodeVersion string
}

type NodeLabelsGetter interface {
	GetNodeLabels() []machine.NodeLabel
}

type LocalPathGetter interface {
	LocalPath() string
}

type K8sNodeVersion struct {
	K8sVersion string // started with v1.26.x
	SHA256     string
}

func (k K8sNodeVersion) String() string {
	if k.SHA256 == "" {
		return fmt.Sprintf("kindest/node:%s", k.K8sVersion)
	}
	return fmt.Sprintf("kindest/node:%s@sha256:%s", k.K8sVersion, k.SHA256)
}
