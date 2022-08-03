package multikf

import (
	"github.com/footprintai/multikf/pkg/machine"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/kind/pkg/log"
)

func NewConnectCommand(logger log.Logger, ioStreams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect",
		Short: "connect with machine via port-forward",
	}

	cmd.AddCommand(newConnectKubeflowCommand(logger, ioStreams))
	return cmd
}

func newConnectKubeflowCommand(logger log.Logger, ioStreams genericclioptions.IOStreams) *cobra.Command {
	var (
		port int // dedicated port
	)
	handle := func(machineName string) error {
		m, err := findMachineByName(machineName, logger)
		if err != nil {
			return err
		}
		destPort := port
		if !(destPort > 1024 && destPort < 65536) {
			logger.V(0).Infof("invaid customized port, use random\n")
			destPort, err = machine.FindFreePort()
			if err != nil {
				return err
			}
		}
		logger.V(0).Infof("now you can open http://localhost:%d\n", destPort)
		return m.GetKubeCli().Portforward(m.GetKubeConfig(), "svc/istio-ingressgateway", "istio-system", 80, destPort)
	}
	cmd := &cobra.Command{
		Use:   "kubeflow",
		Short: "connect with kubeflow via port-forward",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handle(args[0])
		},
	}

	cmd.Flags().IntVar(&port, "port", 0, "customized port number for connect, ranged should be 65535> >1024, default is 0 (random)")
	return cmd
}
