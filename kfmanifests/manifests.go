package kfmanifests

import (
	_ "embed"
)

//go:embed kubeflow-manifest-v1.4.1-template.yaml
var KF14TemplateString string

//go:embed kubeflow-manifest-v1.5.1-template.yaml
var KF15TemplateString string

//go:embed kubeflow-manifest-v1.6.0-template.yaml
var KF16TemplateString string

// NOTE(hsiny): all customized variables used in KF14TemplateString are tagged with [[ xxx ]], whereas default golang template is with {{ yy }}

// run kustomize build base/kf14 > kubeflow-manifest-v1.4.1-template.yaml
// run kustomize build base/kf15 > kubeflow-manifest-v1.5.1-template.yaml
// and replace the default password with template manually
