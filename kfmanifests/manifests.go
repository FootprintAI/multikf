package kfmanifests

import (
	_ "embed"
)

//go:embed kubeflow-manifest-v1.4.1.yaml
var KF14 string

//go:embed kubeflow-manifest-v1.4.1-template.yaml
var KF14TemplateString string

// NOTE(hsiny): all customized variables used in KF14TemplateString are tagged with [[ xxx ]], whereas default golang template is with {{ yy }}
