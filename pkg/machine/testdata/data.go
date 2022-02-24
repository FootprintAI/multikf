package testdata

import (
	_ "embed"
)

//go:embed gpuinfo.xml
var NvidiaSMIGpuInfo string

//go:embed meminfo.text
var MemInfo string

//go:embed cpuinfo.text
var CpuInfo string
