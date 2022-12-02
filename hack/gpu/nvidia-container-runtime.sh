#!/usr/bin/env bash

# run as root
if (( $EUID != 0 )); then
   echo "this script should be running as root identity"
   exit
fi

distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | tee /etc/apt/sources.list.d/nvidia-docker.list

apt-get update && apt-get install -y nvidia-container-toolkit nvidia-container-runtime

# if you were using containerd, please check here: https://github.com/NVIDIA/k8s-device-plugin#configure-containerd
# append /etc/docker/daemon.json with the following config
tee /etc/docker/daemon.json <<EOF
{
    "default-runtime": "nvidia",
    "runtimes": {
        "nvidia": {
            "path": "/usr/bin/nvidia-container-runtime",
            "runtimeArgs": []
        }
    }
}
EOF

# restart dockerd
systemctl daemon-reload
systemctl restart docker

# test docker run with nvidia container
docker run --gpus all nvidia/cuda:10.0-base nvidia-smi > /dev/null
if [ $? -eq 0 ]; then
    echo "dockerd with nvidia is ready!"
else
    echo "dockerd with nvidia failed!"
    exit -1
fi
