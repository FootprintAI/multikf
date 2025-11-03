#!/bin/bash
# Test if the current setup actually works despite missing host-lib/host-dev mounts

set -e

echo "=========================================="
echo "Testing current GPU setup"
echo "=========================================="
echo ""

echo "Step 1: Check device plugin status..."
kubectl get pods -n kube-system -l name=nvidia-device-plugin-ds
echo ""

echo "Step 2: Wait for device plugin to be ready..."
kubectl wait --for=condition=Ready pods -n kube-system -l name=nvidia-device-plugin-ds --timeout=60s || echo "Device plugin not ready yet"
echo ""

echo "Step 3: Check device plugin logs..."
echo "Init container logs:"
kubectl logs -n kube-system -l name=nvidia-device-plugin-ds -c nvidia-driver-installer --tail=10
echo ""
echo "Main container logs:"
kubectl logs -n kube-system -l name=nvidia-device-plugin-ds -c nvidia-device-plugin-ctr --tail=20
echo ""

echo "Step 4: Check if nodes report GPU capacity..."
kubectl get nodes -o=jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.capacity.nvidia\.com/gpu}{"\n"}{end}'
echo ""

echo "Step 5: Deploy test pod..."
kubectl delete pod gpu-test-cuda --ignore-not-found=true
kubectl apply -f podman-rootless/testpod.yaml
echo ""

echo "Step 6: Wait for test pod to start..."
sleep 5
kubectl wait --for=condition=Ready pod/gpu-test-cuda --timeout=60s || echo "Pod not ready yet"
echo ""

echo "Step 7: Check test pod logs..."
kubectl logs gpu-test-cuda
echo ""

echo "=========================================="
echo "Step 8: Manual check - exec into kind node"
echo "=========================================="
echo "Checking what's actually in the kind node:"
echo ""
echo "GPU devices in kind node:"
podman exec gpu-cluster4-control-plane ls -la /dev/nvidia*
echo ""
echo "Can we find NVIDIA libraries?"
podman exec gpu-cluster4-control-plane find /lib /usr/lib -name "libnvidia-ml.so*" 2>/dev/null | head -5 || echo "No NVIDIA libraries found in standard locations"
echo ""
