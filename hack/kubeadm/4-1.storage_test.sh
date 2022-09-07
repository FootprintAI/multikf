#!/usr/bin/env bash

# we leverage ioszone to perform load test, you can run this pod in the newly created cluster

```
# File: iozone.yaml (use default storageclass)
---
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: iozonev1
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 5Gi
---
apiVersion: v1
kind: Pod
metadata:
  name: iozone
  labels:
    app: benchmarktest
spec:
  containers:
  - name: iozone
    image: docker.io/pstauffer/iozone:v1.0
    imagePullPolicy: IfNotPresent
    command: ["iozone"]
    args:
        - -w
        - -c
        - -e
        - -i 0
        - -+n
        - -C
        - -r 64k
        - -s 1g
        - -t 16
        - +p 60
    volumeMounts:
    - mountPath: "/mnt/flusterfs"
      name: csivol
  volumes:
  - name: csivol
    persistentVolumeClaim:
      claimName: iozonev1
  restartPolicy: OnFailure
```

kubectl apply -f iozone.yaml
