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
		provisionerStr string // provider specifies the underly privisoner for virtual machine, either docker (under host) or vagrant
	)
	handle := func(machineName string) error {
		vag, err := newMachineFactoryWithProvisioner(
			machine.MustParseProvisioner(provisionerStr),
			logger,
		)
		if err != nil {
			return err
		}
		m, err := vag.NewMachine(machineName, nil)
		if err != nil {
			return err
		}
		_, err = m.Portforward("svc/istio-ingressgateway", "istio-system", 80)
		if err != nil {
			logger.Errorf("connect: unable to connect %s failed, err:%+v\n", machineName, err)
			return err
		}
		return nil
	}
	cmd := &cobra.Command{
		Use:   "kubeflow",
		Short: "connect with kubeflow via port-forward",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handle(args[0])
		},
	}

	cmd.Flags().StringVar(&provisionerStr, "provisioner", "docker", "provisioner, possible value: docker and vagrant")
	return cmd
}
