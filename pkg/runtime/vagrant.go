package runtime

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/footprintai/multikind/pkg/template"
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

func (vm *VagrantMachines) NewMachine(name string) *VagrantMachine {
	return &VagrantMachine{
		name:              name,
		vagrantMachineDir: filepath.Join(vm.vagrantDir, name),
		verbose:           vm.verbose,
	}
}

type VagrantMachine struct {
	name              string
	vagrantMachineDir string
	verbose           bool
	config            *VagrantMachineConfig
}

func (v *VagrantMachine) AddConfig(config *VagrantMachineConfig) {
	v.config = config
}

type VagrantMachineConfig struct {
	CPUs   int
	Memory int // measured in M bytes
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

func (v *VagrantMachine) Info() (status string, cpuinfo *CpuInfo, meminfo *MemInfo, reterr error) {
	cli, reterr := v.NewVagrantCli()
	if reterr != nil {
		return
	}
	meminfo, reterr = NewMemInfoParserHelper(cli.SshExec("cat /proc/meminfo"))
	if reterr != nil {
		return
	}
	cpuinfo, reterr = NewCpuInfoParserHelper(cli.SshExec("cat /proc/cpuinfo"))
	if reterr != nil {
		return
	}
	status = cli.Status()
	return
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
	sshport, err := findFreeSSHPort()
	if err != nil {
		return err
	}
	kubeport, err := findFreeKubeApiPort()
	if err != nil {
		return err
	}
	log.Infof("vagrantmachine(%s): get port (%d,%d) for ssh and kubeapi\n", sshport, kubeport)
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

func (vm *VagrantMachines) ListMachines() ([]*OutputVagrantMachine, error) {
	machineNamesMap := map[string]*OutputVagrantMachine{}
	vfs := os.DirFS(vm.vagrantDir)
	entries, err := fs.ReadDir(vfs, ".")
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			machineName := entry.Name()

			m := vm.NewMachine(machineName)
			status, cpuinfo, memInfo, err := m.Info()
			if err != nil {
				return nil, err
			}
			machineNamesMap[machineName] = &OutputVagrantMachine{
				Name:              machineName,
				VagrantMachineDir: filepath.Join(vm.vagrantDir, machineName),
				VagrantStatus:     status,
				VagrantCpus:       fmt.Sprintf("%d", cpuinfo.NumCPUs()),
				VagrantMemory:     fmt.Sprintf("%d/%d", memInfo.Free(), memInfo.Total()),
			}
		}
	}

	var out []*OutputVagrantMachine
	for _, v := range machineNamesMap {
		out = append(out, v)
	}
	return out, nil
}

// OutputVagrantMachine defines the output format returned for each VagrantMachine
type OutputVagrantMachine struct {
	Name              string `json:"name"`
	VagrantMachineDir string `json:"dir"`
	VagrantStatus     string `json:"status"`
	VagrantCpus       string `json:"cpus"`
	VagrantMemory     string `json:"memory"`
}

func (o *OutputVagrantMachine) Headers() []string {
	return []string{
		"name",
		"dir",
		"status",
		"cpus",
		"memory",
	}
}

func (o *OutputVagrantMachine) Values() []string {
	return []string{
		o.Name,
		o.VagrantMachineDir,
		o.VagrantStatus,
		o.VagrantCpus,
		o.VagrantMemory,
	}
}
