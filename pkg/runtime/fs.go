package runtime

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"

	"github.com/footprintai/multikind/pkg/template"
	"github.com/spf13/afero"
)

func NewDefaultTemplates() []template.TemplateExecutor {
	return []template.TemplateExecutor{
		template.NewKindTemplate(),
		template.NewDefaultVagrantTemplate(),
	}
}

func NewTemplateFilesFs() *TemplateFilesFs {
	return &TemplateFilesFs{
		vfs: afero.NewMemMapFs(),
	}
}

// TemplateFilesFs holds a set of template files
type TemplateFilesFs struct {
	vfs afero.Fs
}

func (t *TemplateFilesFs) FS() fs.FS {
	return afero.NewIOFS(t.vfs)
}

func (t *TemplateFilesFs) Generate(config *template.TemplateFileConfig, execs []template.TemplateExecutor) error {
	for _, exec := range execs {
		if err := exec.Populate(config); err != nil {
			return err
		}
		bytesInFile := &bytes.Buffer{}
		if err := exec.Execute(bytesInFile); err != nil {
			return err
		}
		if err := afero.WriteFile(t.vfs, exec.Filename(), bytesInFile.Bytes(), 0666); err != nil {
			return err
		}
	}
	return nil
}

func NewVagrantFolder(folderPath string) *VagrantFolder {
	return &VagrantFolder{
		folderFs: afero.NewBasePathFs(afero.NewOsFs(), folderPath),
	}
}

type VagrantFolder struct {
	folderFs afero.Fs
}

// DumpFiles dumps virtualFs into a real file folder which should be the home dir of a vagrant project
func (v *VagrantFolder) DumpFiles(virtualFs ...fs.FS) error {
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
