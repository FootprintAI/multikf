#!/usr/bin/env bash

# install dockerd
apt-get update
apt-get install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io

# install kind
curl -Lo ./kind https://github.com/FootprintAI/kind/raw/gpu/bin/kind-linux && \
    chmod +x ./kind && \
    mv ./kind /usr/local/bin/kind
