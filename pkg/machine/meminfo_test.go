package machine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/footprintai/multikf/pkg/machine/testdata"
)

func TestMemInfo(t *testing.T) {
	meminfo, err := NewMemInfoParser(testdata.MemInfo)
	assert.NoError(t, err)
	assert.EqualValues(t, "976.88 Mib", meminfo.Total())
	assert.EqualValues(t, "68.15 Mib", meminfo.Free())

}
