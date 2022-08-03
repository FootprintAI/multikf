package fs

import (
	"bytes"
	"io/fs"

	"github.com/footprintai/multikf/pkg/template"

	//log "github.com/golang/glog"
	"github.com/spf13/afero"
)

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

func (t *MemoryFilesFs) Generate(config interface{}, execTemplates ...template.TemplateExecutor) error {
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
