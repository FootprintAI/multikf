package config

import (
	"fmt"

	"github.com/footprintai/multikf/pkg/k8s"
	"github.com/footprintai/multikf/pkg/machine"
	"github.com/footprintai/multikf/pkg/mirror"
	"github.com/footprintai/multikf/pkg/template"
)

var (
	_ template.KindConfiger = &DefaultTemplateConfig{}
)

type DefaultTemplateConfig struct {
	name                  string
	cpus                  int // number of cpus allocated
	memory                int // number of bytes memory allocated
	sshPort               int
	kubeApiPort           int
	kubeApiIP             string
	gpus                  int
	exportPorts           []machine.ExportPortPair
	auditEnabled          bool
	auditFileAbsolutePath string
	workerCount           int
	nodeLabels            []machine.NodeLabel
	localPath             string
	nodeVersion           k8s.KindK8sVersion
	registryMirrors       []mirror.Registry // Using the Registry type from mirror package
}

// NewDefaultTemplateConfig creates a default template config
// registryMirrors can be nil or an empty slice if no registry mirrors are needed
func NewDefaultTemplateConfig(
	name string,
	cpus int,
	memory int,
	sshport int,
	kubeApiPort int,
	kubeApiIP string,
	gpus int,
	exportPorts []machine.ExportPortPair,
	auditEnabled bool,
	auditFileAbsolutePath string,
	workerCount int,
	nodeLabels []machine.NodeLabel,
	localPath string,
	nodeVersion k8s.KindK8sVersion,
	registryMirrors []mirror.Registry,
) *DefaultTemplateConfig {
	// If registryMirrors is nil, initialize as empty slice
	if registryMirrors == nil {
		registryMirrors = []mirror.Registry{}
	}

	return &DefaultTemplateConfig{
		name:                  name,
		cpus:                  cpus,
		memory:                memory,
		sshPort:               sshport,
		kubeApiPort:           kubeApiPort,
		kubeApiIP:             kubeApiIP,
		gpus:                  gpus,
		exportPorts:           exportPorts,
		auditEnabled:          auditEnabled,
		auditFileAbsolutePath: auditFileAbsolutePath,
		workerCount:           workerCount,
		nodeLabels:            nodeLabels,
		localPath:             localPath,
		nodeVersion:           nodeVersion,
		registryMirrors:       registryMirrors,
	}
}

// AddRegistryMirror adds a registry mirror to the config
func (t *DefaultTemplateConfig) AddRegistryMirror(registry mirror.Registry) {
	t.registryMirrors = append(t.registryMirrors, registry)
}

// AddAuthenticatedRegistryMirror adds a registry mirror with authentication to the config
func (t *DefaultTemplateConfig) AddAuthenticatedRegistryMirror(source string, mirrorURL string, username string, password string) {
	t.registryMirrors = append(t.registryMirrors, mirror.Registry{
		Source:  source,
		Mirrors: []string{mirrorURL},
		Auth: &mirror.Auth{
			Username: username,
			Password: password,
		},
	})
}

// GetRegistry implements the mirror.Getter interface
func (t *DefaultTemplateConfig) GetRegistry() []mirror.Registry {
	return t.registryMirrors
}

func (t *DefaultTemplateConfig) GetName() string {
	return t.name
}

func (t *DefaultTemplateConfig) GetNodeVersion() k8s.KindK8sVersion {
	return t.nodeVersion
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

func (t *DefaultTemplateConfig) AuditEnabled() bool {
	return t.auditEnabled
}

func (t *DefaultTemplateConfig) AuditFileAbsolutePath() string {
	return t.auditFileAbsolutePath
}

func (t *DefaultTemplateConfig) GetWorkers() []template.Worker {
	ids := make([]template.Worker, t.workerCount, t.workerCount)
	for i := 0; i < t.workerCount; i++ {
		ids[i] = template.Worker{
			Id:          fmt.Sprintf("%d", i),
			UseGPU:      t.GetGPUs() > 0,
			LocalPath:   t.LocalPath(),
			NodeVersion: t.GetNodeVersion().String(),
		}
	}
	return ids
}

func (t *DefaultTemplateConfig) GetNodeLabels() []machine.NodeLabel {
	return t.nodeLabels
}

func (t *DefaultTemplateConfig) LocalPath() string {
	return t.localPath
}
