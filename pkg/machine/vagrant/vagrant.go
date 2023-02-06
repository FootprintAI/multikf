package vagrant

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	vagrantclient "github.com/footprintai/multikf/pkg/client/vagrant"
	machine "github.com/footprintai/multikf/pkg/machine"
	machinecmd "github.com/footprintai/multikf/pkg/machine/cmd"
	"github.com/footprintai/multikf/pkg/machine/fsutil"
	"github.com/footprintai/multikf/pkg/machine/kubectl"
	machinekubectl "github.com/footprintai/multikf/pkg/machine/kubectl"
	"github.com/footprintai/multikf/pkg/machine/vagrant/template"
	"sigs.k8s.io/kind/pkg/log"
)

func NewVagrantMachines(logger log.Logger, vagrantDir string, verbose bool) machine.MachineCURDFactory {
	kubecli, _ := machinekubectl.NewCLI(logger, filepath.Join(vagrantDir, "bin"), verbose)
	return &VagrantMachines{
		logger:     logger,
		vagrantDir: vagrantDir,
		verbose:    verbose,
		kubecli:    kubecli,
	}
}

type VagrantMachines struct {
	logger     log.Logger
	vagrantDir string
	verbose    bool
	kubecli    *machinekubectl.CLI
}

func (vm *VagrantMachines) EnsureRuntime() error {
	_, status, err := machinecmd.NewCmd(vm.logger).Run("vagrant", "--version")
	if err != nil {
		return err
	}
	procStatus := <-status
	if procStatus.Exit != 0 {
		return fmt.Errorf("proc(vagrant): vagrant is not installed? Use `vagrant --version` to verify results")
	}
	return nil
}

func (vm *VagrantMachines) NewMachine(name string, options machine.MachineConfiger) (machine.MachineCURD, error) {
	if err := checkMachineNaming(name); err != nil {
		return nil, err
	}
	return &VagrantMachine{
		logger:            vm.logger,
		mtype:             machine.MachineTypeVagrant,
		name:              name,
		vagrantMachineDir: filepath.Join(vm.vagrantDir, name),
		verbose:           vm.verbose,
		options:           options,
		kubecli:           vm.kubecli,
	}, nil
}

func checkMachineNaming(machinename string) error {
	if strings.Contains(machinename, "-") {
		return fmt.Errorf("vagrant: invalid naming. dash('-') is not allowed ")
	}
	return nil
}

type VagrantMachine struct {
	logger            log.Logger
	mtype             machine.MachineType
	name              string
	vagrantMachineDir string
	verbose           bool
	options           machine.MachineConfiger
	kubecli           *machinekubectl.CLI
}

func (v *VagrantMachine) Type() machine.MachineType {
	return v.mtype
}

func (v *VagrantMachine) GetKubeCli() *kubectl.CLI {
	return v.kubecli
}

func (v *VagrantMachine) HostDir() string {
	return v.vagrantMachineDir
}

func (v *VagrantMachine) GetKubeConfig() string {
	return filepath.Join(v.HostDir(), "kubeconfig.yaml")
}

func (v *VagrantMachine) Up() error {
	// TODO: implement with kubeflow options

	if err := v.ensureVagrantFiles(); err != nil {
		return err
	}
	v.logger.V(0).Infof("vagrantmachine(%s): ready to launch machine\n", v.name)
	cli, err := v.NewVagrantCli()
	if err != nil {
		return err
	}
	if err := cli.TryUp(); err != nil {
		return err
	}
	kubeConfigPath := filepath.Join(v.vagrantMachineDir, "kubeconfig.yaml")
	return v.ExportKubeConfig(kubeConfigPath, true)
}

func (v *VagrantMachine) NewVagrantCli() (*vagrantclient.VagrantCli, error) {
	return vagrantclient.NewVagrantCli(v.name, v.vagrantMachineDir, v.logger, v.verbose)
}

func (v *VagrantMachine) ExportKubeConfig(path string, force bool) error {
	fileExists := false
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		fileExists = true
	}
	if fileExists && !force {
		return fmt.Errorf("kubecfg %s exists, use -f to overwrite it\n", path)
	}
	v.logger.V(0).Infof("vagrantmachine(%s): export kubecfg to path:%s\n", v.name, path)
	cli, err := v.NewVagrantCli()
	if err != nil {
		return err
	}
	return cli.Scp("/home/vagrant/.kube/config", path)
}

func (v *VagrantMachine) Destroy() error {
	v.logger.V(0).Infof("vagrantmachine(%s): ready to destroy\n", v.name)
	cli, err := v.NewVagrantCli()
	if err != nil {
		return err
	}
	return cli.Destroy()
}

func (v *VagrantMachine) Info() (*machine.MachineInfo, error) {
	cli, err := v.NewVagrantCli()
	if err != nil {
		return nil, err
	}
	meminfo, err := machine.NewMemInfoParserHelper(cli.SshExec("cat /proc/meminfo"))
	if err != nil {
		return nil, err
	}
	cpuinfo, err := machine.NewCpuInfoParserHelper(cli.SshExec("cat /proc/cpuinfo"))
	if err != nil {
		return nil, err
	}
	status := cli.Status()
	return &machine.MachineInfo{
		CpuInfo: cpuinfo,
		MemInfo: meminfo,
		GpuInfo: &machine.GpuInfo{},
		Status:  status,
	}, nil
}

func (v *VagrantMachine) Name() string {
	return v.name
}

// TODO(hsiny): add force overwrite options
func (v *VagrantMachine) ensureVagrantFiles() error {
	// only check Vagrantfile
	v.logger.V(0).Infof("vagrantmachine dir:%s\n", v.vagrantMachineDir)
	if !hasVagrantfileInDir(v.vagrantMachineDir) || v.options.GetForceOverwriteConfig() {
		v.logger.V(0).Infof("vagrantmachine(%s): prepare files under %s\n", v.name, v.vagrantMachineDir)
		if err := v.prepareFiles(); err != nil {
			return err
		}
		return nil
	}
	// vagrantfile exists
	v.logger.V(0).Infof("vagrantmachine: Vagrantfile exists, reuse it\n")
	return nil
	//return fmt.Errorf("vagrantmahcine: Vagrantfile exists")
}

func (v *VagrantMachine) prepareFiles() error {
	sshport, err := machine.FindFreeSSHPort()
	if err != nil {
		return err
	}
	kubeport, err := machine.FindFreeKubeApiPort()
	if err != nil {
		return err
	}
	v.logger.V(0).Infof("vagrantmachine(%s): get port (%d,%d) for ssh and kubeapi\n", v.name, sshport, kubeport)
	tmplConfig := template.NewVagrantTemplateConfig(
		v.name,
		v.options.GetCPUs(),
		v.options.GetMemory(),
		sshport,
		kubeport,
		v.options.GetKubeAPIIP(),
		v.options.GetGPUs(),
		v.options.GetExportPorts(),
		v.options.AuditEnabled(),
		"/tmp/audit-policy.yaml", /*for vagrant, we will copy the file under /tmp and run local installation*/
		v.options.GetWorkers(),
		v.options.GetNodeLabels(),
	)

	vfolder := NewVagrantFolder(v.vagrantMachineDir)
	if err := vfolder.GenerateVagrantFiles(tmplConfig); err != nil {
		return err
	}
	return nil
}

func (vm *VagrantMachines) ListMachines() ([]machine.MachineCURD, error) {
	var machines []machine.MachineCURD
	//machineNamesMap := map[string]*OutputVagrantMachine{}
	vfs := os.DirFS(vm.vagrantDir)
	entries, err := fs.ReadDir(vfs, ".")
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !hasVagrantfileInDir(filepath.Join(vm.vagrantDir, entry.Name())) {
			continue
		} else {
			machineName := entry.Name()
			m, _ := vm.NewMachine(machineName, nil)
			machines = append(machines, m)
		}
	}
	return machines, nil
}

func hasVagrantfileInDir(folderPath string) bool {
	if folderPath == "bin" {
		return false
	}
	return fsutil.Exists(os.DirFS(folderPath), "Vagrantfile")
}
