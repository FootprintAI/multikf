#!/usr/bin/env bash

# run as root
if (( $EUID != 0 )); then
   echo "this script should be running as root identity"
   exit
fi

VERSION_STRING=5:28.1.1-1~ubuntu.22.04~focal
apt-get update
apt-get install ca-certificates curl
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get install docker-ce=$VERSION_STRING docker-ce-cli=$VERSION_STRING containerd.io docker-buildx-plugin docker-compose-plugin

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


