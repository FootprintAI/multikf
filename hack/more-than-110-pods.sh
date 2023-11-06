### for some reason you need to run more than 110 pods (default max pods count per node) for your workload
### you can use the following configuration to change it 
### run it with root permission
### ref: https://medium.com/@initcron/how-to-increase-the-number-of-pods-limit-per-kubernetes-node-877dcec5e4fa

```
vim /etc/systemd/system/kubelet.service.d/10-kubeadm.conf

- $KUBELET_CONFIG_ARGS $KUBELET_KUBEADM_ARGS $KUBELET_EXTRA_ARGS
+ $KUBELET_CONFIG_ARGS $KUBELET_KUBEADM_ARGS $KUBELET_EXTRA_ARGS --max-pods=243
```

systemctl restart kubelet
Warning: kubelet.service changed on disk. Run 'systemctl daemon-reload' to reload units.
systemctl daemon-reload
systemctl restart kubelet
