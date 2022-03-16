# multikf
Multi-Kind leverages [Vagrant](https://github.com/hashicorp/vagrant) and [Kind](https://github.com/kubernetes-sigs/kind) (Kubernetes In Docker) to create multiple local kubernetes and kubeflow clusters inside the same host machine, see the following png for simple layout
![flow](./images/intro.png)

#### Why we need this?

As a machine gets more powerful, it is such a waste to have it running just one Kubernetes, especially for the applications which require only a local Kubernetes for practice. One example is our [Kubeflow workshop](https://github.com/footprintai/kubeflow-workshop).
To fully utilize hardware resources, we leverage vagrant to construct a fully isolated environment and install required packages on it (e.g. Kubernetes and Kubeflow and more ...), map ports for kubeApi and ssh, and also export its kubeconfg to host. Therefore, users on the host machine can easily talk to the guest Kube-API via kubectl.

#### When we need this?

We expected the user are under:

- Windows environemnt with docker desktop installed
- Linux environment

and this tool provides abstractions for them to operate clusters.

#### Why Vagrant is required?

Idealy, we could just use Kind which running as a container to provide resource isolation. However, Kind was unable to isolate resources from its underlying kubelet(see [issue](https://github.com/kubernetes-sigs/kind/issues/877)) due to kubelet's implementation. Thus, Vagrant is served as a resource isolation and provide clean guest enviornment.

#### How to use?

##### Add a vagrant machine named test000 with 1 cpu and 1G memory.

```
./multikf add test000 --cpu 1 --memory 1
```

##### Export a vargant machine's kubeconfig
```
./multikf export test000 --kubeconfig_path /tmp/test000.kubeconfig

run kubectl from host

 kubectl get pods --all-namespaces --kubeconfig=/tmp/test000.kubeconfig
```


##### list machines

```
./multikf list

+---------+------------------+---------+------+---------------+
|  NAME   |       DIR        | STATUS  | CPUS |    MEMORY     |
+---------+------------------+---------+------+---------------+
| test000 | .vagrant/test000 | running |    1 | 70720/1000328 |
+---------+------------------+---------+------+---------------+
```

##### delete a machine

```
./multikf delete test000

```

#### connect a machine

```
./multikf connect kubeflow test000

```

#### Roadmap

Fields listed here is on our roadmap.

| Fields | Supported |
|------|------|
| Cpu Isolation | O |
| Memory Isolation | O |
| GPU Isolation | X |


#### Gpu Passthough

For passing gpu to docker container, one approach is to use `--gpus=all` when you launched docker container like

```
docker run -it --gpus=all ubuntu:21.10 /bin/bash

```
where it relies on the host's cuda driver.
However, Kind are NOT supported this approach, see [issue](https://github.com/kubernetes-sigs/kind/pull/1886)
