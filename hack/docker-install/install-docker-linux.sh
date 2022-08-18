#!/usr/bin/env bash

# run as root
if (( $EUID != 0 )); then
   echo "this script should be running as root identity"
   exit
fi

apt-get update
apt-get install -y apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository -y \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io

echo "==============================="
echo "installation completed, please add your user into docker group, something like"
echo "****"
echo "usermod -aG docker ubuntu"
echo "****"
echo "for user ubuntu"
echo "And try to logout/login again, and see if `docker ps` works"
