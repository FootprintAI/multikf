package vagrant

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/footprintai/multikind/assets"
	vagranttemplates "github.com/footprintai/multikind/pkg/machine/vagrant/template"
	"github.com/footprintai/multikind/pkg/template"
	templatefs "github.com/footprintai/multikind/pkg/template/fs"

	//log "github.com/golang/glog"
	"github.com/spf13/afero"
)

func NewDefaultTemplates() []template.TemplateExecutor {
	return []template.TemplateExecutor{
		vagranttemplates.NewKindTemplate(),
		vagranttemplates.NewDefaultVagrantTemplate(),
	}
}

func NewVagrantFolder(folderPath string) *VagrantFolder {
	return &VagrantFolder{
		folderFs:   afero.NewBasePathFs(afero.NewOsFs(), folderPath),
		folderPath: folderPath,
		//forceOverwrite: false,
	}
}

type VagrantFolder struct {
	folderFs   afero.Fs
	folderPath string
	//forceOverwrite bool
}

func (v *VagrantFolder) Root() string {
	return v.folderPath
}

func (v *VagrantFolder) ensureFolder() error {
	_, err := v.folderFs.Stat(".")
	if os.IsNotExist(err) {
		if err := v.folderFs.MkdirAll(".", 0777); err != nil {
			return err
		}
	}
	return nil
}

func (v *VagrantFolder) GenerateVagrantFiles(tmplConfig *vagranttemplates.TemplateFileConfig) error {
	memoryFileFs := templatefs.NewMemoryFilesFs()
	if err := memoryFileFs.Generate(tmplConfig, NewDefaultTemplates()...); err != nil {
		return err
	}
	if err := v.dumpFiles(memoryFileFs.FS(), assets.BootstrapFs); err != nil {
		return err
	}
	return nil
}

// dumpFiles dumps virtualFs into a real file folder which should be the home dir of a vagrant project
func (v *VagrantFolder) dumpFiles(virtualFs ...fs.FS) error {
	if err := v.ensureFolder(); err != nil {
		return err
	}

	for _, vfs := range virtualFs {
		err := fs.WalkDir(vfs, ".", func(path string, d fs.DirEntry, err error) error {
			if path == "." {
				return nil
			}
			if d.IsDir() {
				fmt.Printf("D %s\n", path)
				//os.Mkdir(path)
				if err := v.folderFs.MkdirAll(path, 0777); err != nil {
					return err
				}
				return nil
			}
			fmt.Printf("F %s\n", path)
			fd, err := vfs.Open(path)
			if err != nil {
				return err
			}
			b, err := ioutil.ReadAll(fd)
			if err != nil {
				return err
			}
			fd.Close()
			if err := afero.WriteFile(v.folderFs, path, b, 0666); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
