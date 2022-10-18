#!/usr/bin/env bash

# run as root
if (( $EUID != 0 )); then
   echo "this script should be running as root identity"
   exit
fi

ENV OS=ubuntu2204
ENV cudnn_version=8.6.0.*
ENV cuda_version=cuda11.8

echo " this script runs on ${OS}, for other version please check https://developer.nvidia.com/cuda-downloads"

# purge previous installation
# apt-get purge -y nvidia*

apt-get update
apt-get install -y wget

wget https://developer.download.nvidia.com/compute/cuda/repos/${OS}/x86_64/cuda-${OS}.pin \
    && mv cuda-${OS}.pin /etc/apt/preferences.d/cuda-repository-pin-600
apt-key adv --fetch-keys https://developer.download.nvidia.com/compute/cuda/repos/${OS}/x86_64/3bf863cc.pub
add-apt-repository "deb https://developer.download.nvidia.com/compute/cuda/repos/${OS}/x86_64/ /"

apt-get update
# or use apt-get install -y nvidia-driver-515 to install previous driver version to avoid conflict in cuda11.8
apt-get install -y cuda libcudnn8=${cudnn_version}-1+${cuda_version} libcudnn8-dev=${cudnn_version}-1+${cuda_version}
