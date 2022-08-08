#!/usr/bin/env bash

kustomize build base/kf14 > kubeflow-manifest-v1.4.1-template.yaml
kustomize build base/kf15 > kubeflow-manifest-v1.5.1-template.yaml
