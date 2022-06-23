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

type Plugin interface {
	PluginType() TypePlugin
}

func AddPlugins(m machine.MachineCURD, plugins ...Plugin) error {
	pluginAndFiles := map[TypePlugin]string{}

	memFs := templatefs.NewMemoryFilesFs()
	for _, plugin := range plugins {
		switch plugin.PluginType() {
		case TypePluginKubeflow:
			// handle kubeflow plugins
			temp := kubeflowplugin.NewKubeflowTemplate()
			if err := memFs.Generate(plugin, temp); err != nil {
				return err
			}
			pluginAndFiles[plugin.PluginType()] = temp.Filename()
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
	}
	if err != nil {
		return err
	}
	return nil
}

func RemovePlugins(m machine.MachineCURD, plugins ...Plugin) error {
	return errors.New("plugins: not imp")

}
