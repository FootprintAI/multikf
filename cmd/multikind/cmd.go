package multikind

import (
	goflag "flag"
	"os"
	"path/filepath"

	"github.com/footprintai/multikind/pkg/runtime"
	"github.com/footprintai/multikind/pkg/version"
	log "github.com/golang/glog"
	"github.com/spf13/cobra"
)

var (
	cpus           int    // number of cpus allocated to the vagrant
	memoryInG      int    // number of Gigabytes allocated to the vagrant
	vagrantRootDir string // vagrant root dir which containing multiple vagrant folders, each folder(i.e. $machinename) represents a single virtual machine configuration (default: ./.vagrant)
	forceDelete    bool   // force to deleted the instance (default: false)
	forceCreate    bool   // force to create the instance regardless the instance's status (default: false)
	forceOverwrite bool   // force to overwrite the existing kubeconf file
	verbose        bool   // verbose (default: true)
	kubeconfigPath string // kubeconfig path of a vagrant machine (default: ./.vagrant/$machine/kubeconfig)

	rootCmd = &cobra.Command{
		Use:   "multikind",
		Short: "a multikind cli tool",
		Long:  `multikind is a command-line tool which use vagrant and kind to provision k8s single-node cluster.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// For cobra + glog flags. Available to all subcommands.
			goflag.Parse()
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "version of vagrant machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			version.Print()
			return nil
		},
	}

	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "export kubeconfig from a vagrant machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			run := mustNewRunCmd()
			return run.Export(args[0], kubeconfigPath)
		},
	}

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "add a vagrant machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			run := mustNewRunCmd()
			return run.Add(args[0], cpus, memoryInG)
		},
	}
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete a vagrant machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			run := mustNewRunCmd()
			return run.Delete(args[0])
		},
	}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "list vagrant machines",
		RunE: func(cmd *cobra.Command, args []string) error {
			run := mustNewRunCmd()
			return run.List()
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
	//cfg := &runtime.VagrantMachineConfig{
	//	CPUs:   cpus,
	//	Memory: memoryInG * 1024, // in M egabytes
	//}
	vag := runtime.NewVagrantMachines(vagrantRootDir, verbose)
	return &runCmd{vag: vag}, nil
}

type runCmd struct {
	vag *runtime.VagrantMachines
}

func (r *runCmd) Add(name string, cpus, memoryInG int) error {
	m := r.vag.NewMachine(name)
	m.AddConfig(&runtime.VagrantMachineConfig{
		CPUs:   cpus,
		Memory: memoryInG * 1024, // in M egabytes
	})
	if err := m.Up(forceCreate); err != nil {
		log.Errorf("runcmd: add vagrant node (%s) failed, err:%+v\n", name, err)
		return err
	}
	return nil
}

func (r *runCmd) Export(name string, path string) error {
	if path == "" {
		path = filepath.Join(vagrantRootDir, name, "kubeconfig")
	}
	m := r.vag.NewMachine(name)
	if err := m.ExportKubeConfig(path, forceOverwrite); err != nil {
		log.Errorf("runcmd: export vagrant node (%s) failed, err:%+v\n", name, err)
		return err
	}
	return nil
}

func (r *runCmd) Delete(name string) error {
	m := r.vag.NewMachine(name)
	if err := m.Destroy(forceDelete); err != nil {
		log.Errorf("runcmd: delete vagrant node (%s) failed, err:%+v\n", name, err)
		return err
	}
	return nil
}

var dummyRow = &runtime.OutputVagrantMachine{}

func (r *runCmd) List() error {
	w := NewFormatWriter(os.Stdout, Table)
	var listvalues [][]string
	respList, err := r.vag.ListMachines()
	if err != nil {
		return err
	}
	for _, resp := range respList {
		listvalues = append(listvalues, resp.Values())
	}
	return w.WriteAndClose(dummyRow.Headers(), listvalues)
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

	rootCmd.PersistentFlags().StringVar(&vagrantRootDir, "dir", ".vagrant", "vagrant root dir")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", true, "verbose (default: true)")
	addCmd.Flags().IntVar(&cpus, "cpus", 1, "number of cpus allocated to the vagrant")
	addCmd.Flags().IntVar(&memoryInG, "memoryg", 1, "number of memory in gigabytes allocated to the vagrant")
	addCmd.Flags().BoolVar(&forceCreate, "f", false, "force to create instance regardless the machine status")
	deleteCmd.Flags().BoolVar(&forceDelete, "f", false, "force remove vagrant instance")
	exportCmd.Flags().StringVar(&kubeconfigPath, "kubeconfig_path", "", "force remove vagrant instance")
	exportCmd.Flags().BoolVar(&forceOverwrite, "f", false, "force to overwrite the exiting file")
}
