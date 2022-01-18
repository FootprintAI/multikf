package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFreePorts(t *testing.T) {
	sshport, err := findFreeSSHPort()
	assert.NoError(t, err)
	assert.True(t, sshport >= 2022)

	kubeport, err := findFreeKubeApiPort()
	assert.NoError(t, err)
	assert.True(t, kubeport >= 16443)
}

func TestVagrantCmd(t *testing.T) {
	// only for integration test
	name := "test001"
	config := &VagrantMachineConfig{
		CPUs:   1,
		Memory: 1024 * 2,
	}
	//vagrantDir := newEmptyDir()
	vagrantDir := "/var/folders/0g/nmvss1h170b8wgbkgb9csd180000gn/T/unittest4190648733"
	vm := NewVagrantMachines(vagrantDir, true)
	m := vm.NewMachine(name)
	m.AddConfig(config)
	cli, err := m.NewVagrantCli()
	assert.NoError(t, err)
	assert.NoError(t, cli.TryUp(true))

}
