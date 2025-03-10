## troubleshooting ##

#### Booting issue

#####Error message showing when the machine just booted.
```
/dev/nvme0n1p2: clean xxx/xxxx files, xxxx/xxxxx/blocks
```

This is because the graph card driver issues where the system was loaded in-correct graph card drivers. you can check relevant solutions here[[1](https://askubuntu.com/questions/1277842/ubuntu-20-04-stuck-at-dev-nvme0n1p5-clean-xxx-xxx-files-xxx-xxx-blocks-duri)].

To fix it, boot with recovery mode and update grub config with root permission, and remove `quiet splash` with `nomodeset` to enable boot with lower resolution to avoid driver issue.

```
gedit /etc/default/grub

- GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
+ GRUB_CMDLINE_LINUX_DEFAULT="nomodeset"

```

then update grub with `sudo update-grub`. then reboot would solve the issue.


#### disable ipv6

```
sudo vim /etc/sysctl.conf

+ net.ipv6.conf.all.disable_ipv6=1
+ net.ipv6.conf.default.disable_ipv6=1
+ net.ipv6.conf.lo.disable_ipv6=1

```
then run `sudo sysctl -p` to update such changes into system

#### '/sbin/ldconfig.real': No such file or directory

during installation gpu-operator on kind (kubernetes in docker), `nvidia-operator-validator-67zsz` would failed to run due to `/sbin/ldconfig.real` is not found.
the simple workaround would be fixed with the symbolic link: `ln -s /sbin/ldconfig /sbin/ldconfig.real`

#### etcd keep rebooting ####

when your etcd is keep rebooting, especially you are using kubeadm with v1.27+,
you might want to turn the flag on by editing the file `/etc/containerd/config.toml` :

```
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
    SystemdCgroup = true
```

ref issue: https://github.com/etcd-io/etcd/issues/13670

#### node netowrk issue on kind cluster

when your node has network issue where pods are unable to reach each other,
try: 

```
sysctl net.bridge.bridge-nf-call-iptables=0
sysctl net.bridge.bridge-nf-call-arptables=0
sysctl net.bridge.bridge-nf-call-ip6tables=0
```
 to see how it goes.
 ref: https://github.com/kubernetes-sigs/kind/issues/2886#issuecomment-1219158523
