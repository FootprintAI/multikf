package docker

import (
	"fmt"
	"os"

	//"os"
	"testing"

	"github.com/footprintai/multikf/pkg/k8s"
	"github.com/footprintai/multikf/pkg/machine"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kind/pkg/cmd"
)

func TestHostCmd(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	logger := cmd.NewLogger()

	name := "host001"
	dir := "/tmp/unittest4190648733"
	assert.NoError(t, os.MkdirAll(dir, os.ModePerm))
	//defer os.RemoveAll(dir)
	fmt.Printf("dir:%s\n", dir)

	hostmachines := NewHostMachines(logger, dir, true)
	m, err := hostmachines.NewMachine(name, noConfigurer{})
	assert.NoError(t, err)
	assert.NoError(t, m.Up())
	_, err = m.Info()
	assert.NoError(t, err)
	assert.NoError(t, m.Destroy())
}

type noConfigurer struct{}

var (
	_ machine.MachineConfiger = noConfigurer{}
)

func (n noConfigurer) Info() string {
	return ""
}

func (n noConfigurer) GetCPUs() int {
	return 1
}

func (n noConfigurer) GetMemory() int {
	return 4
}

func (n noConfigurer) GetGPUs() int {
	return 0
}

func (n noConfigurer) GetKubeAPIIP() string {
	return "0.0.0.0"
}

func (n noConfigurer) GetExportPorts() []machine.ExportPortPair {
	return nil
}

func (n noConfigurer) GetForceOverwriteConfig() bool {
	return false
}

func (n noConfigurer) AuditEnabled() bool {
	return false
}

func (n noConfigurer) GetWorkers() int {
	return 0
}

func (n noConfigurer) GetNodeVersion() k8s.KindK8sVersion {
	return k8s.DefaultVersion()

}

func (n noConfigurer) GetNodeLabels() []machine.NodeLabel {
	return nil
}

func (n noConfigurer) GetLocalPath() string {
	return ""
}
