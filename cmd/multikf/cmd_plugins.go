package multikf

import (
	"github.com/footprintai/multikf/pkg/machine/plugins"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/kind/pkg/log"
)

func NewPluginCommand(logger log.Logger, ioStreams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "plugin",
		Long:  `enable plugins to the underlying k8s`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.AddCommand(newAddPluginCommand(logger, ioStreams))
	cmd.AddCommand(newRemovePluginCommand(logger, ioStreams))
	return cmd
}

func newAddPluginCommand(logger log.Logger, ioStreams genericclioptions.IOStreams) *cobra.Command {
	var (
		withKubeflow                bool   // install with kubeflow components
		withKubeflowVersion         string // with kubeflow version
		withKubeflowDefaultPassword string // with kubeflow defaultpassword
	)

	handle := func(machineName string) error {
		m, err := findMachineByName(machineName, logger)
		if err != nil {
			return err
		}
		var installedPlugins []plugins.Plugin
		if withKubeflow {
			installedPlugins = append(installedPlugins,
				kubeflowPlugin{
					withKubeflowDefaultPassword: withKubeflowDefaultPassword,
					kubeflowVersion:             plugins.NewTypePluginVersion(withKubeflowVersion),
				},
			)
		}
		return plugins.AddPlugins(m, installedPlugins...)
	}
	cmd := &cobra.Command{
		Use:   "add <machine-name> --with_kubeflow",
		Short: "add kubeflow plugins to the machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handle(args[0])
		},
	}

	cmd.Flags().BoolVar(&withKubeflow, "with_kubeflow", true, "install kubeflow modules (default: true)")
	cmd.Flags().StringVar(&withKubeflowVersion, "kubeflow_version", "v1.4", "kubeflow version v1.4/v1.5.1")
	cmd.Flags().StringVar(&withKubeflowDefaultPassword, "with_password", "12341234", "with a specific password for default user (default: 12341234)")

	return cmd
}

func newRemovePluginCommand(logger log.Logger, ioStreams genericclioptions.IOStreams) *cobra.Command {
	var (
		removeKubeflow bool // remove kubeflow components
	)

	handle := func(machineName string) error {
		m, err := findMachineByName(machineName, logger)
		if err != nil {
			return err
		}
		var removingPlugins []plugins.Plugin
		if removeKubeflow {
			removingPlugins = append(removingPlugins, kubeflowPlugin{})
		}
		return plugins.RemovePlugins(m, removingPlugins...)
	}
	cmd := &cobra.Command{
		Use:   "remove <machine-name> --remove_kubeflow",
		Short: "rmeove kubeflow plugins to the machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handle(args[0])
		},
	}

	cmd.Flags().BoolVar(&removeKubeflow, "remove_kubeflow", false, "remove kubeflow modules (default: false)")
	return cmd
}
