#!/usr/bin/env bash

# Usage: ./remove-cudadriver-2204.sh [OPTIONS]
# Options:
#   --keep-drivers    Remove CUDA but keep NVIDIA drivers
#   --full            Remove everything (default)
#
# This script removes CUDA toolkit and optionally NVIDIA drivers installed
# by get-cudadriver-2204.sh

# Exit on error, undefined variables, and pipe failures
set -euo pipefail

# run as root
if (( $EUID != 0 )); then
   echo "this script should be running as root identity"
   exit 1
fi

# Parse options
KEEP_DRIVERS=false
if [[ "${1:-}" == "--keep-drivers" ]]; then
    KEEP_DRIVERS=true
    echo "Mode: Removing CUDA toolkit only, keeping NVIDIA drivers"
elif [[ "${1:-}" == "--full" ]] || [[ -z "${1:-}" ]]; then
    KEEP_DRIVERS=false
    echo "Mode: Full removal (CUDA + NVIDIA drivers)"
else
    echo "Error: Unknown option '$1'"
    echo "Usage: $0 [--keep-drivers|--full]"
    echo "  --keep-drivers    Remove CUDA but keep NVIDIA drivers"
    echo "  --full            Remove everything (default)"
    exit 1
fi

echo ""
echo "======================================"
echo "CUDA/NVIDIA Removal Script"
echo "======================================"
echo ""

# Stop running services that might hold GPU
echo "Stopping NVIDIA services..."
systemctl stop nvidia-persistenced 2>/dev/null || true
systemctl stop nvidia-fabricmanager 2>/dev/null || true

# Check if any processes are using the GPU
if command -v nvidia-smi &> /dev/null; then
    echo ""
    echo "Checking for processes using GPU..."
    if nvidia-smi --query-compute-apps=pid,process_name --format=csv,noheader 2>/dev/null | grep -q .; then
        echo "Warning: The following processes are using the GPU:"
        nvidia-smi --query-compute-apps=pid,process_name --format=csv,noheader
        echo ""
        read -p "Do you want to continue? These processes may be terminated. (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Aborted."
            exit 1
        fi
    fi
fi

# Remove CUDA packages
echo ""
echo "Removing CUDA packages..."
apt-get remove --purge -y \
    'cuda-*' \
    'libcublas*' \
    'libcudnn*' \
    'libnccl*' \
    'nsight-*' 2>/dev/null || true

# Remove NVIDIA drivers if not keeping them
if [[ "$KEEP_DRIVERS" == false ]]; then
    echo ""
    echo "Removing NVIDIA drivers and libraries..."

    # Remove DKMS modules first
    echo "Removing NVIDIA DKMS modules..."
    for dkms_module in $(dkms status | grep nvidia | cut -d',' -f1 | cut -d':' -f1 | sort -u); do
        for version in $(dkms status "$dkms_module" | cut -d',' -f2 | cut -d':' -f1 | tr -d ' '); do
            echo "  Removing DKMS module: $dkms_module/$version"
            dkms remove "$dkms_module/$version" --all 2>/dev/null || true
        done
    done

    # Remove all NVIDIA packages
    apt-get remove --purge -y \
        'nvidia-*' \
        'libnvidia-*' 2>/dev/null || true
else
    echo ""
    echo "Keeping NVIDIA drivers as requested..."
fi

# Clean up apt
echo ""
echo "Cleaning up package manager..."
apt-get autoremove -y
apt-get autoclean

# Remove leftover directories
echo ""
echo "Removing leftover directories and files..."

# Remove CUDA directories
if [ -d /usr/local/cuda ]; then
    echo "  Removing /usr/local/cuda*"
    rm -rf /usr/local/cuda*
fi

# Remove CUDA repository configs
if [ -d /var/lib/cuda-repo* ]; then
    echo "  Removing /var/lib/cuda-repo*"
    rm -rf /var/lib/cuda-repo*
fi

# Remove repository files
if ls /etc/apt/sources.list.d/cuda*.list 1> /dev/null 2>&1; then
    echo "  Removing CUDA repository lists"
    rm -f /etc/apt/sources.list.d/cuda*.list
fi

if ls /etc/apt/sources.list.d/nvidia*.list 1> /dev/null 2>&1; then
    echo "  Removing NVIDIA repository lists"
    rm -f /etc/apt/sources.list.d/nvidia*.list
fi

# Remove keyrings
if ls /usr/share/keyrings/cuda*.gpg 1> /dev/null 2>&1; then
    echo "  Removing CUDA keyrings"
    rm -f /usr/share/keyrings/cuda*.gpg
fi

if ls /usr/share/keyrings/nvidia*.gpg 1> /dev/null 2>&1; then
    echo "  Removing NVIDIA keyrings"
    rm -f /usr/share/keyrings/nvidia*.gpg
fi

# Remove preferences
if [ -f /etc/apt/preferences.d/cuda-repository-pin-600 ]; then
    echo "  Removing CUDA repository pin"
    rm -f /etc/apt/preferences.d/cuda-repository-pin-600
fi

# Remove environment configuration
if [ -f /etc/profile.d/cuda.sh ]; then
    echo "  Removing CUDA environment configuration (/etc/profile.d/cuda.sh)"
    rm -f /etc/profile.d/cuda.sh
fi

# Clean up user configurations (optional - commented out for safety)
# echo "  Removing CUDA user configurations"
# rm -rf ~/.nv
# rm -rf ~/.cuda

# Update package cache
echo ""
echo "Updating package lists..."
apt-get update

echo ""
echo "======================================"
echo "Removal completed!"
echo "======================================"

if [[ "$KEEP_DRIVERS" == false ]]; then
    echo ""
    echo "CUDA toolkit and NVIDIA drivers have been removed."
    echo ""
    echo "IMPORTANT: You should reboot the system to complete the removal."
    echo "The nouveau driver (open-source) will be used after reboot."
    echo ""
    echo "If you plan to reinstall NVIDIA drivers, you may want to:"
    echo "  1. Reboot first to ensure clean state"
    echo "  2. Run the installation script again"
else
    echo ""
    echo "CUDA toolkit has been removed, NVIDIA drivers are still installed."
    echo ""
    echo "You can verify with:"
    echo "  - nvidia-smi (should still work)"
    echo "  - nvcc --version (should not be found)"
    echo ""
fi
