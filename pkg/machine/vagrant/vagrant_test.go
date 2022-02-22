package vagrant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVagrantCmd(t *testing.T) {
	t.SkipNow()

	// only for integration test
	name := "test001"
	config := &VagrantMachineConfig{
		CPUs:   1,
		Memory: 1024 * 2,
	}
	//vagrantDir := newEmptyDir()
	vagrantDir := "/var/folders/0g/nmvss1h170b8wgbkgb9csd180000gn/T/unittest4190648733"
	vm := NewVagrantMachines(vagrantDir, true)
	m := vm.NewMachine(name, config)
	cli, err := m.NewVagrantCli()
	assert.NoError(t, err)
	assert.NoError(t, cli.TryUp(true))

}
