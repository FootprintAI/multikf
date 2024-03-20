### configure timezone

```
sudo timedatectl set-timezone Asia/Taipei

// verify timezone with cmd

>> date
>> Mon Dec 25 16:55:59 CST 2023
```

### install ntp


#### install ntp via apt-get 
```
sudo apt-get update && sudo apt-get install -y ntp
```

### allow ntp port on firewall

```
sudo ufw allow 123/udp
```

### check ntp status

```
sudo systemctl status ntp
ntpq -p
```
