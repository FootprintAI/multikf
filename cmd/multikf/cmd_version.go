package multikf

import (
	"github.com/footprintai/multikf/pkg/version"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/kind/pkg/log"
)

func NewVersionCommand(logger log.Logger, ioStreams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "version of multikf ",
		RunE: func(cmd *cobra.Command, args []string) error {
			version.Print()
			return nil
		},
	}
	return cmd
}
