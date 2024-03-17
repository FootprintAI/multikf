#!/usr/bin/env bash

# install kubeadm
# as the package repo has changed, see https://kubernetes.io/docs/tasks/administer-cluster/kubeadm/change-package-repository/
# so we disable the following command, download binaries instead.

# curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
# cat <<EOF | tee /etc/apt/sources.list.d/kubernetes.list
# deb https://apt.kubernetes.io/ kubernetes-xenial main
# EOF
# apt-get update
# apt-get install -y kubelet=1.25.14-00 kubeadm=1.25.14-00 kubectl=1.25.14-00
# apt-mark hold kubeadm kubelet kubectl

echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.25/deb/ /" |  tee /etc/apt/sources.list.d/kubernetes.list
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.25/deb/Release.key |  gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg

## verify pager /etc/apt/sources.list.d/kubernetes.list
## deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.25/deb/ /

apt-get update
apt-get install -y kubelet=1.25.14-1.1 kubeadm=1.25.14-1.1 kubectl=1.25.14-1.1
apt-mark hold kubeadm kubelet kubectl

## version v1.24.17 (this requires additional configuration as apt-get would do that)
#curl -LO https://dl.k8s.io/release/v1.24.17/bin/linux/amd64/kubectl \
#        && chmod +x kubectl \
#        && mv kubectl /usr/bin/
#curl -LO https://dl.k8s.io/release/v1.24.17/bin/linux/amd64/kubeadm \
#        && chmod +x kubeadm \
#        && mv kubeadm /usr/bin/
#curl -LO https://dl.k8s.io/release/v1.24.17/bin/linux/amd64/kubelet \
#        && chmod +x kubelet \
#        && mv kubelet /usr/bin/
# add missing package
# apt-get install ethtool socat conntrack -y


# install crictl
curl -LO https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.26.0/crictl-v1.26.0-linux-amd64.tar.gz \
        && tar xzvf crictl-v1.26.0-linux-amd64.tar.gz \
        && chmod +x crictl \
        && mv crictl /usr/bin/

# if encountered the following error during install kubelet
# W: An error occurred during the signature verification. The repository is not updated and the previous index files will be used. GPG error: https://packages.cloud.google.com/apt kubernetes-xenial InRelease: NO_PUBKEY B53DC80D13EDEF05
# W: https://apt.kubernetes.io/dists/kubernetes-xenial/InRelease...： NO_PUBKEY B53DC80D13EDEF05
# W: Some index files failed to download. They have been ignored, or old ones used instead.
#
# check this page https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/#install-using-native-package-management
# to see how it is fixed
#
##
# In short
# sudo mkdir -p /etc/apt/keyrings
# curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-archive-keyring.gpg
# echo "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
# sudo apt-get update && apt-get install -y kubectl
