package docker

import (
	kfmanifests "github.com/footprintai/multikf/kfmanifests"
	hosttemplates "github.com/footprintai/multikf/pkg/machine/docker/template"
	templatefs "github.com/footprintai/multikf/pkg/template/fs"
	//log "github.com/golang/glog"
)

func NewHostFolder(folderpath string) *HostFolder {
	return &HostFolder{
		folder: templatefs.NewFolder(folderpath),
	}
}

type HostFolder struct {
	folder *templatefs.Folder
}

func (h *HostFolder) GenerateFiles(tmplConfig *hosttemplates.TemplateFileConfig) error {
	memoryFileFs := templatefs.NewMemoryFilesFs()
	if err := memoryFileFs.Generate(tmplConfig, hosttemplates.NewKindTemplate()); err != nil {
		return err
	}
	if err := h.folder.DumpFiles(memoryFileFs.FS(), kfmanifests.FS); err != nil {
		return err
	}
	return nil
}
