package docker

import (
	hosttemplates "github.com/footprintai/multikf/pkg/machine/docker/template"
	pkgtemplate "github.com/footprintai/multikf/pkg/template"
	templatefs "github.com/footprintai/multikf/pkg/template/fs"
)

func NewHostFolder(folderpath string) *HostFolder {
	return &HostFolder{
		folder: templatefs.NewFolder(folderpath),
	}
}

type HostFolder struct {
	folder *templatefs.Folder
}

func (h *HostFolder) GenerateFiles(tmplConfig *hosttemplates.DockerHostmachineTemplateConfig) error {
	memoryFileFs := templatefs.NewMemoryFilesFs()
	if err := memoryFileFs.Generate(tmplConfig, pkgtemplate.NewKindTemplate(), pkgtemplate.NewAuditPolicyTemplate()); err != nil {
		return err
	}
	/*
		if err := memoryFileFs.Generate(tmplConfig, hosttemplates.NewKubeflowTemplate()); err != nil {
			return err
		}
	*/
	if err := h.folder.DumpFiles(true, memoryFileFs.FS()); err != nil {
		return err
	}
	return nil
}
