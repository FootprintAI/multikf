package runtime

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/footprintai/multikind/assets"
	"github.com/footprintai/multikind/pkg/template"
	//log "github.com/golang/glog"
	"github.com/spf13/afero"
)

func NewDefaultTemplates() []template.TemplateExecutor {
	return []template.TemplateExecutor{
		template.NewKindTemplate(),
		template.NewDefaultVagrantTemplate(),
	}
}

func NewMemoryFilesFs() *MemoryFilesFs {
	return &MemoryFilesFs{
		vfs: afero.NewMemMapFs(),
	}
}

// MemoryFilesFs holds a set of template files
type MemoryFilesFs struct {
	vfs afero.Fs
}

func (t *MemoryFilesFs) FS() fs.FS {
	return afero.NewIOFS(t.vfs)
}

func (t *MemoryFilesFs) Generate(config *template.TemplateFileConfig, additionalTemplates ...template.TemplateExecutor) error {
	execTemplates := NewDefaultTemplates()
	execTemplates = append(execTemplates, additionalTemplates...)
	for _, exec := range execTemplates {
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

func (v *VagrantFolder) GenerateVagrantFiles(tmplConfig *template.TemplateFileConfig) error {
	memoryFileFs := NewMemoryFilesFs()
	if err := memoryFileFs.Generate(tmplConfig); err != nil {
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
