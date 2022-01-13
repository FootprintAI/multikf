#!/usr/bin/env bash

# init vagrant vm and install docker/kind/k8s
# vagrant up test0

# copy .kube/config to hst
# ./vagrant-scp.sh test0:/home/vagrant/kubeconfig ./testconfig
# kubectl get node --kubeconfig=testconfig --context kind-cluster1

# teardown test0
# vagrant destroy test0 -f
