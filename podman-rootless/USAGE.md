# GPU Kind Cluster - Usage Guide

## Quick Reference

### Using Pre-built Image (Default)

```bash
# Just run - uses pre-built image by default
./podman-rootless/recreate-cluster.sh
```

**Pros:**
- ✅ Fast cluster creation (no library copying)
- ✅ Consistent environment across all clusters
- ✅ NVIDIA libraries already included

**When to use:** Production, repeated testing, consistent environments

---

### Building Custom Image

```bash
# 1. Build the image
chmod +x podman-rootless/build-kind-gpu-image.sh
./podman-rootless/build-kind-gpu-image.sh

# 2. Push to your registry
podman push asia-east1-docker.pkg.dev/footprintai-dev/kafeido-mlops/kindest/node-cuda:v1.33.2

# 3. Use it (already default)
./podman-rootless/recreate-cluster.sh
```

**When to use:**
- First time setup
- NVIDIA driver version changed on host
- Want to customize the base image

---

### Using Standard Image (Copy Libraries)

```bash
# Override to use standard image
export USE_PREBUILT_IMAGE=false
export KIND_NODE_IMAGE=kindest/node:v1.33.2
./podman-rootless/recreate-cluster.sh
```

**Pros:**
- No custom image needed
- Always uses latest libraries from host

**Cons:**
- ⚠️ Slower cluster creation (copies ~47 libraries per node)
- Requires NVIDIA drivers on host

**When to use:** Testing, debugging, don't have access to pre-built image

---

## Configuration Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `USE_PREBUILT_IMAGE` | `true` | Skip library copying if using pre-built image |
| `KIND_NODE_IMAGE` | `asia-east1-docker.pkg.dev/footprintai-dev/kafeido-mlops/kindest/node-cuda:v1.33.2` | Kind node image to use |

### Examples

```bash
# Use different registry
export KIND_NODE_IMAGE=your-registry.com/kindest/node-cuda:v1.33.2
./podman-rootless/recreate-cluster.sh

# Force library copy even with custom image
export USE_PREBUILT_IMAGE=false
./podman-rootless/recreate-cluster.sh

# Use specific Kubernetes version (requires rebuilding image)
export KIND_NODE_IMAGE=kindest/node:v1.30.0
export USE_PREBUILT_IMAGE=false
./podman-rootless/recreate-cluster.sh
```

---

## Building the Pre-built Image

The `build-kind-gpu-image.sh` script:

1. ✅ Collects NVIDIA libraries from host
2. ✅ Creates Dockerfile based on `kindest/node:v1.33.2`
3. ✅ Copies libraries to `/opt/nvidia/lib` in image
4. ✅ Installs prerequisites (curl, wget, gnupg)
5. ✅ Builds with Podman
6. ✅ Tags as `asia-east1-docker.pkg.dev/footprintai-dev/kafeido-mlops/kindest/node-cuda:v1.33.2`

### Requirements

- NVIDIA drivers installed on host
- Podman with build capabilities
- Access to push to the registry

### What Gets Included

All `libnvidia*.so*` and `libcuda*.so*` from:
- `/lib/x86_64-linux-gnu/`
- `/usr/lib/x86_64-linux-gnu/`

Typically ~47 library files (~150MB total)

---

## Cluster Configuration

The cluster is configured with:

- **2 nodes**: 1 control-plane + 1 worker
- **maxPods**: 250 per node
- **GPU**: 1 GPU advertised per node (if available)
- **NVIDIA libraries**: `/opt/nvidia/lib` (in pre-built image or copied)

### Modifying Configuration

Edit `podman-rootless/gpu-kind-config.yaml`:

```yaml
# Add more workers
nodes:
  - role: control-plane
    kubeadmConfigPatches:
      - |
        kind: KubeletConfiguration
        maxPods: 250
  - role: worker
    kubeadmConfigPatches:
      - |
        kind: KubeletConfiguration
        maxPods: 250
  - role: worker  # Additional worker
    kubeadmConfigPatches:
      - |
        kind: KubeletConfiguration
        maxPods: 250
```

---

## Testing

```bash
# Quick test
./podman-rootless/test-gpu.sh

# Manual test
kubectl apply -f podman-rootless/testpod-with-libs.yaml
kubectl logs -f gpu-test-full
```

Expected output:
```
=== Summary ===
GPU devices: 8 devices found
NVIDIA libraries: 47 libraries mounted

✓ GPU passthrough is working!
```

---

## Troubleshooting

### Image not found

```bash
# Pull the pre-built image first
podman pull asia-east1-docker.pkg.dev/footprintai-dev/kafeido-mlops/kindest/node-cuda:v1.33.2

# Or build it locally
./podman-rootless/build-kind-gpu-image.sh
```

### Libraries not in pre-built image

```bash
# Rebuild the image
./podman-rootless/build-kind-gpu-image.sh

# Or fall back to copying
export USE_PREBUILT_IMAGE=false
./podman-rootless/recreate-cluster.sh
```

### GPU not detected

```bash
# Check device plugin
kubectl logs -n kube-system -l name=nvidia-device-plugin-ds

# Should see: "Detected NVML platform: found NVML library"
```

### Different NVIDIA driver version

```bash
# Rebuild image with current drivers
./podman-rootless/build-kind-gpu-image.sh

# Push updated image
podman push asia-east1-docker.pkg.dev/footprintai-dev/kafeido-mlops/kindest/node-cuda:v1.33.2
```

---

## Best Practices

1. **Use pre-built image** - Faster, more consistent
2. **Rebuild after driver updates** - Ensures compatibility
3. **Version your images** - Tag with date or driver version
4. **Test after building** - Run `test-gpu.sh` to verify
5. **Document your setup** - Note driver versions, image tags

---

## Image Versioning Strategy

Consider tagging with driver version:

```bash
# Build
./podman-rootless/build-kind-gpu-image.sh

# Tag with driver version
DRIVER_VERSION=$(nvidia-smi --query-gpu=driver_version --format=csv,noheader | head -1)
podman tag asia-east1-docker.pkg.dev/footprintai-dev/kafeido-mlops/kindest/node-cuda:v1.33.2 \
            asia-east1-docker.pkg.dev/footprintai-dev/kafeido-mlops/kindest/node-cuda:v1.33.2-driver${DRIVER_VERSION}

# Push both
podman push asia-east1-docker.pkg.dev/footprintai-dev/kafeido-mlops/kindest/node-cuda:v1.33.2
podman push asia-east1-docker.pkg.dev/footprintai-dev/kafeido-mlops/kindest/node-cuda:v1.33.2-driver${DRIVER_VERSION}
```

Then specify version:
```bash
export KIND_NODE_IMAGE=asia-east1-docker.pkg.dev/footprintai-dev/kafeido-mlops/kindest/node-cuda:v1.33.2-driver560.35.05
./podman-rootless/recreate-cluster.sh
```
