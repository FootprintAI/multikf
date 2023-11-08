//go:build linux
// +build linux

package kubectl

var urlBinaryRes = binaryResource{
	Os:      "linux",
	Kind:    "https://github.com/FootprintAI/kind/releases/download/v0.20.0-gpu/kind-linux",
	Kubectl: "https://storage.googleapis.com/kubernetes-release/release/v1.25.11/bin/linux/amd64/kubectl",
}

var localBinaryRes = binaryResource{
	Os:      "linux",
	Kind:    "kind",
	Kubectl: "kubectl",
}
