package multikf

import (
	goflag "flag"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/log"

	_ "github.com/footprintai/multikf/pkg/machine/docker"
	_ "github.com/footprintai/multikf/pkg/machine/vagrant"
)

var (
	guestRootDir string // root dir which containing multiple guest machines, each folder(i.e. $machinename) represents a single virtual machine configuration (default: ./.multikfdir)
	verbose      bool   // verbose (default: true)
)

func NewRootCommand(logger log.Logger, ioStreams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multikf",
		Short: "a multikf cli tool",
		Long:  `multikf is a command-line tool which use vagrant and docker to provision Kubernetes and kubeflow single-node cluster.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// For cobra + glog flags. Available to all subcommands.
			goflag.Parse()
		},
	}
	cmd.AddCommand(NewVersionCommand(logger, ioStreams))
	cmd.AddCommand(NewAddCommand(logger, ioStreams))
	cmd.AddCommand(NewListCommand(logger, ioStreams))
	cmd.AddCommand(NewDeleteCommand(logger, ioStreams))
	cmd.AddCommand(NewConnectCommand(logger, ioStreams))
	cmd.AddCommand(NewPluginCommand(logger, ioStreams))

	cmd.PersistentFlags().StringVar(&guestRootDir, "dir", ".multikfdir", "multikf root dir")
	cmd.PersistentFlags().BoolVar(&verbose, "verbose", true, "verbose (default: true)")
	return cmd
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viperConfigKeyRootDir.Set(guestRootDir)
	viperConfigKeyVerbose.Set(verbose)
}

func Main() {

	logger := cmd.NewLogger()
	if verbose {
		type verbosity interface {
			SetVerbosity(verbosity log.Level)
		}
		_, ok := logger.(verbosity)
		if ok {
			logger.(verbosity).SetVerbosity(log.Level(1))
		}
	}

	NewRootCommand(
		logger,
		genericclioptions.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		},
	).Execute()
}
