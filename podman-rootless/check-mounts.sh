#!/bin/bash
# Check what mounts actually exist in the kind nodes

echo "=========================================="
echo "Checking mounts in control-plane node"
echo "=========================================="
echo "All mounts:"
podman exec gpu-cluster4-control-plane mount | grep -E "(host-lib|host-dev|nvidia)" || echo "No matching mounts found"
echo ""

echo "Directory listing:"
podman exec gpu-cluster4-control-plane ls -la / | grep -E "host" || echo "No host-* directories"
echo ""

echo "Check /dev directly:"
podman exec gpu-cluster4-control-plane ls -la /dev/nvidia* 2>&1
echo ""

echo "Check host libraries location:"
podman exec gpu-cluster4-control-plane find /lib -name "libnvidia*.so*" 2>/dev/null | head -5
podman exec gpu-cluster4-control-plane find /usr -name "libnvidia*.so*" 2>/dev/null | head -5
echo ""

echo "=========================================="
echo "Checking Podman mount info"
echo "=========================================="
podman inspect gpu-cluster4-control-plane --format '{{json .Mounts}}' | python3 -m json.tool 2>/dev/null || podman inspect gpu-cluster4-control-plane --format '{{json .Mounts}}'
echo ""

echo "=========================================="
echo "Checking kind config that was used"
echo "=========================================="
podman exec gpu-cluster4-control-plane cat /kind/kubeadm.conf 2>&1 | head -20
