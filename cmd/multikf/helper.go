package multikf

import (
	"errors"
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

type machineConfig struct {
	logger          log.Logger
	cpus            int
	memoryInG       int
	useGPUs         int
	kubeAPIIP       string
	exportPorts     string
	defaultPassword string
	forceOverwrite  bool
	auditEnabled    bool
}

func (m machineConfig) GetCPUs() int {
	return m.cpus
}

// GetMemory returns memory amount in M bytes
func (m machineConfig) GetMemory() int {
	return m.memoryInG * 1024
}

func (m machineConfig) GetGPUs() int {
	return m.useGPUs
}

func (m machineConfig) GetKubeAPIIP() string {
	return m.kubeAPIIP
}

func (m machineConfig) AuditEnabled() bool {
	return m.auditEnabled
}

func (m machineConfig) GetExportPorts() []machine.ExportPortPair {
	if len(m.exportPorts) == 0 {
		m.logger.V(1).Infof("getexportport: export nothing\n")
		return nil
	}
	tokens := strings.Split(m.exportPorts, ",")
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
	return m.forceOverwrite
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
