package runtime

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/footprintai/multikind/pkg/template"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func newEmptyDir() string {
	tmpdir, err := ioutil.TempDir("", "unittest")
	if err != nil {
		panic(err)
	}
	os.Remove(tmpdir) // ensure the dir is not exist, only path we need
	return tmpdir
}

func TestVagrantFile(t *testing.T) {
	// test with memory fs
	//mockFs := afero.NewMemMapFs()
	//vdir := &VagrantFolder{
	//	folderFs: mockFs,
	//}

	// test with os.fs
	tmpdir := newEmptyDir()
	fmt.Printf("tmpdir:%s\n", tmpdir)
	mockFs := afero.NewBasePathFs(afero.NewOsFs(), tmpdir)
	vdir := NewVagrantFolder(tmpdir)
	assert.NoError(t, vdir.GenerateVagrantFiles(&template.TemplateFileConfig{
		Name:        "unittest",
		CPUs:        2,
		Memory:      1026,
		SSHPort:     1234,
		KubeApiPort: 5678,
	}))

	expectedFiles := []string{
		"Vagrantfile",
		"kind-config.yaml",
		"bootstrap/bootstrap.sh",
		"bootstrap/provision-cluster.sh",
		"bootstrap/provision-kf14.sh",
	}
	for _, expectedfile := range expectedFiles {
		_, err := mockFs.Stat(expectedfile)
		if os.IsNotExist(err) {
			t.Errorf("file \"%s\" does not exist.\n", expectedfile)
		}
	}
}
