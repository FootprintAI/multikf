#!/usr/bin/env bash

# allow iptable see bridged traffic
cat <<EOF | tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF

sysctl --system

# enable ipv4 forwarding
bash -c 'echo 1 > /proc/sys/net/ipv4/ip_forward'


# disable swap memory
# or comment the line on /etc/fstab
#
# vim /etc/fstab
# #/swap.img       none    swap    sw      0       0
# 
# then reboot the machine
# or run the following command to disable temporary
swapoff -a


