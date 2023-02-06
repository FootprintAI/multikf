package template

import (
	"bytes"
	"testing"

	"github.com/footprintai/multikf/pkg/machine"
	"github.com/stretchr/testify/assert"
)

type staticConfig struct{}

func (s staticConfig) GetName() string {
	return "staticconfig"
}

func (s staticConfig) GetKubeAPIPort() int {
	return 8443
}

func (s staticConfig) GetKubeAPIIP() string {
	return "1.2.3.4"
}

func (s staticConfig) GetGPUs() int {
	return 1
}

func (s staticConfig) AuditEnabled() bool {
	return false
}

func (s staticConfig) AuditFileAbsolutePath() string {
	return ""
}

func (s staticConfig) GetWorkerIDs() []int {
	return []int{1, 2, 3}
}

func (s staticConfig) GetNodeLabels() []machine.NodeLabel {
	return []machine.NodeLabel{
		machine.NodeLabel{
			Key:   "a",
			Value: "b",
		},
		machine.NodeLabel{
			Key:   "c",
			Value: "d",
		},
	}
}

func (s staticConfig) GetExportPorts() []machine.ExportPortPair {
	return []machine.ExportPortPair{
		machine.ExportPortPair{
			HostPort:      80,
			ContainerPort: 8081,
		},
		machine.ExportPortPair{
			HostPort:      443,
			ContainerPort: 8083,
		},
	}
}

func TestKindTemplate(t *testing.T) {
	kt := NewKindTemplate()
	assert.NoError(t, kt.Populate(staticConfig{}))
	buf := &bytes.Buffer{}
	assert.NoError(t, kt.Execute(buf))
	assert.EqualValues(t, gold, buf.String())
}

var gold = `
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: staticconfig
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
        node-labels: "a=b"
        node-labels: "c=d"
  image: kindest/node:v1.23.12@sha256:9402cf1330bbd3a0d097d2033fa489b2abe40d479cc5ef47d0b6a6960613148a
  gpus: true
  extraPortMappings:
  - containerPort: 8081
    hostPort: 80
    protocol: TCP
  - containerPort: 8083
    hostPort: 443
    protocol: TCP
- role: worker
  image: kindest/node:v1.23.12@sha256:9402cf1330bbd3a0d097d2033fa489b2abe40d479cc5ef47d0b6a6960613148a
- role: worker
  image: kindest/node:v1.23.12@sha256:9402cf1330bbd3a0d097d2033fa489b2abe40d479cc5ef47d0b6a6960613148a
- role: worker
  image: kindest/node:v1.23.12@sha256:9402cf1330bbd3a0d097d2033fa489b2abe40d479cc5ef47d0b6a6960613148a
networking:
  apiServerAddress: 1.2.3.4
  apiServerPort: 8443
`

type auditConfig struct{}

func (s auditConfig) GetName() string {
	return "auditConfig"
}

func (s auditConfig) GetKubeAPIPort() int {
	return 8443
}

func (s auditConfig) GetKubeAPIIP() string {
	return "1.2.3.4"
}

func (s auditConfig) GetGPUs() int {
	return 0
}

func (s auditConfig) GetExportPorts() []machine.ExportPortPair {
	return []machine.ExportPortPair{
		machine.ExportPortPair{
			HostPort:      80,
			ContainerPort: 8081,
		},
	}
}

func (s auditConfig) AuditEnabled() bool {
	return true
}

func (s auditConfig) AuditFileAbsolutePath() string {
	return "foo.bar.yaml"
}

func (s auditConfig) GetWorkerIDs() []int {
	return []int{}
}

func (s auditConfig) GetNodeLabels() []machine.NodeLabel {
	return []machine.NodeLabel{}
}

func TestKindTemplateWithAudit(t *testing.T) {
	kt := NewKindTemplate()
	assert.NoError(t, kt.Populate(auditConfig{}))
	buf := &bytes.Buffer{}
	assert.NoError(t, kt.Execute(buf))
	assert.EqualValues(t, goldWithAudit, buf.String())
}

var goldWithAudit = `
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: auditConfig
nodes:
- role: control-plane
  kubeadmConfigPatches:
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
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  image: kindest/node:v1.23.12@sha256:9402cf1330bbd3a0d097d2033fa489b2abe40d479cc5ef47d0b6a6960613148a
  gpus: false
  extraPortMappings:
  - containerPort: 8081
    hostPort: 80
    protocol: TCP
  # mount the local file on the control plane
  extraMounts:
  - hostPath: foo.bar.yaml
    containerPath: /etc/kubernetes/policies/audit-policy.yaml
    readOnly: true
networking:
  apiServerAddress: 1.2.3.4
  apiServerPort: 8443
`
