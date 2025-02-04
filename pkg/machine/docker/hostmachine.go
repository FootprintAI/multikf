package docker

import (
	"fmt"
	"path/filepath"

	"github.com/footprintai/multikf/pkg/k8s"
	machine "github.com/footprintai/multikf/pkg/machine"
	machinecmd "github.com/footprintai/multikf/pkg/machine/cmd"
	machinekindcmd "github.com/footprintai/multikf/pkg/machine/cmd/kind"
	machinekubectlcmd "github.com/footprintai/multikf/pkg/machine/cmd/kubectl"
	"github.com/footprintai/multikf/pkg/machine/docker/template"
	"github.com/footprintai/multikf/pkg/machine/fsutil"
	"sigs.k8s.io/kind/pkg/log"
)

func NewHostMachines(logger log.Logger, hostDir string, verbose bool) machine.MachineCURDFactory {
	kindcli, _ := machinekindcmd.NewCLI(logger, filepath.Join(hostDir, "bin"), verbose)
	dockercli, _ := NewDockerCli(logger, verbose)
	return &HostMachines{
		logger:    logger,
		hostDir:   hostDir,
		verbose:   verbose,
		kindcli:   kindcli,
		dockercli: dockercli,
	}
}

func (hm *HostMachines) EnsureRuntime() error {
	_, status, err := machinecmd.NewCmd(hm.logger).Run("docker", "version")
	if err != nil {
		return err
	}
	procStatus := <-status
	if procStatus.Exit != 0 {
		return fmt.Errorf("proc(docker): docker daemon is not running? Use `docker ps` to verify results")
	}
	return nil
}

type HostMachines struct {
	logger    log.Logger
	hostDir   string
	verbose   bool
	kindcli   *machinekindcmd.CLI
	dockercli *DockerCli
}

func (hm *HostMachines) NewMachine(name string, options machine.MachineConfiger) (machine.MachineCURD, error) {
	var nodeVersion k8s.KindK8sVersion
	if options != nil {
		nodeVersion = options.GetNodeVersion()
	}
	kubectlcli, err := machinekubectlcmd.NewCLI(hm.logger, filepath.Join(hm.hostDir, name), hm.verbose, nodeVersion)
	if err != nil {
		return nil, err
	}
	return &HostMachine{
		logger:         hm.logger,
		mtype:          machine.MachineTypeDocker,
		name:           name,
		containername:  NewContainerName(name),
		hostMachineDir: filepath.Join(hm.hostDir, name),
		verbose:        hm.verbose,
		kubeconfig:     filepath.Join(hm.hostDir, name, "kubeconfig.yaml"),
		kubectlcli:     kubectlcli,
		kindcli:        hm.kindcli,
		dockercli:      hm.dockercli,
		options:        options,
	}, nil
}

func (hm *HostMachines) ListMachines() ([]machine.MachineCURD, error) {
	clusternames, err := hm.kindcli.ListClusters()
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
	mtype          machine.MachineType
	containername  ContainerName
	hostMachineDir string
	verbose        bool
	kubeconfig     string // filepath to kubeconfig
	options        machine.MachineConfiger

	kubectlcli *machinekubectlcmd.CLI
	kindcli    *machinekindcmd.CLI
	dockercli  *DockerCli
}

var (
	_ machine.MachineCURD = &HostMachine{}
)

func (h *HostMachine) ensureFiles() error {

	f := filepath.Join(h.hostMachineDir, "kind-config.yaml")
	if !fsutil.FileExists(f) {
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
	tmplConfig := template.NewDockerHostmachineTemplateConfig(
		h.name,
		h.options.GetCPUs(),
		h.options.GetMemory(),
		-1,
		kubeport,
		h.options.GetKubeAPIIP(),
		h.options.GetGPUs(),
		h.options.GetExportPorts(),
		h.options.AuditEnabled(),
		filepath.Join(h.hostMachineDir, "audit-policy.yaml"),
		h.options.GetWorkers(),
		h.options.GetNodeLabels(),
		h.options.GetLocalPath(),
		h.options.GetNodeVersion(),
	)

	vfolder := NewHostFolder(h.hostMachineDir)
	if err := vfolder.GenerateFiles(tmplConfig); err != nil {
		h.logger.Errorf("hostmachine(%s): failed to generate files, err:%+v\n", h.name, err)
		return err
	}
	h.logger.V(1).Infof("hostmachine(%s): configs are prepared\n", h.name)
	return nil
}

func (h *HostMachine) Name() string {
	return h.name
}

func (h *HostMachine) Type() machine.MachineType {
	return h.mtype
}

func (h *HostMachine) GetKubeCli() *machinekubectlcmd.CLI {
	return h.kubectlcli
}

func (h *HostMachine) HostDir() string {
	return h.hostMachineDir
}

func (h *HostMachine) GetKubeConfig() string {
	return filepath.Join(h.HostDir(), "kubeconfig.yaml")
}

func (h *HostMachine) Up() error {
	h.logger.V(1).Infof("hostmachine(%s): configs:%+v\n", h.name, h.options.Info())

	if err := h.ensureFiles(); err != nil {
		return err
	}
	kindConfigPath := filepath.Join(h.hostMachineDir, "kind-config.yaml")
	kubeConfigPath := filepath.Join(h.hostMachineDir, "kubeconfig.yaml")
	h.logger.V(1).Infof("hostmachine(%s): check %s for kubeconfig.yaml\n", h.name, kubeConfigPath)

	if err := h.kindcli.ProvisonCluster(kindConfigPath); err != nil {
		return err
	}
	if err := h.kindcli.GetKubeConfig(h.name, kubeConfigPath); err != nil {
		return err
	}
	return nil
}

func (h *HostMachine) ensureKubeconfig() error {
	if !fsutil.FileExists(h.kubeconfig) {
		return h.kindcli.GetKubeConfig(h.name, h.kubeconfig)
	}
	return nil
}

func (h *HostMachine) ExportKubeConfig(path string, force bool) error {
	if fsutil.FileExists(path) && !force {
		h.logger.Errorf("host: local kubeconfig file %s exists, use -f to force overwrite", path)
		return fmt.Errorf("local kubeconfig file %s exists, use -f to force overwrite", path)
	}
	return h.kindcli.GetKubeConfig(h.name, path)
}

func (h *HostMachine) Destroy() error {
	return h.kindcli.RemoveCluster(h.name)
}

func (h *HostMachine) Info() (*machine.MachineInfo, error) {
	meminfo, err := machine.NewMemInfoParserHelper(h.dockercli.RemoteExec(h.containername, "cat /proc/meminfo"))
	if err != nil {
		return nil, err
	}
	cpuinfo, err := machine.NewCpuInfoParserHelper(h.dockercli.RemoteExec(h.containername, "cat /proc/cpuinfo"))
	if err != nil {
		return nil, err
	}
	gpuinfo, err := machine.NewGpuInfoParserHelper(h.dockercli.RemoteExec(h.containername, "/usr/bin/nvidia-smi -x -q -a"))
	if err != nil {
		h.logger.V(2).Infof("host: get cpu info failed, err:%s\n", err)
		gpuinfo = &machine.GpuInfo{}
	}
	status, err := h.dockercli.GetClusterStatus(h.containername)
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
