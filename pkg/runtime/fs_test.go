package runtime

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/footprintai/multikind/assets"
	"github.com/footprintai/multikind/pkg/template"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestVagrantFile(t *testing.T) {
	// test with memory fs
	//mockFs := afero.NewMemMapFs()
	//vdir := &VagrantFolder{
	//	folderFs: mockFs,
	//}

	// test with os.fs
	tmpdir, err := ioutil.TempDir("", "unittest")
	assert.NoError(t, err)
	fmt.Printf("tmpdir:%s\n", tmpdir)
	mockFs := afero.NewBasePathFs(afero.NewOsFs(), tmpdir)
	vdir := NewVagrantFolder(tmpdir)

	templateFs := NewTemplateFilesFs()
	assert.NoError(t, templateFs.Generate(&template.TemplateFileConfig{
		Name:        "unittest",
		CPUs:        2,
		Memory:      1026,
		SSHPort:     1234,
		KubeApiPort: 5678,
	}, NewDefaultTemplates()))

	assert.NoError(t, vdir.DumpFiles(templateFs.FS(), assets.BootstrapFs))

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
