package template

import (
	"github.com/footprintai/multikf/pkg/k8s"
	"github.com/footprintai/multikf/pkg/machine"
	pkgtemplateconfig "github.com/footprintai/multikf/pkg/template/config"
)

type VagrantTemplateConfig struct {
	*pkgtemplateconfig.DefaultTemplateConfig
}

func NewVagrantTemplateConfig(name string, cpus int, memory int, sshport int, kubeApiPort int, kubeApiIP string, gpus int, exportPorts []machine.ExportPortPair, auditEnabled bool, auditFileAbsolutePath string, workerCount int, nodeLabels []machine.NodeLabel, localPath string, nodeVersion k8s.KindK8sVersion) *VagrantTemplateConfig {
	return &VagrantTemplateConfig{
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
		),
	}
}

func (v *VagrantTemplateConfig) GPUs() int {
	return 0
}
