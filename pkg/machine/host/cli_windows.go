//go:build windows
// +build windows

package host

var urlBinaryRes = urlBinary{
	Os:      "windows",
	Kind:    "https://kind.sigs.k8s.io/dl/v0.11.1/kind-windows-amd64",
	Kubectl: "https://storage.googleapis.com/kubernetes-release/release/v1.20.7/bin/windows/amd64/kubectl.exe",
}
