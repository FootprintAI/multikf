package vagrant

import (
	"github.com/footprintai/multikf/assets"
	vagranttemplates "github.com/footprintai/multikf/pkg/machine/vagrant/template"
	pkgtemplate "github.com/footprintai/multikf/pkg/template"
	templatefs "github.com/footprintai/multikf/pkg/template/fs"
)

func NewVagrantFolder(folderpath string) *VagrantFolder {
	return &VagrantFolder{
		folder: templatefs.NewFolder(folderpath),
	}
}

type VagrantFolder struct {
	folder *templatefs.Folder
}

func (v *VagrantFolder) GenerateVagrantFiles(tmplConfig *vagranttemplates.VagrantTemplateConfig) error {
	memoryFileFs := templatefs.NewMemoryFilesFs()
	if err := memoryFileFs.Generate(tmplConfig, vagranttemplates.NewDefaultVagrantTemplate()); err != nil {
		return err
	}
	if err := memoryFileFs.Generate(tmplConfig, pkgtemplate.NewKindTemplate(), pkgtemplate.NewAuditPolicyTemplate()); err != nil {
		return err
	}
	if err := v.folder.DumpFiles(true, memoryFileFs.FS(), assets.BootstrapFs); err != nil {
		return err
	}
	return nil
}
