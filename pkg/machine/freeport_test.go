package machine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFreePorts(t *testing.T) {
	sshport, err := FindFreeSSHPort()
	assert.NoError(t, err)
	assert.True(t, sshport >= 2022)

	kubeport, err := FindFreeKubeApiPort()
	assert.NoError(t, err)
	assert.True(t, kubeport >= 16443)
}
