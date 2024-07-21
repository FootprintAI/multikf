#!/usr/bin/env bash

# run as root
if (( $EUID != 0 )); then
   echo "this script should be running as root identity"
   exit
fi

OS=ubuntu2204

echo " this script runs on $OS, for other version please check https://developer.nvidia.com/cuda-downloads"

# purge previous installation
# apt-get purge -y nvidia*

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

