#!/usr/bin/env bash


# install prerequsite

# install nfs-common
kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v1.3.2/deploy/prerequisite/longhorn-nfs-installation.yaml

# install open-iScsi
kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v1.3.2/deploy/prerequisite/longhorn-iscsi-installation.yaml

kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v1.3.2/deploy/longhorn.yaml

# uninstall
# kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v1.3.2/uninstall/uninstall.yaml
#
# waiting for the uninstall job running
# 
# kubectl delete -f https://raw.githubusercontent.com/longhorn/longhorn/v1.3.2/deploy/longhorn.yaml
# kubectl delete -f https://raw.githubusercontent.com/longhorn/longhorn/v1.3.2/uninstall/uninstall.yaml
#
# check this for longhorn deadlock on finalizer
# https://avasdream.engineer/kubernetes-longhorn-stuck-terminating

# make longhore storageclass as default
kubectl patch storageclass longhorn -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'

# Create default authentication for webui
# longhore provides a webui for us to fix data and replicas, but the following steps 
# will create an internet-facing ingress for longhore admin portal
# you may NOT want this to be public, so use portforward for internal usage.
#
# kubectl port-forward svc/longhorn-frontend -n longhorn-system 8080:80
#
#

USER=admin; PASSWORD=admin; echo "${USER}:$(openssl passwd -stdin -apr1 <<< ${PASSWORD})" >> auth
kubectl -n longhorn-system create secret generic basic-auth --from-file=auth


```longhorn-ingress.yml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: longhorn-ingress
  namespace: longhorn-system
  annotations:
    # type of authentication
    nginx.ingress.kubernetes.io/auth-type: basic
    # prevent the controller from redirecting (308) to HTTPS
    nginx.ingress.kubernetes.io/ssl-redirect: 'false'
    # name of the secret that contains the user/password definitions
    nginx.ingress.kubernetes.io/auth-secret: basic-auth
    # message to display with an appropriate context why the authentication is required
    nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required '
spec:
  rules:
  - http:
      paths:
      - pathType: Prefix
        path: "/"
        backend:
          service:
            name: longhorn-frontend
            port:
              number: 80
```
kubectl apply -f longhorn-ingress.yml

