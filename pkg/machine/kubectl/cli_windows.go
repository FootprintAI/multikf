//go:build windows
// +build windows

package kubectl

var urlBinaryRes = binaryResource{
	Os:      "windows",
	Kind:    "https://github.com/FootprintAI/kind/releases/download/v0.16.0-gpu-master/kind-windows",
	Kubectl: "https://storage.googleapis.com/kubernetes-release/release/v1.21.2/bin/windows/amd64/kubectl.exe",
}

var localBinaryRes = binaryResource{
	Os:      "windows",
	Kind:    "kind.exe",
	Kubectl: "kubectl.exe",
}
