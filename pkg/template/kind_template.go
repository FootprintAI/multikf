package template

import (
	"fmt"
	"html/template"
	"io"

	"github.com/footprintai/multikf/pkg/machine"
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
	NameGetter
	KubeAPIPortGetter
	KubeAPIIPGetter
	GpuGetter
	ExportPortsGetter
}

func (k *KindFileTemplate) Populate(v interface{}) error {
	if _, isKindConfiger := v.(kindConfig); !isKindConfiger {
		return fmt.Errorf("not implements kindConfig interface")
	}
	c := v.(kindConfig)
	k.Name = c.GetName()
	k.KubeAPIPort = c.GetKubeAPIPort()
	k.KubeAPIIP = c.GetKubeAPIIP()
	k.UseGPU = c.GetGPUs() > 0
	k.ExportPorts = c.GetExportPorts()

	return nil
}

type KindFileTemplate struct {
	Name             string
	KubeAPIIP        string
	KubeAPIPort      int
	UseGPU           bool
	kindFileTemplate string
	ExportPorts      []machine.ExportPortPair
}

var kindDefaultFileTemplate string = `
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: {{.Name}}
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  image: kindest/node:v1.21.2
  gpus: {{.UseGPU}}
  {{if .ExportPorts}}extraPortMappings:{{end}}
  {{- range $i, $p := .ExportPorts}}
  - containerPort: {{ $p.ContainerPort }}
    hostPort: {{ $p.HostPort }}
    protocol: TCP
  {{- end}}
networking:
  apiServerAddress: {{.KubeAPIIP}}
  apiServerPort: {{.KubeAPIPort}}
`
