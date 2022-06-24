package multikf

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/kind/pkg/log"
)

func NewDeleteCommand(logger log.Logger, ioStreams genericclioptions.IOStreams) *cobra.Command {
	handle := func(machineName string) error {
		m, err := findMachineByName(machineName, logger)
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
	return cmd
}
