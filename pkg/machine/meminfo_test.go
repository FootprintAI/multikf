package machine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemInfo(t *testing.T) {
	meminfo, err := NewMemInfoParser(meminfo)
	assert.NoError(t, err)
	assert.EqualValues(t, 1000328, meminfo.Total())
	assert.EqualValues(t, 69784, meminfo.Free())

}

var meminfo = `
MemTotal:        1000328 kB
MemFree:           69784 kB
MemAvailable:     144376 kB
Buffers:            8872 kB
Cached:           185232 kB
SwapCached:            0 kB
Active:           680748 kB
Inactive:          93604 kB
Active(anon):     584056 kB
Inactive(anon):     6268 kB
Active(file):      96692 kB
Inactive(file):    87336 kB
Unevictable:       18764 kB
Mlocked:           18764 kB
SwapTotal:             0 kB
SwapFree:              0 kB
Dirty:               256 kB
Writeback:             0 kB
AnonPages:        599052 kB
Mapped:           146652 kB
Shmem:              8208 kB
KReclaimable:      36312 kB
Slab:             106272 kB
SReclaimable:      36312 kB
SUnreclaim:        69960 kB
KernelStack:        6220 kB
PageTables:         6264 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:      500164 kB
Committed_AS:    3846388 kB
VmallocTotal:   34359738367 kB
VmallocUsed:       14420 kB
VmallocChunk:          0 kB
Percpu:             1048 kB
HardwareCorrupted:     0 kB
AnonHugePages:     12288 kB
ShmemHugePages:        0 kB
ShmemPmdMapped:        0 kB
FileHugePages:         0 kB
FilePmdMapped:         0 kB
CmaTotal:              0 kB
CmaFree:               0 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
Hugetlb:               0 kB
DirectMap4k:      118720 kB
DirectMap2M:      929792 kB
`
