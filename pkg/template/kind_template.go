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

type KindConfiger interface {
	NameGetter
	KubeAPIPortGetter
	KubeAPIIPGetter
	GpuGetter
	ExportPortsGetter
	AuditEnabler
	WorkersGetter
	NodeLabelsGetter
	LocalPathGetter
}

func (k *KindFileTemplate) Populate(v interface{}) error {
	if _, isKindConfiger := v.(KindConfiger); !isKindConfiger {
		return fmt.Errorf("not implements kindConfig interface")
	}
	c := v.(KindConfiger)
	k.Name = c.GetName()
	k.KubeAPIPort = c.GetKubeAPIPort()
	k.KubeAPIIP = c.GetKubeAPIIP()
	k.UseGPU = c.GetGPUs() > 0
	k.ExportPorts = c.GetExportPorts()
	k.AuditEnabled = c.AuditEnabled()
	k.AuditFileAbsolutePath = c.AuditFileAbsolutePath()
	k.LocalPath = c.LocalPath()
	k.Workers = c.GetWorkers()

	nodeLabels := c.GetNodeLabels()
	k.NodeLabels = make([]string, len(nodeLabels), len(nodeLabels))
	for idx := 0; idx < len(nodeLabels); idx++ {
		k.NodeLabels[idx] = fmt.Sprintf("%s=%s", nodeLabels[idx].Key, nodeLabels[idx].Value)
	}

	return nil
}

type KindFileTemplate struct {
	Name                  string
	KubeAPIIP             string
	KubeAPIPort           int
	UseGPU                bool
	kindFileTemplate      string
	ExportPorts           []machine.ExportPortPair
	AuditEnabled          bool
	AuditFileAbsolutePath string
	LocalPath             string
	Workers               []Worker
	NodeLabels            []string
}

var kindDefaultFileTemplate string = `
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: {{.Name}}
nodes:
- role: control-plane
  kubeadmConfigPatches:
  {{- if .AuditEnabled}}
  - |
    kind: ClusterConfiguration
    apiServer:
      # enable auditing flags on the API server
      extraArgs:
        audit-log-path: /var/log/kubernetes/kube-apiserver-audit.log
        audit-policy-file: /etc/kubernetes/policies/audit-policy.yaml
        audit-log-maxage: "30"
        audit-log-maxbackup: "10"
        audit-log-maxsize: "100"
      # mount new files / directories on the control plane
      extraVolumes:
        - name: audit-policies
          hostPath: /etc/kubernetes/policies
          mountPath: /etc/kubernetes/policies
          readOnly: true
          pathType: "DirectoryOrCreate"
        - name: "audit-logs"
          hostPath: "/var/log/kubernetes"
          mountPath: "/var/log/kubernetes"
          readOnly: false
          pathType: DirectoryOrCreate
  {{- end}}
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
        {{- range $i, $p := .NodeLabels}}
        node-labels: "{{$p}}"
        {{- end}}
  image: kindest/node:v1.23.17@sha256:e5fd1d9cd7a9a50939f9c005684df5a6d145e8d695e78463637b79464292e66c
  gpus: {{.UseGPU}}
  {{if .ExportPorts}}extraPortMappings:{{end}}
  {{- range $i, $p := .ExportPorts}}
  - containerPort: {{ $p.ContainerPort }}
    hostPort: {{ $p.HostPort }}
    protocol: TCP
  {{- end}}
  {{- if or .AuditEnabled .LocalPath}}
  extraMounts:
  {{- if or .AuditEnabled }}
  - hostPath: {{.AuditFileAbsolutePath}}
    containerPath: /etc/kubernetes/policies/audit-policy.yaml
    readOnly: true
  {{- end}}
  {{- if ne .LocalPath ""}}
  - hostPath: {{.LocalPath}}
    containerPath: /var/local-path-provisioner
  {{- end}}
  {{- end}}
{{- range .Workers }}
- role: worker
  image: kindest/node:v1.23.17@sha256:e5fd1d9cd7a9a50939f9c005684df5a6d145e8d695e78463637b79464292e66c
  gpus: {{ .UseGPU}}
  {{- if ne .LocalPath ""}}
  extraMounts:
  - hostPath: {{.LocalPath}}
    containerPath: /var/local-path-provisioner
  {{- end}}
{{- end}}
networking:
  apiServerAddress: {{.KubeAPIIP}}
  apiServerPort: {{.KubeAPIPort}}
`
