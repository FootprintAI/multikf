package vagrant

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	vagrantclient "github.com/footprintai/multikind/pkg/client/vagrant"
	machine "github.com/footprintai/multikind/pkg/machine"
	"github.com/footprintai/multikind/pkg/machine/vagrant/template"
	log "github.com/golang/glog"
)

func NewVagrantMachines(vagrantDir string, verbose bool) *VagrantMachines {
	return &VagrantMachines{
		vagrantDir: vagrantDir,
		verbose:    verbose,
	}
}

type VagrantMachines struct {
	vagrantDir string
	verbose    bool
}

func (vm *VagrantMachines) NewMachine(name string, config *VagrantMachineConfig) *VagrantMachine {
	return &VagrantMachine{
		name:              name,
		vagrantMachineDir: filepath.Join(vm.vagrantDir, name),
		verbose:           vm.verbose,
		config:            config,
	}
}

type VagrantMachine struct {
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

func (v *VagrantMachine) Up(forceDeleteIfNecessary bool) error {
	if v.config == nil {
		return fmt.Errorf("vagrantmachine requires config when Up")
	}
	if err := v.ensureVagrantFiles(); err != nil {
		return err
	}
	log.Infof("vagrantmachine(%s): ready to launch machine\n", v.name)
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
	log.Infof("vagrantmachine(%s): export kubecfg to path:%s\n", v.name, path)
	cli, err := v.NewVagrantCli()
	if err != nil {
		return err
	}
	return cli.Scp("/home/vagrant/.kube/config", path)
}

func (v *VagrantMachine) Destroy(force bool) error {
	log.Infof("vagrantmachine(%s): ready to destroy\n", v.name)
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
		Status:  status,
	}, nil
}

func (v *VagrantMachine) Name() string {
	return v.name
}

func (v *VagrantMachine) ensureVagrantFiles() error {
	// only check Vagrantfile
	f := filepath.Join(v.vagrantMachineDir, "Vagrantfile")
	if _, err := os.Stat(f); os.IsNotExist(err) {
		log.Infof("vagrantmachine(%s): prepare files under %s\n", v.name, v.vagrantMachineDir)
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
	log.Infof("vagrantmachine(%s): get port (%d,%d) for ssh and kubeapi\n", v.name, sshport, kubeport)
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

			m := vm.NewMachine(machineName, nil)
			machines = append(machines, m)
		}
	}
	return machines, nil
}
