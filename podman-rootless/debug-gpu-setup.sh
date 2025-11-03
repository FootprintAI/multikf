#!/bin/bash
# Debug script for GPU setup in kind cluster

echo "=========================================="
echo "1. Check device plugin status"
echo "=========================================="
kubectl get pods -n kube-system -l name=nvidia-device-plugin-ds
echo ""

echo "=========================================="
echo "2. Check device plugin logs (init container)"
echo "=========================================="
kubectl logs -n kube-system -l name=nvidia-device-plugin-ds -c nvidia-driver-installer --tail=20
echo ""

echo "=========================================="
echo "3. Check device plugin logs (main container)"
echo "=========================================="
kubectl logs -n kube-system -l name=nvidia-device-plugin-ds -c nvidia-device-plugin-ctr --tail=30
echo ""

echo "=========================================="
echo "4. Check node GPU capacity"
echo "=========================================="
kubectl get nodes -o json | grep -A 10 "nvidia.com/gpu" || echo "No GPU resources found on nodes"
echo ""

echo "=========================================="
echo "5. Check GPU devices in kind control-plane node"
echo "=========================================="
podman exec gpu-cluster4-control-plane ls -la /dev/nvidia* 2>&1
echo ""

echo "=========================================="
echo "6. Check host-lib mount in kind node"
echo "=========================================="
podman exec gpu-cluster4-control-plane ls -la /host-lib/libnvidia*.so* 2>&1 | head -10
echo ""

echo "=========================================="
echo "7. Check /host-dev mount in kind node"
echo "=========================================="
podman exec gpu-cluster4-control-plane ls -la /host-dev/nvidia* 2>&1
echo ""

echo "=========================================="
echo "8. Check if NVIDIA libraries are in kind node"
echo "=========================================="
podman exec gpu-cluster4-control-plane bash -c "ldconfig -p | grep nvidia" 2>&1 | head -10
echo ""

echo "=========================================="
echo "9. Check worker node (if exists)"
echo "=========================================="
if podman ps --filter "name=gpu-cluster4-worker" --format "{{.Names}}" | grep -q worker; then
    echo "Worker node found, checking devices..."
    podman exec gpu-cluster4-worker ls -la /dev/nvidia* 2>&1
    echo ""
    podman exec gpu-cluster4-worker ls -la /host-lib/libnvidia*.so* 2>&1 | head -10
else
    echo "No worker node found"
fi
echo ""

echo "=========================================="
echo "10. Check allocatable resources on nodes"
echo "=========================================="
kubectl describe nodes | grep -A 5 "Allocatable:" | grep -E "(Allocatable:|nvidia)"
echo ""
