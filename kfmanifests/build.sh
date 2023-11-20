#!/usr/bin/env bash

#kustomize build base/kf16 > manifests/kubeflow-manifest-v1.6.1-template.yaml
#kustomize build base/kf17 > manifests/kubeflow-manifest-v1.7.0-template.yaml
kustomize build base/kf18 > manifests/kubeflow-manifest-v1.8.0-template.yaml
