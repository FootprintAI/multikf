//go:build linux
// +build linux

package docker

var urlBinaryRes = binaryResource{
	Os:      "linux",
	Kind:    "https://github.com/FootprintAI/kind/raw/gpu/bin/kind-linux",
	Kubectl: "https://storage.googleapis.com/kubernetes-release/release/v1.21.2/bin/linux/amd64/kubectl",
}

var localBinaryRes = binaryResource{
	Os:      "linux",
	Kind:    "kind",
	Kubectl: "kubectl",
}
