#### ssh local forwarding ####

Host ----(22/8222)----> Proxy ----(80/8080,443/8443)----> Host

Host:

```
ssh -R 8222:localhost:22 -i ~/.ssh/<proxy-rsa> ubuntu@<ip-address> -p 9527
```

Proxy:

```
ssh -L 0.0.0.0:8080:localhost:80 -L 0.0.0.0:8443:localhost:443 ubuntu@localhost -p 8222
```

Run check on proxy server:

```
netstat -antup|grep 8080


tcp        0      0 0.0.0.0:8080            0.0.0.0:*               LISTEN      1679/ssh
tcp        0      0 192.168.1.210:8080      192.168.1.243:49742     FIN_WAIT2   1679/ssh
```

noted that `0.0.0.0:8080` must be there so the telnet can work.
