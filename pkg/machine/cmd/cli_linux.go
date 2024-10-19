//go:build linux
// +build linux

package cmd

import (
	"fmt"

	"github.com/footprintai/multikf/pkg/k8s"
)

var OSUrlBinaryRes = BinaryResource{
	Os:   "linux",
	Kind: "https://github.com/FootprintAI/kind/releases/download/v0.24.0-gpu/kind-linux",

	Kubectl: func(v k8s.KindK8sVersion) string {
		return fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/kubectl", v.Version())
	},
}

var OSLocalBinaryRes = BinaryResource{
	Os:   "linux",
	Kind: "kind",
	Kubectl: func(_ k8s.KindK8sVersion) string {
		return "kubectl"
	},
}
