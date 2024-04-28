package multikf

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/footprintai/multikf/pkg/machine"
	"github.com/footprintai/multikf/pkg/machine/plugins"
	"sigs.k8s.io/kind/pkg/log"
)

func findMachineByName(name string, logger log.Logger) (machine.MachineCURD, error) {
	//for _, provisioner := []machine.Provisioner {}
	var found machine.MachineCURD
	var outErr error = errors.New("machine: not found")

	machine.ForEachProvisioner(func(p machine.Provisioner) {
		vag, err := machine.NewMachineFactory(
			p,
			logger,
			viperConfigKeyRootDir.GetString(),
			viperConfigKeyVerbose.GetBool(),
		)
		if err != nil {
			logger.Errorf("machine.find: failed, err:%+v\n", err)
			return
		}
		machines, err := vag.ListMachines()
		if err != nil {
			logger.Errorf("machine.find: failed, err:%+v\n", err)
			return
		}
		for _, machine := range machines {
			if machine.Name() == name {
				found = machine
				outErr = nil
				return
			}
		}
	})
	return found, outErr
}

func newMachineFactoryWithProvisioner(p machine.Provisioner, logger log.Logger) (machine.MachineCURDFactory, error) {
	vag, err := machine.NewMachineFactory(
		p,
		logger,
		viperConfigKeyRootDir.GetString(),
		viperConfigKeyVerbose.GetBool(),
	)
	if err != nil {
		return nil, err
	}
	if err := vag.EnsureRuntime(); err != nil {
		return nil, err
	}
	return vag, nil
}

var (
	_ machine.MachineConfiger = &machineConfig{}
)

type machineConfig struct {
	logger          log.Logger
	Cpus            int            `json:"cpus"`
	MemoryInG       int            `json:"memoryInG"`
	UseGPUs         int            `json:useGpus`
	KubeAPIIP       string         `json:"kubeapi_ip"`
	ExportPorts     string         `json:"export_ports"`
	DefaultPassword string         `json:"default_password"`
	ForceOverwrite  bool           `json:"force_overwrite"`
	IsAuditEnabled  bool           `json:"audit_enabled"`
	Workers         int            `json:"workers"`
	NodeLabels      string         `json:"node_labels"`
	LocalPath       string         `json:"local_path"`
	NodeVersion     K8sNodeVersion `json:"node_version"`
}

type K8sNodeVersion struct {
	K8sVersion string `json:"k8s_version"` // started with v1.26.x
	SHA256     string `json:"sha256"`
}

func (k K8sNodeVersion) String() string {
	return fmt.Sprintf("kindest/node:%s@sha256:%s", k.K8sVersion, k.SHA256)
}

func (m machineConfig) Info() string {
	bb, _ := json.Marshal(m)
	return string(bb)

}

func (m machineConfig) GetNodeVersion() string {
	return m.NodeVersion.String()
}

func (m machineConfig) GetCPUs() int {
	return m.Cpus
}

// GetMemory returns memory amount in M bytes
func (m machineConfig) GetMemory() int {
	return m.MemoryInG * 1024
}

func (m machineConfig) GetGPUs() int {
	return m.UseGPUs
}

func (m machineConfig) GetKubeAPIIP() string {
	return m.KubeAPIIP
}

func (m machineConfig) AuditEnabled() bool {
	return m.IsAuditEnabled
}

func (m machineConfig) GetExportPorts() []machine.ExportPortPair {
	if len(m.ExportPorts) == 0 {
		m.logger.V(1).Infof("getexportport: export nothing\n")
		return nil
	}
	tokens := strings.Split(m.ExportPorts, ",")
	var exportPorts []machine.ExportPortPair
	for _, token := range tokens {
		subtokens := strings.Split(token, ":")
		if len(subtokens) != 2 {
			m.logger.Errorf("getexportport: parse failed, expect: a:b but got:%s\n", token)
			continue
		}
		hostport, err := strconv.Atoi(subtokens[0])
		if err != nil {
			m.logger.Errorf("getexportport: parse failed, err:%+v\n", err)
			continue
		}
		containerport, err := strconv.Atoi(subtokens[1])
		if err != nil {
			m.logger.Errorf("getexportport: parse failed, err:%+v\n", err)
			continue
		}
		exportPorts = append(exportPorts, machine.ExportPortPair{HostPort: hostport, ContainerPort: containerport})
	}
	m.logger.V(1).Infof("getexportport: export ports:%+v\n", exportPorts)
	return exportPorts
}

func (m machineConfig) GetForceOverwriteConfig() bool {
	return m.ForceOverwrite
}

func (m machineConfig) GetWorkers() int {
	return m.Workers
}

// a=b,c=d
func (m machineConfig) GetNodeLabels() []machine.NodeLabel {
	if len(m.NodeLabels) == 0 {
		m.logger.V(1).Infof("getnodelabel: no label\n")
		return nil
	}
	tokens := strings.Split(m.NodeLabels, ",")
	var nodeLabels []machine.NodeLabel
	for _, token := range tokens {
		subtokens := strings.Split(token, "=")
		if len(subtokens) != 2 {
			m.logger.Errorf("getnodelabel: parse failed, expect: key=value but got:%s\n", token)
			continue
		}
		nodeLabels = append(nodeLabels, machine.NodeLabel{Key: subtokens[0], Value: subtokens[1]})
	}
	m.logger.V(1).Infof("getnodelabel: labels:%+v\n", nodeLabels)
	return nodeLabels
}

func (m machineConfig) GetLocalPath() string {
	return m.LocalPath
}

type kubeflowPlugin struct {
	withKubeflowDefaultPassword string
	kubeflowVersion             plugins.TypePluginVersion
}

func (k kubeflowPlugin) PluginType() plugins.TypePlugin {
	return plugins.TypePluginKubeflow
}

func (k kubeflowPlugin) PluginVersion() plugins.TypePluginVersion {
	return k.kubeflowVersion

}

func (k kubeflowPlugin) GetDefaultPassword() string {
	return k.withKubeflowDefaultPassword
}
