#!/usr/bin/env bash

kubeadm join $MASTERIP:8443 --token $KUBEADM_TOKEN \
    --discovery-token-unsafe-skip-ca-verification

