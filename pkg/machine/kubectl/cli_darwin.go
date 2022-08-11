//go:build darwin
// +build darwin

package kubectl

var urlBinaryRes = binaryResource{
	Os:      "darwin",
	Kind:    "https://github.com/FootprintAI/kind/releases/download/v0.16.0-gpu-master/kind-darwin",
	Kubectl: "https://storage.googleapis.com/kubernetes-release/release/v1.21.2/bin/darwin/amd64/kubectl",
}

var localBinaryRes = binaryResource{
	Os:      "darwin",
	Kind:    "kind",
	Kubectl: "kubectl",
}
