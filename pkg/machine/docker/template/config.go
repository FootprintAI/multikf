package template

import (
	"github.com/footprintai/multikf/pkg/k8s"
	"github.com/footprintai/multikf/pkg/machine"
	"github.com/footprintai/multikf/pkg/mirror"
	pkgtemplateconfig "github.com/footprintai/multikf/pkg/template/config"
)

type DockerHostmachineTemplateConfig struct {
	*pkgtemplateconfig.DefaultTemplateConfig
}

func NewDockerHostmachineTemplateConfig(
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
) *DockerHostmachineTemplateConfig {
	return &DockerHostmachineTemplateConfig{
		DefaultTemplateConfig: pkgtemplateconfig.NewDefaultTemplateConfig(
			name,
			cpus,
			memory,
			sshport,
			kubeApiPort,
			kubeApiIP,
			gpus,
			exportPorts,
			auditEnabled,
			auditFileAbsolutePath,
			workerCount,
			nodeLabels,
			localPath,
			nodeVersion,
			registryMirrors,
		),
	}
}

func (d *DockerHostmachineTemplateConfig) GetSSHPort() int {
	// no ssh port for docker hostmachine
	return -1
}

func (d *DockerHostmachineTemplateConfig) GetMemory() int {
	// NOTE: for hostmachine like docker, we wont be able to control cpu/memory as kubelet will break the jail sent by dockerd
	return -1
}

func (d *DockerHostmachineTemplateConfig) GetCPUs() int {
	// NOTE: for hostmachine like docker, we wont be able to control cpu/memory as kubelet will break the jail sent by dockerd
	return -1
}
