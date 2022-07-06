#!/usr/bin/env bash

/usr/local/bin/kind create cluster --config kind-config.yaml

sudo cp -r /root/.kube /home/vagrant/.kube
sudo chown -R vagrant:vagrant /home/vagrant/.kube
