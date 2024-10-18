//go:build windows
// +build windows

package cmd

import (
	"fmt"

	"github.com/footprintai/multikf/pkg/k8s"
)

var OSUrlBinaryRes = BinaryResource{
	Os:   "windows",
	Kind: "https://github.com/FootprintAI/kind/releases/download/v0.24.0-gpu/kind-windows",
	Kubectl: func(v k8s.KindK8sVersion) string {
		return fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/windows/amd64/kubectl.exe", v.Version())
	},
}

var OSLocalBinaryRes = BinaryResource{
	Os:   "windows",
	Kind: "kind.exe",
	Kubectl: func(_ k8s.KindK8sVersion) string {
		return "kubectl.exe"
	},
}
