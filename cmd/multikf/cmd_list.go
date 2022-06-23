package multikf

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/kind/pkg/log"

	"github.com/footprintai/multikf/pkg/machine"
)

func NewListCommand(logger log.Logger, ioStreams genericclioptions.IOStreams) *cobra.Command {
	handle := func() error {
		machineNamesMap := map[string]*OutputMachineInfo{}
		machine.ForEachProvisioner(func(p machine.Provisioner) {
			vag, err := newMachineFactoryWithProvisioner(p, logger)
			if err != nil {
				logger.Errorf("list machine: failed, err:+%v\n", err)
				return
			}
			machineList, err := vag.ListMachines()
			if err != nil {
				logger.Errorf("list machine: failed, err:+%v\n", err)
				return
			}
			for _, m := range machineList {
				info, err := m.Info()
				if err != nil {
					//logger.Errorf("grap info from individual machine failed, err:+%v\n", err)
					//return
					continue
				}
				machineNamesMap[m.Name()] = &OutputMachineInfo{
					Name:       m.Name(),
					Type:       m.Type().String(),
					MachineDir: m.HostDir(),
					Status:     info.Status,
					Gpus:       fmt.Sprintf("%s", info.GpuInfo.Info()),
					KubeApi:    info.KubeApi,
					Cpus:       fmt.Sprintf("%d", info.CpuInfo.NumCPUs()),
					Memory:     fmt.Sprintf("%s/%s", info.MemInfo.Free(), info.MemInfo.Total()),
				}
			}
		})

		var csvValues [][]string
		for _, v := range machineNamesMap {
			csvValues = append(csvValues, v.Values())
		}
		var dummyRow = &OutputMachineInfo{}
		return NewFormatWriter(ioStreams.Out, Table).WriteAndClose(
			dummyRow.Headers(),
			csvValues,
		)
	}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list guest machines",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handle()
		},
	}
	return cmd
}

// OutputMachineInfo defines the output format returned for each Machine
type OutputMachineInfo struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
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
		"type",
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
		o.Type,
		o.MachineDir,
		o.Status,
		o.Gpus,
		o.KubeApi,
		o.Cpus,
		o.Memory,
	}
}
