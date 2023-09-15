#!/usr/bin/env bash

# create containerd configuration
#
containerd config default | tee /etc/containerd/config.toml
systemctl daemon-reload

# add crictl argument
#
cat << EOF | tee /etc/crictl.yaml
runtime-endpoint: unix:///run/containerd/containerd.sock
image-endpoint: unix:///run/containerd/containerd.sock
timeout: 10
debug: false
EOF

# as our installation is from docker, need to restart docker as well
systemctl restart docker containerd
