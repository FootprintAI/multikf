### renew certificates inside cluster###

running on master node

```
kubeadm certs renew all
```

then restart kubelet running on the master node (or container if in kind env) with

```
systemctl restart kubelet
```


run this command to verify cluster certificates

```

kubeadm certs check-expiration

CERTIFICATE                EXPIRES                  RESIDUAL TIME   CERTIFICATE AUTHORITY   EXTERNALLY MANAGED
admin.conf                 Mar 01, 2025 02:27 UTC   364d            ca                      no
apiserver                  Mar 01, 2025 02:27 UTC   364d            ca                      no
apiserver-etcd-client      Mar 01, 2025 02:27 UTC   364d            etcd-ca                 no
apiserver-kubelet-client   Mar 01, 2025 02:27 UTC   364d            ca                      no
controller-manager.conf    Mar 01, 2025 02:27 UTC   364d            ca                      no
etcd-healthcheck-client    Mar 01, 2025 02:27 UTC   364d            etcd-ca                 no
etcd-peer                  Mar 01, 2025 02:27 UTC   364d            etcd-ca                 no
etcd-server                Mar 01, 2025 02:27 UTC   364d            etcd-ca                 no
front-proxy-client         Mar 01, 2025 02:27 UTC   364d            front-proxy-ca          no
scheduler.conf             Mar 01, 2025 02:27 UTC   364d            ca                      no

CERTIFICATE AUTHORITY   EXPIRES                  RESIDUAL TIME   EXTERNALLY MANAGED
ca                      Feb 15, 2033 17:36 UTC   8y              no
etcd-ca                 Feb 15, 2033 17:36 UTC   8y              no
front-proxy-ca          Feb 15, 2033 17:36 UTC   8y              no
```

Update admin.conf

`
mv etc/kubernetes/admin.conf /home/ubuntu/.kube/config
`

