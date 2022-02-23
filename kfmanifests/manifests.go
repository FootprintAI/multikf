package kfmanifests

import (
	"embed"
)

//go:embed kubeflow-manifest-v1.4.1.yaml
var FS embed.FS
