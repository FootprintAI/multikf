#!/usr/bin/env bash

kubeadm init --apiserver-bind-port=8443 \
--apiserver-advertise-address=$PUBLICIP \
--kubernetes-version=$KUBECTL_VERSION \
--pod-network-cidr=10.244.0.0/16 \
--service-cidr=10.96.0.0/12 \
--token=$KUBEADM_TOKEN \
--ignore-preflight-errors=NumCPU \
--apiserver-cert-extra-sans="$LOCALHOST,$PUBLICIP"

# move to .kube, so kubectl can work
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# add network plugins
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml


# untaint master node
kubectl taint nodes <node-name> node-role.kubernetes.io/master-
kubectl taint nodes <node-name> node-role.kubernetes.io/control-plane-

## taint master
# kubectl taint nodes <node-name> node-role.kubernetes.io/master:NoSchedule
# kubectl taint nodes <node-name> node-role.kubernetes.io/control-plane:NoSchedule



##### cilium

kubeadm init --apiserver-bind-port=8443 \
  --apiserver-advertise-address=$PUBLICIP \
  --kubernetes-version=$KUBECTL_VERSION \
  --pod-network-cidr=10.0.0.0/16 \  # Recommended default for Cilium
  --service-cidr=10.96.0.0/12 \
  --token=$KUBEADM_TOKEN \
  --ignore-preflight-errors=NumCPU \
  --apiserver-cert-extra-sans="$LOCALHOST,$PUBLICIP"

# Setup kubectl config
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config


# Install Cilium (via CLI installer)
curl -L --remote-name https://github.com/cilium/cilium-cli/releases/latest/download/cilium-linux-amd64.tar.gz
tar xzvf cilium-linux-amd64.tar.gz
sudo mv cilium /usr/local/bin/

# Install Cilium into the cluster
cilium install

# Optionally verify installation
cilium status --wait
