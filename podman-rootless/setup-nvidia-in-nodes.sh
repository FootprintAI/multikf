#!/bin/bash
# Manually copy NVIDIA libraries into kind nodes (workaround for rootless podman mount issues)

set -e

echo "=========================================="
echo "Setting up NVIDIA libraries in kind nodes"
echo "=========================================="
echo ""

# Function to setup a node
setup_node() {
    local node=$1
    echo "Setting up node: $node"
    echo "-------------------"

    # Install prerequisites and copy libraries
    podman exec -i $node bash << 'SETUP'
set -e

echo "1. Installing prerequisites..."
apt-get update -qq
apt-get install -y curl gnupg wget ldconfig > /dev/null 2>&1

echo "2. Creating library directories..."
mkdir -p /usr/lib/x86_64-linux-gnu /usr/local/nvidia/lib64

echo "3. Checking what's in /host-lib (should be empty/wrong)..."
ls /host-lib 2>&1 | head -3

echo "4. Since bind mounts don't work in rootless podman, we need to copy from host..."
echo "   This will be done via podman cp from the host script"

SETUP

    echo ""
    echo "5. Copying NVIDIA libraries from host to node..."
    # Copy from host /lib/x86_64-linux-gnu
    for lib in /lib/x86_64-linux-gnu/libnvidia*.so* /lib/x86_64-linux-gnu/libcuda*.so*; do
        if [ -e "$lib" ]; then
            podman cp "$lib" "$node:/usr/lib/x86_64-linux-gnu/" 2>/dev/null || true
        fi
    done

    # Copy from host /usr/lib/x86_64-linux-gnu (critical - has libnvidia-ml.so)
    for lib in /usr/lib/x86_64-linux-gnu/libnvidia*.so* /usr/lib/x86_64-linux-gnu/libcuda*.so*; do
        if [ -e "$lib" ]; then
            podman cp "$lib" "$node:/usr/lib/x86_64-linux-gnu/" 2>/dev/null || true
        fi
    done

    # Also copy to nvidia lib path
    podman exec $node bash << 'SETUP2'
echo "6. Copying to /usr/local/nvidia/lib64..."
cp -a /usr/lib/x86_64-linux-gnu/libnvidia*.so* /usr/local/nvidia/lib64/ 2>/dev/null || true
cp -a /usr/lib/x86_64-linux-gnu/libcuda*.so* /usr/local/nvidia/lib64/ 2>/dev/null || true

echo "7. Updating library cache..."
cat > /etc/ld.so.conf.d/nvidia.conf << 'EOF'
/usr/lib/x86_64-linux-gnu
/usr/local/nvidia/lib64
EOF
ldconfig

echo "8. Verifying libraries..."
ldconfig -p | grep -E "nvidia-ml|libcuda" | head -5

echo "9. Checking GPU devices..."
ls -la /dev/nvidia* 2>&1 | head -5

echo "Setup complete for this node!"
SETUP2

    echo ""
}

# Setup control-plane
setup_node "gpu-cluster4-control-plane"
echo ""

# Setup worker if it exists
if podman ps --filter "name=gpu-cluster4-worker" --format "{{.Names}}" | grep -q worker; then
    setup_node "gpu-cluster4-worker"
fi

echo "=========================================="
echo "NVIDIA setup complete!"
echo "=========================================="
echo ""
echo "Next: Deploy the device plugin"
echo "  kubectl apply -f podman-rootless/nvidia-device-plugin-simple.yaml"
