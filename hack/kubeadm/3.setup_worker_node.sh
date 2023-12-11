#!/usr/bin/env bash

## or run
#
# kubeadm token create --print-join-command
#
#

kubeadm join $MASTERIP:8443 --token $KUBEADM_TOKEN \
    --discovery-token-unsafe-skip-ca-verification

