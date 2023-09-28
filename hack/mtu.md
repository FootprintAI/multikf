#### MTU


####  Test your MTU

Ping your gateway(192.168.1.1) with "do not fragment" option

```
ping 192.168.1.1 -c 2 -M do -s 2000

ping: local error: message too long, mtu=1452
```

and minues 28 bytes for IP header and ICMP header overhead.

```
ping 192.168.1.1 -c 2 -M do -s 1372
```

#### Update mtu permanently

##### static ip configuration

```
sudo vim /etc/netplan/00-installer-config.yaml
```



```
# This is the network config written by 'subiquity'
network:
  ethernets:
    enp1s0:
      dhcp4: true
    enp2s0:
      addresses:
      - 192.168.1.201/24
      gateway4: 192.168.1.1
      nameservers:
        addresses:
        - 8.8.8.8
        search: []
+      mtu: 1372
  version: 2
```

then reboot would works

##### dhcp configuation
sudo vim /etc/dhcp/dhclient.conf

```
interface "enp2s0" {
  default interface-mtu 1372;
  supersede interface-mtu 1372;
}
```

restart network manager

```
 sudo systemctl restart systemd-networkd
```

