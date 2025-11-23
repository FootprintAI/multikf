#!/usr/bin/env bash

# NVIDIA Container Runtime installation script for Docker
# Usage: ./nvidia-container-runtime.sh
# Configures Docker to use NVIDIA GPU support with nvidia-container-runtime

set -e

# Check root
if (( $EUID != 0 )); then
   echo "Error: This script must be run as root"
   exit 1
fi

echo "=========================================="
echo "NVIDIA Container Runtime Setup"
echo "=========================================="
echo ""

# Detect distribution
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)

# Setup NVIDIA Docker repository
echo "[1/5] Setting up NVIDIA Docker repository..."
wget -qO- https://nvidia.github.io/nvidia-docker/gpgkey | \
    gpg --dearmor -o /usr/share/keyrings/nvidia-docker-archive-keyring.gpg

wget -qO- https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | \
    sed "s#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-docker-archive-keyring.gpg] https://#g" | \
    tee /etc/apt/sources.list.d/nvidia-docker.list > /dev/null

# Install NVIDIA container packages
echo "[2/5] Installing NVIDIA container toolkit and runtime..."
apt-get update -qq
apt-get install -y nvidia-container-toolkit nvidia-container-runtime > /dev/null

# Configure Docker daemon
echo "[3/5] Configuring Docker daemon..."
# Note: if you were using containerd, please check: https://github.com/NVIDIA/k8s-device-plugin#configure-containerd
tee /etc/docker/daemon.json > /dev/null <<EOF
{
    "mtu": 1374,
    "exec-opts": ["native.cgroupdriver=systemd"],
    "default-runtime": "nvidia",
    "runtimes": {
        "nvidia": {
            "path": "/usr/bin/nvidia-container-runtime",
            "runtimeArgs": []
        }
    }
}
EOF

# Restart Docker service
echo "[4/5] Restarting Docker daemon..."
systemctl daemon-reload
systemctl restart docker

# Verify installation
echo "[5/5] Verifying NVIDIA container runtime..."
if docker run --rm --gpus all nvidia/cuda:12.2.0-base-ubuntu20.04 nvidia-smi > /dev/null 2>&1; then
    echo ""
    echo "=========================================="
    echo "Success! NVIDIA Docker runtime is ready"
    echo "=========================================="
    echo ""
    echo "You can now run GPU-enabled containers with:"
    echo "  docker run --gpus all <image> <command>"
    echo ""
else
    echo ""
    echo "=========================================="
    echo "Error: NVIDIA container runtime verification failed"
    echo "=========================================="
    echo ""
    echo "Troubleshooting steps:"
    echo "1. Verify NVIDIA drivers are installed: nvidia-smi"
    echo "2. Check Docker logs: journalctl -u docker"
    echo "3. Check /etc/docker/daemon.json configuration"
    echo ""
    exit 1
fi
