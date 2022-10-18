#!/usr/bin/env bash

# run as root
if (( $EUID != 0 )); then
   echo "this script should be running as root identity"
   exit
fi

firewall-cmd --add-port=10250/tcp --permanent

# kubedns
firewall-cmd --add-masquerade --permanent

# flannel
firewall-cmd --add-port=8285/udp --permanent
firewall-cmd --add-port=8472/udp --permanent

# longhorn
firewall-cmd --add-port=9443/tcp --permanent
firewall-cmd --add-port=9500/tcp --permanent

# kube api
firewall-cmd --add-port=8443/tcp --permanent


# reloading
firewall-cmd --reload

