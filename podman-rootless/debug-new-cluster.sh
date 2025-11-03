#!/bin/bash

echo "=========================================="
echo "Debugging new cluster mounts"
echo "=========================================="
echo ""

echo "1. Check podman mounts for control-plane:"
podman inspect gpu-cluster4-control-plane --format '{{json .Mounts}}' | python3 -m json.tool 2>/dev/null || podman inspect gpu-cluster4-control-plane --format '{{json .Mounts}}'
echo ""

echo "2. Check what directories exist in the node:"
podman exec gpu-cluster4-control-plane ls -la / | grep host
echo ""

echo "3. Check actual mounts inside the node:"
podman exec gpu-cluster4-control-plane mount | grep -E "(host|nvidia)"
echo ""

echo "4. Check if GPU devices are at least accessible:"
podman exec gpu-cluster4-control-plane ls -la /dev/nvidia*
echo ""

echo "5. Check device plugin logs:"
kubectl logs -n kube-system -l name=nvidia-device-plugin-ds -c nvidia-driver-installer --tail=30
echo ""

echo "6. Check device plugin main container logs:"
kubectl logs -n kube-system -l name=nvidia-device-plugin-ds -c nvidia-device-plugin-ctr --tail=20
echo ""
