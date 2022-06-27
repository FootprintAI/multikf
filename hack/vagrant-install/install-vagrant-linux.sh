#!/usr/bin/env bash

# run as root
if (( $EUID != 0 )); then
   echo "this script should be running as root identity"
   exit
fi

apt-get update
apt-get install virtualbox -y

curl -O https://releases.hashicorp.com/vagrant/2.2.9/vagrant_2.2.9_x86_64.deb
apt-get install ./vagrant_2.2.9_x86_64.deb

echo "======================="
echo "installation completed, use `vagrant --version` to check vagrant installation"
