#!/usr/bin/env bash

kind create cluster --image=kindest/node:v1.20.7 --config /tmp/kind-config.yaml

sudo cp -r /root/.kube /home/vagrant/.kube
sudo chown -R vagrant:vagrant /home/vagrant/.kube
