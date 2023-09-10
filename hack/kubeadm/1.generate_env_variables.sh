#!/usr/bin/env bash

MYIP=$(hostname --ip-address)
MASTERIP="edit.master.ip"

cat << EOF | tee ./env.sh
KUBEADM_TOKEN=3chl6w.ymb8xtge15qndyfk
LOCALHOST=127.0.0.1
PUBLICIP="$MYIP"
MASTERIP="$MASTERIP"
KUBECTL_VERSION=v1.24.17
EOF

source ./env.sh
