package plugins

import (
	"errors"
	"path/filepath"

	kfmanifests "github.com/footprintai/multikf/kfmanifests"
	"github.com/footprintai/multikf/pkg/machine"
	kubeflowplugin "github.com/footprintai/multikf/pkg/machine/plugins/kubeflow"
	"github.com/footprintai/multikf/pkg/template"
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

type templateMakerFunc func() template.TemplateExecutor

var (
	noVersion             = NewTypePluginVersion("v0.0.0")
	typePluginVersionKF14 = NewTypePluginVersion("v1.4")
)

func NewTypePluginVersion(s string) TypePluginVersion {
	return TypePluginVersion(s)
}

func kubeflowPluginVersionTemplate(s TypePluginVersion) (TypePluginVersion, template.TemplateExecutor) {
	version := s.String()
	manifests, err := kfmanifests.GetVersion(version)
	if err != nil {
		return noVersion, nil
	}
	return s, kubeflowplugin.NewKubeflowTemplateExecutor(kfmanifests.VersionBaseFileName(version), manifests)
}

type Plugin interface {
	PluginType() TypePlugin
	PluginVersion() TypePluginVersion
}

type TypeHostFilePath string

func (t TypeHostFilePath) String() string {
	return string(t)
}

func NewTypeHostFilePath(s string) TypeHostFilePath {
	return TypeHostFilePath(s)
}

func generatePluginsManifestsMapping(m machine.MachineCURD, dumpToFile bool, plugins ...Plugin) (map[Plugin]template.TemplateExecutor, error) {
	pluginAndTmpls := map[Plugin]template.TemplateExecutor{}
	for _, plugin := range plugins {
		switch plugin.PluginType() {
		case TypePluginKubeflow:
			// handle kubeflow plugins
			_, tmpl := kubeflowPluginVersionTemplate(plugin.PluginVersion())
			if tmpl == nil {
				return nil, errors.New("plugins: no version found")
			}
			pluginAndTmpls[plugin] = tmpl
		default:
			return nil, errors.New("plugins: no available plugins")
		}
	}
	if dumpToFile {
		memFs := templatefs.NewMemoryFilesFs()
		for plugin, tmpl := range pluginAndTmpls {
			if err := memFs.Generate(plugin, tmpl); err != nil {
				return nil, err
			}
		}
		// TODO: check whether we want to overwrite exsiting or not
		if err := templatefs.NewFolder(m.HostDir()).DumpFiles(true, memFs.FS()); err != nil {
			return nil, err
		}
	}
	return pluginAndTmpls, nil
}

func AddPlugins(m machine.MachineCURD, plugins ...Plugin) error {
	var err error
	pluginAndTmpls, err := generatePluginsManifestsMapping(m, true, plugins...)
	if err != nil {
		return nil
	}
	for plugin, tmpl := range pluginAndTmpls {
		if plugin.PluginType() == TypePluginKubeflow {
			err = m.GetKubeCli().InstallKubeflow(m.GetKubeConfig(), filepath.Join(m.HostDir(), tmpl.Filename()))
			if err != nil {
				return err
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func RemovePlugins(m machine.MachineCURD, plugins ...Plugin) error {
	var err error
	pluginAndTmpls, err := generatePluginsManifestsMapping(m, false, plugins...)
	if err != nil {
		return err
	}
	for plugin, tmpl := range pluginAndTmpls {
		if plugin.PluginType() == TypePluginKubeflow {
			err = m.GetKubeCli().RemoveKubeflow(m.GetKubeConfig(), filepath.Join(m.HostDir(), tmpl.Filename()))
		}
	}
	return err

}
