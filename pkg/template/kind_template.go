package template

import (
	"html/template"
	"io"
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

func (k *KindFileTemplate) Populate(v *TemplateFileConfig) error {
	k.KubeAPIPort = v.KubeApiPort
	return nil
}

type KindFileTemplate struct {
	KubeAPIPort      int
	kindFileTemplate string
}

var kindDefaultFileTemplate string = `
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: cluster1
nodes:
- role: control-plane
  image: kindest/node:v1.20.7
networking:
  apiServerAddress: "0.0.0.0"
  apiServerPort: {{.KubeAPIPort}}
`
