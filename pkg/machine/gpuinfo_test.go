package machine

import (
	"testing"

	"github.com/footprintai/multikind/pkg/machine/testdata"
	"github.com/stretchr/testify/assert"
)

func TestGpuInfo(t *testing.T) {
	gpuinfo, err := NewGpuInfoParser(testdata.NvidiaSMIGpuInfo)
	assert.NoError(t, err)
	assert.EqualValues(t, "450.36.06", gpuinfo.DriverVersion)
	assert.EqualValues(t, 2, len(gpuinfo.Gpus))
	assert.EqualValues(t, "Tesla V100-PCIE-16GB", gpuinfo.Gpus[0].GpuName)
	assert.EqualValues(t, "16160 MiB", gpuinfo.Gpus[0].GpuMemInfo.Free())
	assert.EqualValues(t, "16160 MiB", gpuinfo.Gpus[0].GpuMemInfo.Total())

	assert.EqualValues(t, "Tesla V100-PCIE-16GB", gpuinfo.Gpus[1].GpuName)
	assert.EqualValues(t, "16160 MiB", gpuinfo.Gpus[1].GpuMemInfo.Free())
	assert.EqualValues(t, "16160 MiB", gpuinfo.Gpus[1].GpuMemInfo.Total())
}
