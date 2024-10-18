//go:build darwin
// +build darwin

package cmd

import (
	"fmt"

	"github.com/footprintai/multikf/pkg/k8s"
)

var OSUrlBinaryRes = BinaryResource{
	Os:   "darwin",
	Kind: "https://github.com/FootprintAI/kind/releases/download/v0.24.0-gpu/kind-darwin",
	Kubectl: func(v k8s.KindK8sVersion) string {
		return fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/darwin/amd64/kubectl", v.Version())
	},
}

var OSLocalBinaryRes = BinaryResource{
	Os:   "darwin",
	Kind: "kind",
	Kubectl: func(_ k8s.KindK8sVersion) string {
		return "kubectl"
	},
}
