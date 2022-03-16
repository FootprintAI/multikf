package vagrant

import (
	"github.com/footprintai/multikf/assets"
	vagranttemplates "github.com/footprintai/multikf/pkg/machine/vagrant/template"
	templatefs "github.com/footprintai/multikf/pkg/template/fs"
	//log "github.com/golang/glog"
)

func NewVagrantFolder(folderpath string) *VagrantFolder {
	return &VagrantFolder{
		folder: templatefs.NewFolder(folderpath),
	}
}

type VagrantFolder struct {
	folder *templatefs.Folder
}

func (v *VagrantFolder) GenerateVagrantFiles(tmplConfig *vagranttemplates.TemplateFileConfig) error {
	memoryFileFs := templatefs.NewMemoryFilesFs()
	if err := memoryFileFs.Generate(tmplConfig, vagranttemplates.NewKindTemplate(), vagranttemplates.NewDefaultVagrantTemplate()); err != nil {
		return err
	}
	if err := v.folder.DumpFiles(memoryFileFs.FS(), assets.BootstrapFs); err != nil {
		return err
	}
	return nil
}
