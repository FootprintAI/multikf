#!/usr/bin/env bash

# install nginx-ingress-controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.3.1/deploy/static/provider/baremetal/deploy.yaml

cat << EOF | tee haproxy.cfg
global
        # Maximum number of connections
        maxconn 3000
        # OS user to HAProxy
        user haproxy
        # OS group to HAProxy
        group haproxy
        # daemon mode on (background)
        daemon
        # Amount of processor cores used by HAProxy
        nbproc 1
        log /dev/log    local0
        log /dev/log    local1 notice
        chroot /var/lib/haproxy
        stats timeout 30s
# Parameters for frontend and backend
defaults
        log     global
        mode    http
        #option  httplog
        option  dontlognull
        timeout connect 5000
        timeout client  50000
        timeout server  50000
        errorfile 400 /etc/haproxy/errors/400.http
        errorfile 403 /etc/haproxy/errors/403.http
        errorfile 408 /etc/haproxy/errors/408.http
        errorfile 500 /etc/haproxy/errors/500.http
        errorfile 502 /etc/haproxy/errors/502.http
        errorfile 503 /etc/haproxy/errors/503.http
        errorfile 504 /etc/haproxy/errors/504.http
        # We are going to serve HTTP pages
        mode http
        # Enable statistics
        #stats enable
        # Define user:password to access statistics page (CHANGE IT and record it) <<<<<<<<<<
        #stats auth iamauser:iamapasswd
        # The statistics should be refreshed for each one second
        #stats refresh 1s
        # The URI to access statistics page
        #stats uri /stats                                                           <<<<<<<<<<<
        # Load balance method: Static Round Robin (Allow define different weights for different servers )
        balance static-rr
        # Forces HTTP 1.0
        option httpclose
        # Sets HAProxy to forward the user's IP to the application server
        # option forwardfor
        # Maximum backend connection time
        timeout connect 3000ms
        # Maximum wait time of the backend response
        timeout server 50000ms
        # Maximum waiting time of the user's communication to frontend (the firewall)
        timeout client 50000ms
        # Other parameters can be setted here
# Frontend definitions (entry host)
frontend http
        bind *:80
        # Maximum connections on the frontend
        maxconn 3000
        # The backend that should serve this port
        default_backend httpingresscontroller
frontend https
        bind *:443
        mode tcp
        # Maximum connections on the frontend
        maxconn 3000
        # The backend that should serve this port
        default_backend httpsingresscontroller
# Backend definitions (http servers: containers)
# The "*-big" servers have twice processing power and memory than "*-small" servers, so the weights are 2 and 1 respectively
backend httpingresscontroller
        server defaultcontroller   <use-ingress-nginx-controller-svc-private-ip>:80
backend httpsingresscontroller
        mode tcp
        server defaultcontroller   <use-ingress-nginx-controller-svc-private-ip>:443 check
EOF

```
kubectl get svc -n ingress-nginx
NAME                                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             NodePort    10.108.208.223   <none>        80:30523/TCP,443:30245/TCP   22h
ingress-nginx-controller-admission   ClusterIP   10.96.81.50      <none>        443/TCP                      22h
```


# check haproxy cfg
docker run -it \
    -v $(pwd)/haproxy.cfg:/usr/local/etc/haproxy/haproxy.cfg:ro \
    haproxytech/haproxy-alpine:2.3.2 \
    -c -f /usr/local/etc/haproxy/haproxy.cfg

# run haproxy 
docker run -d -p 80:80 -p 443:443 \
   -v $(pwd)/haproxy.cfg:/usr/local/etc/haproxy/haproxy.cfg:ro \
   haproxytech/haproxy-alpine:2.3.2
