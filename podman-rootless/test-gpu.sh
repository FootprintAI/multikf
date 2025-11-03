#!/bin/bash
# Test GPU access in the cluster

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================="
echo "Testing GPU access in Kubernetes"
echo "=========================================="
echo ""

echo "Step 1: Checking GPU capacity on nodes..."
kubectl get nodes -o=jsonpath='{range .items[*]}{.metadata.name}{": "}{.status.capacity.nvidia\.com/gpu}{" GPU(s)\n"}{end}'
echo ""

echo "Step 2: Deploying GPU test pod with library mounts..."
kubectl delete pod gpu-test-full --ignore-not-found
kubectl apply -f "$SCRIPT_DIR/testpod-with-libs.yaml"
echo ""

echo "Step 3: Waiting for pod to complete..."
kubectl wait --for=condition=Ready pod/gpu-test-full --timeout=30s || echo "Pod may still be starting..."
sleep 3
echo ""

echo "Step 4: Checking test results..."
echo "=========================================="
kubectl logs gpu-test-full
echo "=========================================="
echo ""

echo "Test complete!"
echo ""
echo "If you see '✓ GPU passthrough is working!' above, you're all set!"
echo ""
echo "To run your own GPU workloads:"
echo "1. Add 'nvidia.com/gpu: 1' to resources.limits"
echo "2. Mount /opt/nvidia/lib as a volume"
echo "3. Set LD_LIBRARY_PATH to include /opt/nvidia/lib"
echo ""
