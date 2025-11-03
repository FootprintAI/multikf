#!/bin/bash

echo "Checking what's actually mounted where..."
echo ""

echo "1. All mounts related to host-lib, host-usr-lib, host-dev:"
podman exec gpu-cluster4-control-plane mount | grep -E "host-lib|host-usr-lib|host-dev"
echo ""

echo "2. What device is /dev/vda2?"
podman exec gpu-cluster4-control-plane df -h | grep vda2
echo ""

echo "3. Check if the mount points are being used by something else:"
podman exec gpu-cluster4-control-plane bash -c "mount | grep '/host-lib '"
podman exec gpu-cluster4-control-plane bash -c "mount | grep '/host-usr-lib '"
echo ""

echo "4. Let's check the root filesystem:"
podman exec gpu-cluster4-control-plane df -h /
echo ""

echo "5. Try accessing the directory with full path:"
podman exec gpu-cluster4-control-plane bash -c "ls -la /host-lib/ 2>&1 | head -10"
echo ""

echo "6. Check if it's a symlink or something:"
podman exec gpu-cluster4-control-plane bash -c "file /host-lib /host-usr-lib"
echo ""
