# GPU Support for Kind with Rootless Podman

This directory contains scripts and configurations for running GPU-enabled Kubernetes workloads in a Kind cluster using rootless Podman.

## ✅ Current Status

GPU passthrough is **working**! The setup provides:
- ✅ GPU device access (`/dev/nvidia*`) in pods
- ✅ NVIDIA driver libraries (47 libraries available)
- ✅ CUDA support via libcuda.so
- ✅ Device plugin successfully detecting and advertising GPUs

## 🚀 Quick Start

### Option 1: Using Pre-built Image (Recommended)

```bash
cd ~/multikf/multikf
# Uses pre-built image by default, skips library copying
./podman-rootless/recreate-cluster.sh
```

### Option 2: Building from Scratch

```bash
# Build the custom image with NVIDIA libraries
./podman-rootless/build-kind-gpu-image.sh

# Push to registry
podman push asia-east1-docker.pkg.dev/footprintai-dev/kafeido-mlops/kindest/node-cuda:v1.33.2

# Use the image
./podman-rootless/recreate-cluster.sh
```

### Option 3: Using Standard Image + Library Copy

```bash
# Set environment variables to use standard image and copy libraries
export USE_PREBUILT_IMAGE=false
export KIND_NODE_IMAGE=kindest/node:v1.33.2
./podman-rootless/recreate-cluster.sh
```

The `recreate-cluster.sh` script will:
1. Delete any existing cluster
2. Create a new Kind cluster with GPU support
3. Either use pre-built image OR copy NVIDIA libraries to `/opt/nvidia/lib`
4. Deploy the NVIDIA device plugin
5. Verify GPU capacity on nodes

### Test GPU Access

```bash
./podman-rootless/test-gpu.sh
```

This will deploy a test pod and verify GPU devices and libraries are accessible.

## 📋 How It Works

### The Rootless Podman Challenge

Rootless Podman doesn't support bind mounts the same way Docker does. When using Kind's `extraMounts`, the mounts resolve to the container's root filesystem instead of the host directories.

### Our Solution

1. **Isolated Library Directory**: Copy only NVIDIA/CUDA libraries to `/opt/nvidia/lib` in each Kind node
   - Avoids glibc version conflicts
   - No system library interference

2. **Direct Library Copying**: Use `podman cp` to copy libraries from host to nodes
   - Bypasses broken bind mount behavior
   - Works reliably in rootless mode

3. **Device Plugin Library Access**: Mount `/opt/nvidia/lib` into the device plugin pod
   - Device plugin can load `libnvidia-ml.so.1`
   - Successfully detects and advertises GPUs

## 🔧 Running GPU Workloads

To use GPUs in your pods, follow this pattern:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-gpu-workload
spec:
  volumes:
  - name: nvidia-libs
    hostPath:
      path: /opt/nvidia/lib
  containers:
  - name: my-container
    image: your-cuda-image:latest
    command: ["your-command"]
    env:
    - name: LD_LIBRARY_PATH
      value: "/opt/nvidia/lib:/usr/local/nvidia/lib:/usr/local/nvidia/lib64"
    volumeMounts:
    - name: nvidia-libs
      mountPath: /opt/nvidia/lib
      readOnly: true
    resources:
      limits:
        nvidia.com/gpu: 1
    securityContext:
      privileged: true
```

**Key requirements:**
1. Mount `/opt/nvidia/lib` as a volume
2. Set `LD_LIBRARY_PATH` to include `/opt/nvidia/lib`
3. Request GPU in `resources.limits`
4. Use `privileged: true` for device access

## 📁 Files in This Directory

### Main Scripts
- **`recreate-cluster.sh`** - Main script to create/recreate the cluster with GPU support
- **`test-gpu.sh`** - Test GPU access and verify setup
- **`apply-fixed-device-plugin.sh`** - Redeploy device plugin after changes

### Configuration Files
- **`gpu-kind-config.yaml`** - Kind cluster configuration with extraMounts (note: doesn't work fully in rootless Podman)
- **`nvidia-device-plugin-simple.yaml`** - NVIDIA device plugin DaemonSet
- **`testpod.yaml`** - Basic GPU test pod (devices only)
- **`testpod-with-libs.yaml`** - Full GPU test pod (devices + libraries)

### Debug Scripts
- **`debug-gpu-setup.sh`** - Comprehensive GPU setup diagnostics
- **`debug-libs-in-nodes.sh`** - Check if libraries were copied correctly
- **`check-mounts.sh`** - Verify mount points in Kind nodes
- **`find-nvidia-libs.sh`** - Locate NVIDIA libraries on host
- **`check-host-lib-content.sh`** - Check mounted directory contents

### Other Files
- **`install.txt`** - Original installation notes and manual steps
- **`README.md`** - This file

## 🐛 Troubleshooting

### GPU not detected in pods

Check device plugin logs:
```bash
kubectl logs -n kube-system -l name=nvidia-device-plugin-ds
```

Should see: `Detected NVML platform: found NVML library`

### Libraries not found

Verify libraries were copied:
```bash
./podman-rootless/debug-libs-in-nodes.sh
```

Check that `/opt/nvidia/lib` contains NVIDIA libraries in both control-plane and worker nodes.

### Device plugin failing to start

Check for glibc conflicts:
```bash
kubectl logs -n kube-system -l name=nvidia-device-plugin-ds | grep GLIBC
```

If you see glibc errors, the libraries directory might contain system libraries. Only NVIDIA/CUDA libraries should be in `/opt/nvidia/lib`.

### Recreate from scratch

```bash
./podman-rootless/recreate-cluster.sh
```

This will clean up and recreate everything.

## 📝 Technical Details

### Why Not Use Standard Approaches?

1. **CDI (Container Device Interface)**: Requires NVIDIA Container Toolkit in the container runtime, which is complex to set up in rootless Podman
2. **NVIDIA GPU Operator**: Designed for standard Kubernetes, not Kind with rootless Podman
3. **Direct bind mounts**: Don't work correctly in rootless Podman due to namespace mapping

### Our Approach Benefits

- ✅ Works with rootless Podman
- ✅ No NVIDIA Container Toolkit required in Podman
- ✅ Simple, reproducible setup
- ✅ Easy to debug and understand
- ✅ Libraries isolated to avoid conflicts

### Limitations

- Libraries must be manually copied (automated in our scripts)
- Each workload pod must mount `/opt/nvidia/lib`
- Requires `privileged: true` for GPU device access
- Not suitable for multi-tenant environments (due to privileged requirement)

## 🎯 Next Steps

Now that GPU support is working, you can:

1. Deploy CUDA-based workloads
2. Run ML/AI training jobs
3. Use GPU-accelerated applications
4. Test with frameworks like PyTorch, TensorFlow, etc.

Example ML workload pattern in `testpod-with-libs.yaml`.

## 📚 References

- [Kind Documentation](https://kind.sigs.k8s.io/)
- [NVIDIA Device Plugin](https://github.com/NVIDIA/k8s-device-plugin)
- [Podman Documentation](https://docs.podman.io/)
- Original setup notes: `install.txt`

---

**Status**: ✅ Working and tested
**Last Updated**: 2025-11-03
**NVIDIA Driver Version**: 560.35.05
