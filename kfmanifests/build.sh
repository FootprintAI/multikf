#!/usr/bin/env bash

kustomize build base/kf14 > manifests/kubeflow-manifest-v1.4.1-template.yaml
kustomize build base/kf15 > manifests/kubeflow-manifest-v1.5.1-template.yaml
kustomize build base/kf16 > manifests/kubeflow-manifest-v1.6.1-template.yaml
kustomize build base/kf16-lite > manifests/kubeflow-manifest-v1.6.1-lite-template.yaml
