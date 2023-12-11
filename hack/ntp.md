### install ntp


#### install ntp via apt-get 
```
sudo apt-get update && apt-get install -y ntp
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
