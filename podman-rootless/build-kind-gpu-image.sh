#!/bin/bash
# Build a custom Kind node image with NVIDIA libraries pre-installed

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
IMAGE_NAME="asia-east1-docker.pkg.dev/footprintai-dev/kafeido-mlops/kindest/node-cuda:v1.33.2"
BUILD_DIR="/tmp/kind-gpu-build"

echo "=========================================="
echo "Building custom Kind node image with GPU support"
echo "=========================================="
echo ""

# Create temporary build directory
echo "Step 1: Preparing build directory..."
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR/nvidia-libs"
echo ""

# Copy NVIDIA libraries from host
echo "Step 2: Collecting NVIDIA libraries from host..."
echo "  Copying from /lib/x86_64-linux-gnu..."
cp -L /lib/x86_64-linux-gnu/libnvidia*.so* "$BUILD_DIR/nvidia-libs/" 2>/dev/null || true
cp -L /lib/x86_64-linux-gnu/libcuda*.so* "$BUILD_DIR/nvidia-libs/" 2>/dev/null || true

echo "  Copying from /usr/lib/x86_64-linux-gnu..."
cp -L /usr/lib/x86_64-linux-gnu/libnvidia*.so* "$BUILD_DIR/nvidia-libs/" 2>/dev/null || true
cp -L /usr/lib/x86_64-linux-gnu/libcuda*.so* "$BUILD_DIR/nvidia-libs/" 2>/dev/null || true

LIB_COUNT=$(ls -1 "$BUILD_DIR/nvidia-libs/" 2>/dev/null | wc -l)
echo "  Found $LIB_COUNT NVIDIA library files"

if [ "$LIB_COUNT" -eq 0 ]; then
    echo ""
    echo "ERROR: No NVIDIA libraries found on the host!"
    echo "Please ensure NVIDIA drivers are installed."
    exit 1
fi
echo ""

# Create Dockerfile
echo "Step 3: Creating Dockerfile..."
cat > "$BUILD_DIR/Dockerfile" << 'EOF'
FROM kindest/node:v1.33.2

# Install prerequisites
RUN apt-get update && \
    apt-get install -y curl gnupg wget && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Create NVIDIA library directory
RUN mkdir -p /opt/nvidia/lib

# Copy NVIDIA libraries
COPY nvidia-libs/* /opt/nvidia/lib/

# Set up ldconfig
RUN echo "/opt/nvidia/lib" > /etc/ld.so.conf.d/nvidia.conf

# Add labels
LABEL description="Kind node image with NVIDIA GPU support"
LABEL nvidia.enabled="true"
LABEL version="1.33.2-gpu"
LABEL base.image="kindest/node:v1.33.2"
EOF
echo ""

# Build the image with Podman
echo "Step 4: Building image with Podman..."
cd "$BUILD_DIR"
podman build -t "$IMAGE_NAME" .
echo ""

# Verify the image
echo "Step 5: Verifying image..."
podman images | grep "node-cuda" || echo "Checking image..."
echo ""

# Test the image
echo "Step 6: Testing image..."
TEST_CONTAINER="test-kind-gpu-$(date +%s)"
podman run --rm --name "$TEST_CONTAINER" "$IMAGE_NAME" bash -c "
    echo 'Testing NVIDIA libraries in /opt/nvidia/lib:'
    ls -1 /opt/nvidia/lib | grep -E '(nvidia-ml|libcuda)' | head -5
    echo ''
    echo 'Total library files: '
    ls -1 /opt/nvidia/lib | wc -l
    echo ''
    echo 'Checking prerequisites:'
    which curl && echo '✓ curl installed'
    which wget && echo '✓ wget installed'
"
echo ""

# Cleanup
echo "Step 7: Cleaning up build directory..."
rm -rf "$BUILD_DIR"
echo ""

echo "=========================================="
echo "Build complete!"
echo "=========================================="
echo ""
echo "Image: $IMAGE_NAME"
echo ""
echo "Next steps:"
echo "1. Push the image to the registry:"
echo "   podman push $IMAGE_NAME"
echo ""
echo "2. The image is already configured in recreate-cluster.sh"
echo "   Just run: ./podman-rootless/recreate-cluster.sh"
echo ""
