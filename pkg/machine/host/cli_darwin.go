//go:build darwin
// +build darwin

package host

var urlBinaryRes = binaryResource{
	Os:      "darwin",
	Kind:    "https://github.com/FootprintAI/kind/raw/gpu/bin/kind-darwin",
	Kubectl: "https://storage.googleapis.com/kubernetes-release/release/v1.21.2/bin/darwin/amd64/kubectl",
}

var localBinaryRes = binaryResource{
	Os:      "darwin",
	Kind:    "kind",
	Kubectl: "kubectl",
}
