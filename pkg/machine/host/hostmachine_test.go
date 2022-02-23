package host

import (
	"fmt"
	//"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostCmd(t *testing.T) {

	name := "host001"
	dir := "/var/folders/0g/nmvss1h170b8wgbkgb9csd180000gn/T/unittest4190648733"
	assert.NoError(t, os.MkdirAll(dir, os.ModePerm))
	defer os.RemoveAll(dir)
	fmt.Printf("dir:%s\n", dir)

	hostmachines := NewHostMachines(dir, true)
	m, err := hostmachines.NewMachine(name)
	assert.NoError(t, err)
	assert.NoError(t, m.Up(true))
	_, err := m.Info()
	assert.NoError(t, err)
	assert.NoError(t, m.Destroy(true))
}
