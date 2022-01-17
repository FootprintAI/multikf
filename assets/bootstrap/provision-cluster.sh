#!/usr/bin/env bash

kind create cluster --image=kindest/node:v1.20.7 --config /tmp/kind-config.yaml

cp /root/.kube /home/vagrant/.kube
chown -R vagrant:vagrant /home/vagrant/.kube
