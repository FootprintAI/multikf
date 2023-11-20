package vagrant

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	vagranttemplates "github.com/footprintai/multikf/pkg/machine/vagrant/template"
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
	assert.NoError(t, vdir.GenerateVagrantFiles(vagranttemplates.NewVagrantTemplateConfig(
		"unittest",
		2,
		1026,
		1234,
		5678,
		"1.2.3.4",
		0,
		nil,
		false,
		"",
		0,
		nil,
		"",
	),
	))

	expectedFiles := []string{
		"Vagrantfile",
		"kind-config.yaml",
		"bootstrap/bootstrap.sh",
		"bootstrap/provision-cluster.sh",
	}
	for _, expectedfile := range expectedFiles {
		_, err := mockFs.Stat(expectedfile)
		if os.IsNotExist(err) {
			t.Errorf("file \"%s\" does not exist.\n", expectedfile)
		}
	}
}
