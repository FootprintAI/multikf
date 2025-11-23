#!/usr/bin/env bash

# Simple CUDA installation script for Ubuntu 22.04
# Usage: ./get-cudadriver-2204.sh [CUDA_VERSION]
# Example: ./get-cudadriver-2204.sh 12-0
# Example: ./get-cudadriver-2204.sh (installs latest)

set -e

# Check root
if (( $EUID != 0 )); then
   echo "Error: This script must be run as root"
   exit 1
fi

# Parse version
CUDA_VERSION=${1:-""}
OS=ubuntu2204

if [[ -n "$CUDA_VERSION" ]] && [[ ! "$CUDA_VERSION" =~ ^[0-9]+-[0-9]+$ ]]; then
    echo "Error: Invalid CUDA version format. Use format like '12-0' or '12-6'"
    echo "Usage: $0 [CUDA_VERSION]"
    exit 1
fi

echo "=========================================="
echo "CUDA Installation for Ubuntu 22.04"
echo "=========================================="
echo ""

# Update and install prerequisites
echo "[1/5] Installing prerequisites..."
apt-get update -qq
apt-get install -y wget gnupg > /dev/null

# Setup CUDA repository
echo "[2/5] Setting up CUDA repository..."
wget -q https://developer.download.nvidia.com/compute/cuda/repos/$OS/x86_64/cuda-$OS.pin
mv cuda-$OS.pin /etc/apt/preferences.d/cuda-repository-pin-600

wget -qO- https://developer.download.nvidia.com/compute/cuda/repos/$OS/x86_64/3bf863cc.pub | \
    gpg --dearmor -o /usr/share/keyrings/cuda-archive-keyring.gpg

echo "deb [signed-by=/usr/share/keyrings/cuda-archive-keyring.gpg] https://developer.download.nvidia.com/compute/cuda/repos/$OS/x86_64/ /" | \
    tee /etc/apt/sources.list.d/cuda-$OS.list > /dev/null

apt-get update -qq

# Install CUDA (includes drivers)
echo "[3/5] Installing CUDA toolkit and drivers..."
if [[ -n "$CUDA_VERSION" ]]; then
    echo "Installing CUDA $CUDA_VERSION..."
    apt-get install -y cuda-$CUDA_VERSION
else
    echo "Installing latest CUDA..."
    apt-get install -y cuda
fi

# Setup environment
echo "[4/5] Setting up environment..."
cat > /etc/profile.d/cuda.sh <<'EOF'
export PATH=/usr/local/cuda/bin${PATH:+:${PATH}}
export LD_LIBRARY_PATH=/usr/local/cuda/lib64${LD_LIBRARY_PATH:+:${LD_LIBRARY_PATH}}
export CUDA_HOME=/usr/local/cuda
EOF
chmod +x /etc/profile.d/cuda.sh

echo "[5/5] Installation complete!"
echo ""
echo "=========================================="
echo "Next steps:"
echo "=========================================="
echo "1. Reboot your system: sudo reboot"
echo "2. After reboot, verify:"
echo "   nvidia-smi"
echo "   nvcc --version"
echo ""
