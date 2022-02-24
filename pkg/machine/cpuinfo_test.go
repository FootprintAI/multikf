package machine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/footprintai/multikind/pkg/machine/testdata"
)

func TestCpuInfo(t *testing.T) {

	cpuinfo, err := NewCpuInfoParser(testdata.CpuInfo)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cpuinfo.NumCPUs())

}
