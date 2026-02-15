# Container Image Mirror Tool

A Python utility for mirroring container images from public registries to private registries. This tool makes it easy to create a local mirror of required images, which is useful for air-gapped environments, rate limit mitigation, or speeding up deployments with a local cache.

## Features

- Mirror individual container images or process images in batch
- Support for authenticated registry access
- Preserves image path structure in the target registry
- Handles various image naming formats (with or without explicit registry, repo paths, tags, or digests)
- Option to execute commands directly or just print them for manual execution

## Installation

No special installation is required beyond Python 3. Simply download the script and make it executable:

```bash
chmod +x mirror_image.py
```

## Usage

### Basic Usage

Mirror a single image:

```bash
./mirror_image.py --image docker.io/kubeflow/training-operator:v1.5.0 --mirror reg.footprint-ai.com/kubeflow-mirror
```

This will output the commands needed to pull, tag, and push the image:

```bash
docker pull docker.io/kubeflow/training-operator:v1.5.0
docker tag docker.io/kubeflow/training-operator:v1.5.0 reg.footprint-ai.com/kubeflow-mirror/kubeflow/training-operator:v1.5.0
docker push reg.footprint-ai.com/kubeflow-mirror/kubeflow/training-operator:v1.5.0
```

### Mirror with Authentication

If your private registry requires authentication:

```bash
./mirror_image.py --image docker.io/kubeflow/training-operator:v1.5.0 --mirror reg.footprint-ai.com/kubeflow-mirror --username myuser --password mypass
```

This will add a login command before the other operations:

```bash
echo mypass | docker login reg.footprint-ai.com/kubeflow-mirror --username myuser --password-stdin
docker pull docker.io/kubeflow/training-operator:v1.5.0
docker tag docker.io/kubeflow/training-operator:v1.5.0 reg.footprint-ai.com/kubeflow-mirror/kubeflow/training-operator:v1.5.0
docker push reg.footprint-ai.com/kubeflow-mirror/kubeflow/training-operator:v1.5.0
```

### Batch Processing

To mirror multiple images, create a text file with one image per line:

```
# images.txt
docker.io/kubeflow/training-operator:v1.5.0
docker.io/kubeflow/katib-controller:v0.15.0
k8s.gcr.io/kube-scheduler:v1.21.0
quay.io/coreos/kube-state-metrics:v1.9.7
nginx:latest
```

Then process the batch file:

```bash
./mirror_image.py --batch-file images.txt --mirror reg.footprint-ai.com/kubeflow-mirror
```

### Execute Commands

To directly execute the commands instead of just printing them:

```bash
./mirror_image.py --image docker.io/kubeflow/training-operator:v1.5.0 --mirror reg.footprint-ai.com/kubeflow-mirror --execute
```

Output:

```
Executing: docker pull docker.io/kubeflow/training-operator:v1.5.0
Executing: docker tag docker.io/kubeflow/training-operator:v1.5.0 reg.footprint-ai.com/kubeflow-mirror/kubeflow/training-operator:v1.5.0
Executing: docker push reg.footprint-ai.com/kubeflow-mirror/kubeflow/training-operator:v1.5.0
```

## Command Line Arguments

| Argument | Description |
|----------|-------------|
| `--image` | Source image path (e.g., docker.io/kubeflow/training-operator:v1.5.0) |
| `--mirror` | Mirror registry (e.g., reg.footprint-ai.com/kubeflow-mirror) |
| `--username` | Registry username for authentication |
| `--password` | Registry password for authentication |
| `--batch-file` | File containing list of images to mirror (one per line) |
| `--execute` | Execute the commands instead of printing them |

## Image Name Handling

The tool handles various image formats:

- Images with explicit registry: `docker.io/kubeflow/training-operator:v1.5.0`
- Images with implicit registry: `kubeflow/training-operator:v1.5.0` (assumes docker.io)
- Official images: `nginx:latest` (assumes docker.io/library)
- Images with digest: `docker.io/kubeflow/training-operator@sha256:123abc...`
- Images without tag: `kubeflow/training-operator` (assumes :latest tag)
- Images with multi-level paths: `docker.io/kubeflow/common/katib-controller:v0.15.0`

## Integration with multikf

This tool complements the registry mirror feature in multikf. You can:

1. Use this script to populate your private registry with required images
2. Configure multikf to use your mirror registry:

```bash
multikf add my-cluster --with_registry_mirrors="docker.io|https://reg.footprint-ai.com/kubeflow-mirror"
```

## Use Cases

### Preparing for Air-gapped Deployments

```bash
# Mirror all required images
./mirror_image.py --batch-file required-images.txt --mirror registry.internal --execute

# Configure multikf to use the mirror
multikf add airgap-cluster --with_registry_mirrors="docker.io|https://registry.internal,k8s.gcr.io|https://registry.internal/k8s-mirror"
```

### Working Around Rate Limits

```bash
# Mirror images that might hit rate limits
./mirror_image.py --batch-file frequently-used-images.txt --mirror reg.footprint-ai.com/cache --execute

# Configure multikf to use the mirror
multikf add dev-cluster --with_registry_mirrors="docker.io|https://reg.footprint-ai.com/cache"
```

### Creating Project-specific Image Bundles

```bash
# Mirror project-specific images to a dedicated project path
./mirror_image.py --batch-file kubeflow-images.txt --mirror reg.footprint-ai.com/kubeflow-mirror --execute

# Configure multikf to use the project-specific mirror
multikf add kf-cluster --with_registry_mirrors="docker.io|https://reg.footprint-ai.com/kubeflow-mirror"
```

## License

[Apache License 2.0](LICENSE)

This project is licensed under the Apache License, Version 2.0. See the LICENSE file for the full license text.
