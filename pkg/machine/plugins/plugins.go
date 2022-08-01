package plugins

import (
	"errors"
	"path/filepath"

	"github.com/footprintai/multikf/pkg/machine"
	kubeflowplugin "github.com/footprintai/multikf/pkg/machine/plugins/kubeflow"
	templatefs "github.com/footprintai/multikf/pkg/template/fs"
)

type TypePlugin string

const (
	TypePluginKubeflow TypePlugin = "kubeflow"
)

type TypePluginVersion string

func (t TypePluginVersion) String() string {
	return string(t)
}

type kubeflowTemplateMakerFunc func() *kubeflowplugin.KubeflowFileTemplate

var (
	noVersion         = NewTypePluginVersion("v0.0.0")
	availableVersions = map[TypePluginVersion]kubeflowTemplateMakerFunc{
		NewTypePluginVersion("v1.4"):   kubeflowplugin.NewKubeflow14Template,
		NewTypePluginVersion("v1.5.1"): kubeflowplugin.NewKubeflow15Template,
	}
)

func NewTypePluginVersion(s string) TypePluginVersion {
	return TypePluginVersion(s)
}

func KubeflowPluginVersionTemplate(s TypePluginVersion) (TypePluginVersion, *kubeflowplugin.KubeflowFileTemplate) {
	templateMaker, hasVersion := availableVersions[s]
	if !hasVersion {
		return noVersion, nil
	}
	return s, templateMaker()
}

type Plugin interface {
	PluginType() TypePlugin
	PluginVersion() TypePluginVersion
}

func AddPlugins(m machine.MachineCURD, plugins ...Plugin) error {
	pluginAndFiles := map[TypePlugin]string{}

	memFs := templatefs.NewMemoryFilesFs()
	for _, plugin := range plugins {
		switch plugin.PluginType() {
		case TypePluginKubeflow:
			// handle kubeflow plugins
			_, tmpl := KubeflowPluginVersionTemplate(plugin.PluginVersion())
			if tmpl == nil {
				return errors.New("plugins: no version found")
			}
			if err := memFs.Generate(plugin, tmpl); err != nil {
				return err
			}
			pluginAndFiles[plugin.PluginType()] = tmpl.Filename()
		default:
			return errors.New("plugins: no available plugins")
		}
	}
	if err := templatefs.NewFolder(m.HostDir()).DumpFiles(true, memFs.FS()); err != nil {
		return err
	}
	var err error
	_, hasKf := pluginAndFiles[TypePluginKubeflow]
	if hasKf {
		err = m.GetKubeCli().InstallKubeflow(m.GetKubeConfig(), filepath.Join(m.HostDir(), pluginAndFiles[TypePluginKubeflow]))
		if err == nil {
			err = m.GetKubeCli().PatchKubeflow(m.GetKubeConfig())
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func RemovePlugins(m machine.MachineCURD, plugins ...Plugin) error {
	return errors.New("plugins: not imp")

}
