#!/usr/bin/env bash

# run as root
if (( $EUID != 0 )); then
   echo "this script should be running as root identity"
   exit
fi

VERSION_STRING=5:23.0.4-1~ubuntu.20.04~focal

apt-get update
apt-get install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository -y \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
apt-get update
apt-get install -y docker-ce=$VERSION_STRING docker-ce-cli=$VERSION_STRING containerd.io

echo "==============================="
echo "installation completed, please add your user into docker group, something like"
echo "****"
echo "usermod -aG docker ubuntu"
echo "****"
echo "for user ubuntu"
echo "And try to logout/login again, and see if `docker ps` works"

# noted(hsiny): use low mtu 1374 to enforce this won't hit any router limits (default 1424 - 50 overhead)
# see issue https://github.com/harvester/harvester/issues/3822
tee /etc/docker/daemon.json <<EOF
{
    "exec-opts": ["native.cgroupdriver=systemd"],
    "mtu": 1374
}
EOF

# restart dockerd
systemctl daemon-reload
systemctl restart docker


