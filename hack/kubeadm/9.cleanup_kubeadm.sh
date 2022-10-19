#!/usr/bin/env bash


# run kubeadm reset to kick off reset process
kubeadm reset -f

# stop kubelet and docker
systemctl stop kubelet
systemctl stop docker

# remove orphant files
rm -rf /var/lib/cni/
rm -rf /var/lib/kubelet/*
rm -rf /etc/cni/

# stop network inferface
ifconfig cni0 down
ifconfig flannel.1 down
ifconfig docker0 down

# remove network interface
ip link delete cni0
ip link delete flannel.1

# restart docker
systemctl restart kubelet
systemctl restart docker
