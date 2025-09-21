#!/usr/bin/env bash

# Usage: ./get-cudadriver-2204.sh [CUDA_VERSION]
# Example: ./get-cudadriver-2204.sh 11-8
# Example: ./get-cudadriver-2204.sh 12-0
# If no version specified, installs latest available

# run as root
if (( $EUID != 0 )); then
   echo "this script should be running as root identity"
   exit 1
fi

# Parse CUDA version parameter
CUDA_VERSION=${1:-"latest"}

# Validate CUDA version format (should be like 11-8, 12-0, etc.)
if [[ "$CUDA_VERSION" != "latest" ]] && [[ ! "$CUDA_VERSION" =~ ^[0-9]+-[0-9]+$ ]]; then
    echo "Error: Invalid CUDA version format. Use format like '11-8' or '12-0'"
    echo "Available versions: 11-8, 12-0, 12-1, 12-2, 12-3, 12-4, 12-5, 12-6"
    echo "Usage: $0 [CUDA_VERSION]"
    echo "Example: $0 11-8"
    exit 1
fi

if [[ "$CUDA_VERSION" != "latest" ]]; then
    echo "Installing CUDA version: $CUDA_VERSION"
else
    echo "Installing latest CUDA version"
fi

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

apt-get update
apt-get install -y wget

wget https://developer.download.nvidia.com/compute/cuda/repos/$OS/x86_64/cuda-$OS.pin \
    && mv cuda-$OS.pin /etc/apt/preferences.d/cuda-repository-pin-600
apt-key adv --fetch-keys https://developer.download.nvidia.com/compute/cuda/repos/$OS/x86_64/3bf863cc.pub
add-apt-repository "deb https://developer.download.nvidia.com/compute/cuda/repos/$OS/x86_64/ /" -y

apt-get update

# install older driver
# apt-get install -y nvidia-driver-450 for k80

# or use apt-get install -y nvidia-driver-515 to install previous driver version to avoid conflict in cuda11.8
apt-get install -y nvidia-driver-530
apt-mark hold nvidia-driver-530


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

