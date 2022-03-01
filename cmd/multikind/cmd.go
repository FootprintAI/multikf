package multikind

import (
	goflag "flag"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/golang/glog"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/footprintai/multikind/pkg/machine"
	_ "github.com/footprintai/multikind/pkg/machine/host"
	_ "github.com/footprintai/multikind/pkg/machine/vagrant"
	"github.com/footprintai/multikind/pkg/version"
)

var (
	cpus           int    // number of cpus allocated to the geust machine
	memoryInG      int    // number of Gigabytes allocated to the guest machine
	provisionerStr string // provider specifies the underly privisoner for virtual machine, either docker (under host) or vagrant
	guestRootDir   string // root dir which containing multiple guest machines, each folder(i.e. $machinename) represents a single virtual machine configuration (default: ./.multilind)
	forceDelete    bool   // force to deleted the instance (default: false)
	forceCreate    bool   // force to create the instance regardless the instance's status (default: false)
	forceOverwrite bool   // force to overwrite the existing kubeconf file
	verbose        bool   // verbose (default: true)
	kubeconfigPath string // kubeconfig path of a guest machine (default: ./.mulitkind/$machine/kubeconfig)

	rootCmd = &cobra.Command{
		Use:   "multikind",
		Short: "a multikind cli tool",
		Long:  `multikind is a command-line tool which use vagrant and docker to provision k8s single-node cluster.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// For cobra + glog flags. Available to all subcommands.
			goflag.Parse()
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "version of multikind ",
		RunE: func(cmd *cobra.Command, args []string) error {
			version.Print()
			return nil
		},
	}

	exportCmd = &cobra.Command{
		Use:   "export <machine-name>",
		Short: "export kubeconfig from a guest machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			run := mustNewRunCmd()
			return run.Export(args[0], kubeconfigPath)
		},
	}

	addCmd = &cobra.Command{
		Use:   "add <machine-name>",
		Short: "add a guest machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			run := mustNewRunCmd()
			return run.Add(args[0], cpus, memoryInG)
		},
	}
	deleteCmd = &cobra.Command{
		Use:   "delete <machine-name>",
		Short: "delete a guest machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			run := mustNewRunCmd()
			return run.Delete(args[0])
		},
	}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "list guest machines",
		RunE: func(cmd *cobra.Command, args []string) error {
			run := mustNewRunCmd()
			return run.List()
		},
	}
	kubeflowCmd = &cobra.Command{
		Use:   "kubeflow command",
		Short: "kubeflow command line",
	}
	connectCmd = &cobra.Command{
		Use:   "connect",
		Short: "connect",
		RunE: func(cmd *cobra.Command, args []string) error {
			run := mustNewRunCmd()
			return run.ConnectKubeflow(args[0])
		},
	}
)

func mustNewRunCmd() *runCmd {
	cmd, err := newRunCmd()
	if err != nil {
		panic(err)
	}
	return cmd
}

func newRunCmd() (*runCmd, error) {
	p, err := machine.ParseProvisioner(provisionerStr)
	if err != nil {
		return nil, err
	}
	vag, err := machine.NewMachineFactory(p, guestRootDir, verbose)
	if err != nil {
		return nil, err
	}
	return &runCmd{vag: vag}, nil
}

type runCmd struct {
	vag machine.MachinesCURD
}

type machineConfig struct {
	cpus      int
	memoryInG int
}

func (m machineConfig) GetCPUs() int {
	return m.cpus
}

func (m machineConfig) GetMemory() int {
	return m.memoryInG
}

func (r *runCmd) Add(name string, cpus, memoryInG int) error {
	m, err := r.vag.NewMachine(name, machineConfig{cpus: cpus, memoryInG: memoryInG})
	if err != nil {
		return err
	}
	if err := m.Up(forceCreate); err != nil {
		log.Errorf("runcmd: add node (%s) failed, err:%+v\n", name, err)
		return err
	}
	return nil
}

func (r *runCmd) Export(name string, path string) error {
	if path == "" {
		path = filepath.Join(guestRootDir, name, "kubeconfig")
	}
	m, err := r.vag.NewMachine(name, nil)
	if err != nil {
		return err
	}
	if err := m.ExportKubeConfig(path, forceOverwrite); err != nil {
		log.Errorf("runcmd: export node (%s) failed, err:%+v\n", name, err)
		return err
	}
	return nil
}

func (r *runCmd) Delete(name string) error {
	m, err := r.vag.NewMachine(name, nil)
	if err != nil {
		return err
	}
	if err := m.Destroy(forceDelete); err != nil {
		log.Errorf("runcmd: delete node (%s) failed, err:%+v\n", name, err)
		return err
	}
	return nil
}

// OutputMachineInfo defines the output format returned for each Machine
type OutputMachineInfo struct {
	Name       string `json:"name"`
	MachineDir string `json:"dir"`
	Status     string `json:"status"`
	Cpus       string `json:"cpus"`
	Gpus       string `json:"gpus"`
	KubeApi    string `json:"kubeAPI"`
	Memory     string `json:"memory"`
}

func (o *OutputMachineInfo) Headers() []string {
	return []string{
		"name",
		"dir",
		"status",
		"gpus",
		"kubeAPI",
		"cpus",
		"memory",
	}
}

func (o *OutputMachineInfo) Values() []string {
	return []string{
		o.Name,
		o.MachineDir,
		o.Status,
		o.Gpus,
		o.KubeApi,
		o.Cpus,
		o.Memory,
	}
}

var dummyRow = &OutputMachineInfo{}

func (r *runCmd) List() error {
	w := NewFormatWriter(os.Stdout, Table)
	machineList, err := r.vag.ListMachines()
	if err != nil {
		return err
	}
	machineNamesMap := map[string]*OutputMachineInfo{}
	for _, m := range machineList {
		info, err := m.Info()
		if err != nil {
			return err
		}
		machineNamesMap[m.Name()] = &OutputMachineInfo{
			Name:       m.Name(),
			MachineDir: m.HostDir(),
			Status:     info.Status,
			Gpus:       fmt.Sprintf("%s", info.GpuInfo.Info()),
			KubeApi:    info.KubeApi,
			Cpus:       fmt.Sprintf("%d", info.CpuInfo.NumCPUs()),
			Memory:     fmt.Sprintf("%s/%s", info.MemInfo.Free(), info.MemInfo.Total()),
		}
	}

	var csvValues [][]string
	for _, v := range machineNamesMap {
		csvValues = append(csvValues, v.Values())
	}
	return w.WriteAndClose(dummyRow.Headers(), csvValues)
}

func (r *runCmd) ConnectKubeflow(name string) error {
	m, err := r.vag.NewMachine(name, nil)
	if err != nil {
		return err
	}
	destPort, err := m.Portforward("svc/istio-ingressgateway", "istio-system", 80)
	if err != nil {
		log.Errorf("runcmd: unable to connect %s failed, err:%+v\n", name, err)
		return err
	}
	log.Infof("now open http://localhost:%d", destPort)
	return nil
}

func Main() {
	defer log.Flush()

	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(kubeflowCmd)
	kubeflowCmd.AddCommand(connectCmd)

	rootCmd.PersistentFlags().StringVar(&guestRootDir, "dir", ".multikind", "multikind root dir")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", true, "verbose (default: true)")
	rootCmd.PersistentFlags().StringVar(&provisionerStr, "provisioner", "docker", "provisioner, possible value: docker and vagrant")
	addCmd.Flags().IntVar(&cpus, "cpus", 1, "number of cpus allocated to the guest machine")
	addCmd.Flags().IntVar(&memoryInG, "memoryg", 1, "number of memory in gigabytes allocated to the guest machine")
	addCmd.Flags().BoolVar(&forceCreate, "f", false, "force to create instance regardless the machine status")
	deleteCmd.Flags().BoolVar(&forceDelete, "f", false, "force remove the guest instance")
	exportCmd.Flags().StringVar(&kubeconfigPath, "kubeconfig_path", "", "force remove the guest instance")
	exportCmd.Flags().BoolVar(&forceOverwrite, "f", false, "force to overwrite the exiting file")

	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}
