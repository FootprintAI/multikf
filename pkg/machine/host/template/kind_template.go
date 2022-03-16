package template

import (
	"fmt"
	"html/template"
	"io"

	pkgtemplate "github.com/footprintai/multikf/pkg/template"
)

func NewKindTemplate() *KindFileTemplate {
	return &KindFileTemplate{
		kindFileTemplate: kindDefaultFileTemplate,
	}
}

func (k *KindFileTemplate) Filename() string {
	return "kind-config.yaml"
}

func (k *KindFileTemplate) Execute(w io.Writer) error {
	tmpl, err := template.New("kindconfig").Parse(k.kindFileTemplate)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(w, k); err != nil {
		return err
	}
	return nil
}

type kindConfig interface {
	pkgtemplate.NameGetter
	pkgtemplate.KubeAPIPortGetter
}

func (k *KindFileTemplate) Populate(v interface{}) error {
	if _, isKindConfiger := v.(kindConfig); !isKindConfiger {
		return fmt.Errorf("not implements kindConfig interface")
	}
	c := v.(kindConfig)
	k.Name = c.GetName()
	k.KubeAPIPort = c.GetKubeAPIPort()
	return nil
}

type KindFileTemplate struct {
	Name             string
	KubeAPIPort      int
	kindFileTemplate string
}

var kindDefaultFileTemplate string = `
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: {{.Name}}
nodes:
- role: control-plane
  image: kindest/node:v1.20.7
networking:
  apiServerAddress: "0.0.0.0"
  apiServerPort: {{.KubeAPIPort}}
`
