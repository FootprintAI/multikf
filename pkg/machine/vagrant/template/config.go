package template

import (
	"github.com/footprintai/multikf/pkg/machine"
	pkgtemplateconfig "github.com/footprintai/multikf/pkg/template/config"
)

type VagrantTemplateConfig struct {
	*pkgtemplateconfig.DefaultTemplateConfig
}

func NewVagrantTemplateConfig(name string, cpus int, memory int, sshport int, kubeApiPort int, kubeApiIP string, gpus int, exportPorts []machine.ExportPortPair, auditEnabled bool, auditFileAbsolutePath string) *VagrantTemplateConfig {
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
		),
	}
}

func (v *VagrantTemplateConfig) GPUs() int {
	return 0
}
