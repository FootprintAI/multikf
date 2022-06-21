package multikf

import (
	"github.com/footprintai/multikf/pkg/machine"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/kind/pkg/log"
)

func NewDeleteCommand(logger log.Logger, ioStreams genericclioptions.IOStreams) *cobra.Command {
	var (
		provisionerStr string // provider specifies the underly privisoner for virtual machine, either docker (under host) or vagrant
		//alsoRemoveConfigFile bool
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
		if err := m.Destroy(); err != nil {
			logger.Errorf("del: delete node (%s) failed, err:%+v\n", machineName, err)
		}
		return nil
	}
	cmd := &cobra.Command{
		Use:   "delete <machine-name>",
		Short: "delete a guest machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handle(args[0])
		},
	}
	cmd.Flags().StringVar(&provisionerStr, "provisioner", "docker", "provisioner, possible value: docker and vagrant")
	return cmd
}
