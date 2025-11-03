#!/bin/bash
# Script to recreate the kind cluster with GPU support for rootless Podman
# This script handles the rootless Podman mount limitation by manually copying libraries

set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================="
echo "Recreating kind cluster with GPU support"
echo "=========================================="
echo ""

# Delete existing cluster
echo "Step 1: Deleting existing cluster..."
kind delete cluster --name gpu-cluster4 || echo "Cluster doesn't exist, continuing..."
echo ""

# Create new cluster (basic config, mounts won't work in rootless podman)
echo "Step 2: Creating new cluster..."
KIND_EXPERIMENTAL_PROVIDER=podman kind create cluster \
  --name gpu-cluster4 \
  --config "$SCRIPT_DIR/gpu-kind-config.yaml" \
  --image kindest/node:v1.33.2
echo ""

# Wait for cluster to be ready
echo "Step 3: Waiting for cluster to be ready..."
kubectl wait --for=condition=Ready nodes --all --timeout=120s
echo ""

# Setup NVIDIA libraries in nodes (workaround for rootless podman mount issues)
echo "Step 4: Setting up NVIDIA libraries in kind nodes..."
echo "Note: Bind mounts don't work properly in rootless Podman, so we copy libraries directly"
echo ""

# Function to setup a node
setup_node() {
    local node=$1
    echo "Setting up node: $node"
    echo "-------------------"

    # Install prerequisites and create directories
    podman exec -i $node bash << 'SETUP'
set -e
apt-get update -qq
apt-get install -y curl gnupg wget > /dev/null 2>&1
mkdir -p /opt/nvidia/lib
SETUP

    echo "  Copying NVIDIA libraries from host to isolated directory..."
    # Copy from host /lib/x86_64-linux-gnu - only NVIDIA and CUDA libraries
    for lib in /lib/x86_64-linux-gnu/libnvidia*.so* /lib/x86_64-linux-gnu/libcuda*.so*; do
        if [ -e "$lib" ]; then
            podman cp "$lib" "$node:/opt/nvidia/lib/" 2>/dev/null || true
        fi
    done

    # Copy from host /usr/lib/x86_64-linux-gnu (critical - has libnvidia-ml.so) - only NVIDIA and CUDA libraries
    for lib in /usr/lib/x86_64-linux-gnu/libnvidia*.so* /usr/lib/x86_64-linux-gnu/libcuda*.so*; do
        if [ -e "$lib" ]; then
            podman cp "$lib" "$node:/opt/nvidia/lib/" 2>/dev/null || true
        fi
    done

    # Verify
    podman exec $node bash << 'SETUP2'
echo "  Verifying libraries in /opt/nvidia/lib..."
ls -la /opt/nvidia/lib/ | grep -E "nvidia-ml|libcuda" | head -5 || echo "  WARNING: Critical libraries not found!"
echo "  Checking GPU devices..."
ls -la /dev/nvidia* 2>&1 | head -3 || echo "  WARNING: GPU devices not found!"
SETUP2
    echo ""
}

# Setup control-plane
setup_node "gpu-cluster4-control-plane"

# Setup worker if it exists
if podman ps --filter "name=gpu-cluster4-worker" --format "{{.Names}}" | grep -q worker; then
    setup_node "gpu-cluster4-worker"
fi

# Deploy device plugin
echo "Step 5: Deploying NVIDIA device plugin..."
kubectl apply -f "$SCRIPT_DIR/nvidia-device-plugin-simple.yaml"
echo ""

# Wait for device plugin to be ready
echo "Step 6: Waiting for device plugin to be ready..."
sleep 5
kubectl wait --for=condition=Ready pods -n kube-system -l name=nvidia-device-plugin-ds --timeout=120s || echo "Device plugin not ready yet, check logs"
echo ""

# Check GPU capacity
echo "Step 7: Verifying GPU capacity on nodes..."
kubectl get nodes -o=jsonpath='{range .items[*]}{.metadata.name}{": "}{.status.capacity.nvidia\.com/gpu}{"\n"}{end}'
echo ""

echo "=========================================="
echo "Cluster recreation complete!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Verify device plugin logs:"
echo "   kubectl logs -n kube-system -l name=nvidia-device-plugin-ds"
echo ""
echo "2. Test with the GPU test pod:"
echo "   kubectl delete pod gpu-test-cuda --ignore-not-found"
echo "   kubectl apply -f $SCRIPT_DIR/testpod.yaml"
echo "   kubectl logs -f gpu-test-cuda"
echo ""
