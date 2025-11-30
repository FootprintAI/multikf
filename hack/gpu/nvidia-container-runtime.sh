#!/usr/bin/env bash

# run as root
if (( $EUID != 0 )); then
   echo "this script should be running as root identity"
   exit 1
fi

echo "Installing NVIDIA Container Toolkit..."

# Install prerequisites
apt-get update && apt-get install -y --no-install-recommends curl gnupg2

# Configure the production repository with GPG key (modern approach)
echo "Configuring NVIDIA Container Toolkit repository..."
curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg \
  && curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | \
    sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | \
    tee /etc/apt/sources.list.d/nvidia-container-toolkit.list

# Update package list and install toolkit
echo "Installing nvidia-container-toolkit..."
apt-get update
apt-get install -y nvidia-container-toolkit

# Configure Docker to use NVIDIA runtime (replaces manual daemon.json editing)
echo "Configuring Docker runtime..."
nvidia-ctk runtime configure --runtime=docker

# Optionally preserve custom MTU and cgroup settings if needed
# You may need to manually merge these settings into /etc/docker/daemon.json:
# {
#     "mtu": 1374,
#     "exec-opts": ["native.cgroupdriver=systemd"]
# }

# if you were using containerd, please check here: https://github.com/NVIDIA/k8s-device-plugin#configure-containerd

# restart dockerd
systemctl daemon-reload
systemctl restart docker

# test docker run with nvidia container
docker run --gpus all nvidia/cuda:12.2.0-base-ubuntu20.04 nvidia-smi > /dev/null
if [ $? -eq 0 ]; then
    echo "dockerd with nvidia is ready!"
else
    echo "dockerd with nvidia failed!"
    exit -1
fi
