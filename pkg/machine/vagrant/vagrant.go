package vagrant

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	vagrantclient "github.com/footprintai/multikf/pkg/client/vagrant"
	machine "github.com/footprintai/multikf/pkg/machine"
	"github.com/footprintai/multikf/pkg/machine/vagrant/template"
	"sigs.k8s.io/kind/pkg/log"
)

func NewVagrantMachines(logger log.Logger, vagrantDir string, verbose bool) machine.MachinesCURD {
	return &VagrantMachines{
		logger:     logger,
		vagrantDir: vagrantDir,
		verbose:    verbose,
	}
}

type VagrantMachines struct {
	logger     log.Logger
	vagrantDir string
	verbose    bool
}

func (vm *VagrantMachines) NewMachine(name string, options machine.MachineConfiger) (machine.MachineCURD, error) {
	return &VagrantMachine{
		logger:            vm.logger,
		name:              name,
		vagrantMachineDir: filepath.Join(vm.vagrantDir, name),
		verbose:           vm.verbose,
		config: &VagrantMachineConfig{
			CPUs:   options.GetCPUs(),
			Memory: options.GetMemory(),
		},
	}, nil
}

type VagrantMachine struct {
	logger            log.Logger
	name              string
	vagrantMachineDir string
	verbose           bool
	config            *VagrantMachineConfig
}

type VagrantMachineConfig struct {
	CPUs   int
	Memory int // measured in M bytes
}

func (v *VagrantMachine) HostDir() string {
	return v.vagrantMachineDir
}

func (v *VagrantMachine) Up(forceDeleteIfNecessary bool, withKubeflow bool) error {
	// TODO: implement with kubeflow options
	if v.config == nil {
		return fmt.Errorf("vagrantmachine requires config when Up")
	}
	if err := v.ensureVagrantFiles(); err != nil {
		return err
	}
	v.logger.V(0).Infof("vagrantmachine(%s): ready to launch machine\n", v.name)
	cli, err := v.NewVagrantCli()
	if err != nil {
		return err
	}
	return cli.TryUp(forceDeleteIfNecessary)
}

func (v *VagrantMachine) NewVagrantCli() (*vagrantclient.VagrantCli, error) {
	return vagrantclient.NewVagrantCli(v.name, v.vagrantMachineDir, v.verbose)
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

func (v *VagrantMachine) Destroy(force bool) error {
	v.logger.V(0).Infof("vagrantmachine(%s): ready to destroy\n", v.name)
	cli, err := v.NewVagrantCli()
	if err != nil {
		return err
	}
	return cli.Destroy(force)
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

func (v *VagrantMachine) Portforward(svc, namespace string, fromPort int) (int, error) {
	return 0, errors.New("todo")
}

func (v *VagrantMachine) GetPods(namespace string) error {
	return errors.New("todo")
}

func (v *VagrantMachine) Name() string {
	return v.name
}

func (v *VagrantMachine) ensureVagrantFiles() error {
	// only check Vagrantfile
	f := filepath.Join(v.vagrantMachineDir, "Vagrantfile")
	if _, err := os.Stat(f); os.IsNotExist(err) {
		v.logger.V(0).Infof("vagrantmachine(%s): prepare files under %s\n", v.name, v.vagrantMachineDir)
		if err := v.prepareFiles(); err != nil {
			return err
		}
		return nil
	}
	// vagrantfile exists
	return fmt.Errorf("vagrantmahcine: Vagrantfile exists")
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
	tmplConfig := &template.TemplateFileConfig{
		Name:        v.name,
		CPUs:        v.config.CPUs,
		Memory:      v.config.Memory,
		SSHPort:     sshport,
		KubeApiPort: kubeport,
	}

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
		if entry.IsDir() {
			machineName := entry.Name()

			m, _ := vm.NewMachine(machineName, nil)
			machines = append(machines, m)
		}
	}
	return machines, nil
}
