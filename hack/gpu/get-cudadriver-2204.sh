#!/usr/bin/env bash

# Usage: ./get-cudadriver-2204.sh [CUDA_VERSION] [NVIDIA_DRIVER_VERSION]
# Example: ./get-cudadriver-2204.sh 11-8 530
# Example: ./get-cudadriver-2204.sh 12-0 535
# If no version specified, installs latest CUDA with driver 535

# Exit on error, undefined variables, and pipe failures
set -euo pipefail

# run as root
if (( $EUID != 0 )); then
   echo "this script should be running as root identity"
   exit 1
fi

# Parse CUDA version parameter
CUDA_VERSION=${1:-"latest"}
NVIDIA_DRIVER_VERSION=${2:-"535"}

# Validate CUDA version format (should be like 11-8, 12-0, etc.)
if [[ "$CUDA_VERSION" != "latest" ]] && [[ ! "$CUDA_VERSION" =~ ^[0-9]+-[0-9]+$ ]]; then
    echo "Error: Invalid CUDA version format. Use format like '11-8' or '12-0'"
    echo "Available versions: 11-8, 12-0, 12-1, 12-2, 12-3, 12-4, 12-5, 12-6"
    echo "Usage: $0 [CUDA_VERSION] [NVIDIA_DRIVER_VERSION]"
    echo "Example: $0 11-8 530"
    echo "Example: $0 12-0 535"
    echo "Default: latest CUDA with driver 535"
    exit 1
fi

if [[ "$CUDA_VERSION" != "latest" ]]; then
    echo "Installing CUDA version: $CUDA_VERSION"
else
    echo "Installing latest CUDA version"
fi
echo "Installing NVIDIA driver version: $NVIDIA_DRIVER_VERSION"

OS=ubuntu2204

echo " this script runs on $OS, for other version please check https://developer.nvidia.com/cuda-downloads"

# purge previous installation
# apt-get remove --purge '^nvidia-.*'
# apt-get remove --purge '^libnvidia-.*'
# apt-get remove --purge '^cuda-.*'

# fix DKMS driver
# $ dkms status
#  nvidia/535.154.05: added
# $ dkms build nvidia/535.154.05
# $ dkms install nvidia/535.154.05 --force
# then reboot

echo "Updating package lists..."
apt-get update

echo "Installing prerequisites..."
apt-get install -y wget gnupg

# Setup CUDA repository with proper keyring (apt-key is deprecated in Ubuntu 22.04)
echo "Downloading CUDA repository pin..."
wget https://developer.download.nvidia.com/compute/cuda/repos/$OS/x86_64/cuda-$OS.pin \
    && mv cuda-$OS.pin /etc/apt/preferences.d/cuda-repository-pin-600

# Add CUDA GPG key using keyring method
echo "Adding CUDA GPG key..."
wget -O- https://developer.download.nvidia.com/compute/cuda/repos/$OS/x86_64/3bf863cc.pub | \
    gpg --dearmor -o /usr/share/keyrings/cuda-archive-keyring.gpg

if [ ! -f /usr/share/keyrings/cuda-archive-keyring.gpg ]; then
    echo "Error: Failed to create CUDA keyring"
    exit 1
fi

# Add repository with signed-by keyring
echo "Adding CUDA repository..."
echo "deb [signed-by=/usr/share/keyrings/cuda-archive-keyring.gpg] https://developer.download.nvidia.com/compute/cuda/repos/$OS/x86_64/ /" | \
    tee /etc/apt/sources.list.d/cuda-$OS.list > /dev/null

echo "Updating package lists with CUDA repository..."
apt-get update

# Pre-installation cleanup
echo ""
echo "Checking for conflicting packages..."

# Check for held packages
HELD_PACKAGES=$(apt-mark showhold | grep -E '^(nvidia|cuda)' || true)
if [ -n "$HELD_PACKAGES" ]; then
    echo "Found held packages that may cause conflicts:"
    echo "$HELD_PACKAGES"
    echo ""
    read -p "Do you want to unhold these packages? (recommended) (Y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        echo "Unholding packages..."
        echo "$HELD_PACKAGES" | xargs -r apt-mark unhold
    fi
fi

# Check for broken packages
echo "Checking for broken packages..."
if dpkg -l | grep -E '^(iU|iF)' | grep -qE '(nvidia|cuda)'; then
    echo "Found broken NVIDIA/CUDA packages. Attempting to fix..."
    apt-get install -f -y || true
    dpkg --configure -a || true
fi

# Remove any conflicting partial installations
echo "Cleaning up any partial installations..."
apt-get autoremove -y
apt-get autoclean

echo ""
# install NVIDIA driver
# Older drivers: nvidia-driver-450 for k80, nvidia-driver-515 for cuda11.8 compatibility
# Modern drivers: nvidia-driver-535, nvidia-driver-550, nvidia-driver-560
echo "Installing nvidia-driver-$NVIDIA_DRIVER_VERSION..."

DRIVER_INSTALLED=false

# Method 1: Try Ubuntu repository driver
echo "Attempting Method 1: Ubuntu repository (nvidia-driver-$NVIDIA_DRIVER_VERSION)..."
if apt-get install -y nvidia-driver-$NVIDIA_DRIVER_VERSION 2>/dev/null; then
    echo "✓ Successfully installed nvidia-driver-$NVIDIA_DRIVER_VERSION from Ubuntu repository"
    apt-mark hold nvidia-driver-$NVIDIA_DRIVER_VERSION
    DRIVER_INSTALLED=true
else
    echo "✗ Failed to install from Ubuntu repository (dependency conflicts)"
fi

# Method 2: Try CUDA repository driver with specific version
if [[ "$DRIVER_INSTALLED" == false ]]; then
    echo ""
    echo "Attempting Method 2: CUDA repository (cuda-drivers-$NVIDIA_DRIVER_VERSION)..."
    if apt-get install -y cuda-drivers-$NVIDIA_DRIVER_VERSION 2>/dev/null; then
        echo "✓ Successfully installed cuda-drivers-$NVIDIA_DRIVER_VERSION from CUDA repository"
        DRIVER_INSTALLED=true
    else
        echo "✗ Package cuda-drivers-$NVIDIA_DRIVER_VERSION not available"
    fi
fi

# Method 3: Try generic cuda-drivers metapackage
if [[ "$DRIVER_INSTALLED" == false ]]; then
    echo ""
    echo "Attempting Method 3: Generic cuda-drivers metapackage..."
    if apt-get install -y cuda-drivers 2>/dev/null; then
        echo "✓ Successfully installed cuda-drivers from CUDA repository"
        DRIVER_INSTALLED=true
    else
        echo "✗ Failed to install cuda-drivers metapackage"
    fi
fi

# Method 4: Try alternative stable driver version (535)
if [[ "$DRIVER_INSTALLED" == false ]] && [[ "$NVIDIA_DRIVER_VERSION" != "535" ]]; then
    echo ""
    echo "Attempting Method 4: Fallback to stable driver (nvidia-driver-535)..."
    if apt-get install -y nvidia-driver-535 2>/dev/null; then
        echo "✓ Successfully installed nvidia-driver-535 as fallback"
        apt-mark hold nvidia-driver-535
        DRIVER_INSTALLED=true
        NVIDIA_DRIVER_VERSION="535"
    else
        echo "✗ Failed to install nvidia-driver-535"
    fi
fi

# Final check
if [[ "$DRIVER_INSTALLED" == false ]]; then
    echo ""
    echo "=========================================="
    echo "ERROR: All driver installation methods failed"
    echo "=========================================="
    echo ""
    echo "Running diagnostics..."
    echo ""

    echo "--- Available CUDA driver packages ---"
    apt-cache search cuda-drivers | head -10 || true
    echo ""

    echo "--- Available NVIDIA drivers in Ubuntu repo ---"
    apt-cache search nvidia-driver | grep '^nvidia-driver-[0-9]' | head -10 || true
    echo ""

    echo "--- Currently held packages ---"
    apt-mark showhold | grep -E '(nvidia|cuda)' || echo "None"
    echo ""

    echo "--- Broken packages ---"
    dpkg -l | grep -E '^(iU|iF)' | grep -E '(nvidia|cuda)' || echo "None"
    echo ""

    echo "Please try one of the following manual approaches:"
    echo ""
    echo "1. Clean removal and retry:"
    echo "   ./remove-cudadriver-2204.sh --full"
    echo "   reboot"
    echo "   ./get-cudadriver-2204.sh"
    echo ""
    echo "2. Manually specify a different driver version:"
    echo "   ./get-cudadriver-2204.sh latest 525"
    echo "   ./get-cudadriver-2204.sh latest 530"
    echo ""
    echo "3. Check for system updates:"
    echo "   apt-get update && apt-get upgrade -y"
    echo ""
    exit 1
fi

echo ""
echo "Driver installation successful!"
echo "Installed driver version: $NVIDIA_DRIVER_VERSION"
echo ""


# install CUDA toolkit
if [[ "$CUDA_VERSION" != "latest" ]]; then
    echo "Installing CUDA toolkit version: cuda-$CUDA_VERSION"
    apt-get install -y cuda-$CUDA_VERSION
    apt-mark hold cuda-$CUDA_VERSION
else
    echo "Installing latest CUDA toolkit"
    apt-get install -y cuda
fi

# install cuda related lib
## cublas for cuda12, ref: https://developer.nvidia.com/nvidia-hpc-sdk-releases
## curl https://developer.download.nvidia.com/hpc-sdk/ubuntu/DEB-GPG-KEY-NVIDIA-HPC-SDK | gpg --dearmor -o /usr/share/keyrings/nvidia-hpcsdk-archive-keyring.gpg
## echo 'deb [signed-by=/usr/share/keyrings/nvidia-hpcsdk-archive-keyring.gpg] https://developer.download.nvidia.com/hpc-sdk/ubuntu/amd64 /' | tee /etc/apt/sources.list.d/nvhpc.list
## apt-get update -y
## apt-get install -y nvhpc-24-5

## install cudnn8 for cuda12, ref: https://developer.nvidia.com/rdp/cudnn-archive
## https://developer.nvidia.com/downloads/compute/cudnn/secure/8.9.7/local_installers/12.x/cudnn-local-repo-ubuntu2204-8.9.7.29_1.0-1_amd64.deb/
## you got to login to be able to download it
##
## then run
## dpkg -i <>.deb

# Setup CUDA environment variables
echo "Setting up CUDA environment variables..."
cat > /etc/profile.d/cuda.sh <<'EOF'
# CUDA environment setup
export PATH=/usr/local/cuda/bin${PATH:+:${PATH}}
export LD_LIBRARY_PATH=/usr/local/cuda/lib64${LD_LIBRARY_PATH:+:${LD_LIBRARY_PATH}}
export CUDA_HOME=/usr/local/cuda
EOF

chmod +x /etc/profile.d/cuda.sh

echo ""
echo "======================================"
echo "Installation completed!"
echo "======================================"
echo "NVIDIA driver: $NVIDIA_DRIVER_VERSION"
if [[ "$CUDA_VERSION" != "latest" ]]; then
    echo "CUDA version: $CUDA_VERSION"
else
    echo "CUDA version: latest"
fi
echo ""
echo "Environment variables have been set in /etc/profile.d/cuda.sh"
echo ""
echo "IMPORTANT: You must reboot the system for the NVIDIA driver to take effect."
echo "After reboot, verify the installation with:"
echo "  - nvidia-smi (check driver and GPU status)"
echo "  - nvcc --version (check CUDA toolkit version)"
echo ""
echo "To use CUDA in the current session without rebooting, source the environment:"
echo "  source /etc/profile.d/cuda.sh"
echo ""
