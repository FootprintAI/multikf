package host

import (
	"fmt"
	"path/filepath"

	machine "github.com/footprintai/multikf/pkg/machine"
	"github.com/footprintai/multikf/pkg/machine/host/template"
	"sigs.k8s.io/kind/pkg/log"
)

func NewHostMachines(logger log.Logger, hostDir string, verbose bool) machine.MachinesCURD {
	cli, _ := NewCLI(logger, filepath.Join(hostDir, "bin"), verbose)
	return &HostMachines{
		logger:  logger,
		hostDir: hostDir,
		verbose: verbose,
		cli:     cli,
	}
}

type HostMachines struct {
	logger  log.Logger
	hostDir string
	verbose bool
	cli     *CLI
}

func (hm *HostMachines) NewMachine(name string, options machine.MachineConfiger) (machine.MachineCURD, error) {
	return &HostMachine{
		logger:         hm.logger,
		name:           name,
		containername:  NewContainerName(name),
		hostMachineDir: filepath.Join(hm.hostDir, name),
		verbose:        hm.verbose,
		kubeconfig:     filepath.Join(hm.hostDir, name, "kubeconfig.yaml"),
		cli:            hm.cli,
		options:        options,
	}, nil
}

func (hm *HostMachines) ListMachines() ([]machine.MachineCURD, error) {
	clusternames, err := hm.cli.ListClusters()
	if err != nil {
		return nil, err
	}
	var machines []machine.MachineCURD
	for _, clustername := range clusternames {
		m, _ := hm.NewMachine(clustername, nil)
		machines = append(machines, m)
	}
	return machines, nil
}

type HostMachine struct {
	logger         log.Logger
	name           string
	containername  ContainerName
	hostMachineDir string
	verbose        bool
	kubeconfig     string // filepath to kubeconfig
	options        machine.MachineConfiger

	cli *CLI
}

func (h *HostMachine) ensureFiles() error {
	f := filepath.Join(h.hostMachineDir, "kind-config.yaml")
	if !fileExists(f) {
		h.logger.V(1).Infof("hostmachine(%s): prepare files under %s\n", h.name, h.hostMachineDir)
		if err := h.prepareFiles(); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (h *HostMachine) prepareFiles() error {
	kubeport, err := machine.FindFreeKubeApiPort()
	if err != nil {
		return err
	}
	h.logger.V(1).Infof("hostmachine(%s): get port (%d) for kubeapi\n", h.name, kubeport)
	tmplConfig := &template.TemplateFileConfig{
		Name:        h.name,
		KubeApiPort: kubeport,
		KubeApiIP:   h.options.GetKubeAPIIP(),
		GPUs:        h.options.GetGPUs(),
	}

	vfolder := NewHostFolder(h.hostMachineDir)
	if err := vfolder.GenerateFiles(tmplConfig); err != nil {
		return err
	}
	return nil
}

func (h *HostMachine) Name() string {
	return h.name
}

func (h *HostMachine) HostDir() string {
	return h.hostMachineDir
}

func (h *HostMachine) Up(forceDeletedIfNecessary bool, withKubeflow bool) error {
	if err := h.ensureFiles(); err != nil {
		return err
	}
	kindConfigPath := filepath.Join(h.hostMachineDir, "kind-config.yaml")
	kubeConfigPath := filepath.Join(h.hostMachineDir, "kubeconfig.yaml")
	h.logger.V(1).Infof("hostmachine(%s): check %s for kubeconfig.yaml\n", h.name, kubeConfigPath)

	if err := h.cli.ProvisonCluster(kindConfigPath); err != nil {
		return err
	}
	// install required pkgs
	if err := h.cli.InstallRequiredPkgs(h.containername); err != nil {
		return err
	}
	if err := h.cli.GetKubeConfig(h.name, kubeConfigPath); err != nil {
		return err
	}
	if withKubeflow {
		kfManifestPath := filepath.Join(h.hostMachineDir, "kubeflow-manifest-v1.4.1.yaml")
		if err := h.cli.InstallKubeflow(kubeConfigPath, kfManifestPath); err != nil {
			return err
		}
		if err := h.cli.PatchKubeflow(kubeConfigPath); err != nil {
			return err
		}
	}
	return nil
}

func (h *HostMachine) ensureKubeconfig() error {
	if !fileExists(h.kubeconfig) {
		return h.cli.GetKubeConfig(h.name, h.kubeconfig)
	}
	return nil
}

func (h *HostMachine) ExportKubeConfig(path string, force bool) error {
	if fileExists(path) && !force {
		h.logger.Errorf("host: local kubeconfig file %s exists, use -f to force overwrite", path)
		return fmt.Errorf("local kubeconfig file %s exists, use -f to force overwrite", path)
	}
	return h.cli.GetKubeConfig(h.name, path)
}

func (h *HostMachine) Destroy(force bool) error {
	return h.cli.RemoveCluster(h.name)
}

func (h *HostMachine) Info() (*machine.MachineInfo, error) {
	meminfo, err := machine.NewMemInfoParserHelper(h.cli.RemoteExec(h.containername, "cat /proc/meminfo"))
	if err != nil {
		return nil, err
	}
	cpuinfo, err := machine.NewCpuInfoParserHelper(h.cli.RemoteExec(h.containername, "cat /proc/cpuinfo"))
	if err != nil {
		return nil, err
	}
	gpuinfo, err := machine.NewGpuInfoParserHelper(h.cli.RemoteExec(h.containername, "/usr/bin/nvidia-smi -x -q -a"))
	if err != nil {
		h.logger.Errorf("host: get cpu info failed, err:%s\n", err)
		gpuinfo = &machine.GpuInfo{}
	}
	status, err := h.cli.GetClusterStatus(h.containername)
	if err != nil {
		return nil, err
	}
	return &machine.MachineInfo{
		CpuInfo: cpuinfo,
		MemInfo: meminfo,
		GpuInfo: gpuinfo,
		Status:  status,
	}, nil
}

func (h *HostMachine) Portforward(svc, namespace string, fromPort int) (int, error) {
	if err := h.ensureKubeconfig(); err != nil {
		return 0, err
	}
	destPort, err := machine.FindFreePort()
	if err != nil {
		return 0, err
	}
	h.logger.V(0).Infof("now you can open http://localhost:%d\n", destPort)
	return destPort, h.cli.Portforward(h.kubeconfig, svc, namespace, fromPort, destPort)
}

func (h *HostMachine) GetPods(namespace string) error {
	if err := h.ensureKubeconfig(); err != nil {
		return err
	}
	return h.cli.GetPods(h.kubeconfig, namespace)
}
