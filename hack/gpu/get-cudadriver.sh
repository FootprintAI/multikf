#!/usr/bin/env bash

# run as root
if (( $EUID != 0 )); then
   echo "this script should be running as root identity"
   exit
fi

OS=ubuntu2004
cudnn_version=8.6.0.*
cuda_version=cuda11.8

echo " this script runs on $OS, for other version please check https://developer.nvidia.com/cuda-downloads"

# purge previous installation
# apt-get remove --purge '^nvidia-.*'
# apt-get remove --purge '^libnvidia-.*'
# apt-get remove --purge '^cuda-.*'

apt-get update
apt-get install -y wget

wget https://developer.download.nvidia.com/compute/cuda/repos/$OS/x86_64/cuda-$OS.pin \
    && mv cuda-$OS.pin /etc/apt/preferences.d/cuda-repository-pin-600
apt-key adv --fetch-keys https://developer.download.nvidia.com/compute/cuda/repos/$OS/x86_64/3bf863cc.pub
add-apt-repository "deb https://developer.download.nvidia.com/compute/cuda/repos/$OS/x86_64/ /"

apt-get update

# install older driver
# apt-get install -y nvidia-driver-450 for k80

# or use apt-get install -y nvidia-driver-515 to install previous driver version to avoid conflict in cuda11.8
apt-get install -y nvidia-driver-530 libcudnn8=$cudnn_version-1+$cuda_version libcudnn8-dev=$cudnn_version-1+$cuda_version
apt-mark hold nvidia-driver-530
