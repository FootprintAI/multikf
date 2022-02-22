//go:build darwin
// +build darwin

package host

var urlBinaryRes = urlBinary{
	Os:      "darwin",
	Kind:    "https://kind.sigs.k8s.io/dl/v0.11.1/kind-darwin-amd64",
	Kubectl: "https://storage.googleapis.com/kubernetes-release/release/v1.20.7/bin/darwin/amd64/kubectl",
}
