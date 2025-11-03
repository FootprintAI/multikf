#!/bin/bash

echo "Mounts exist in podman inspect, but ls fails. Debugging..."
echo ""

echo "1. Check if directories exist:"
podman exec gpu-cluster4-control-plane bash -c "ls -ld /host-lib /host-usr-lib /host-dev 2>&1"
echo ""

echo "2. Check what's IN the directories:"
podman exec gpu-cluster4-control-plane bash -c "ls /host-lib 2>&1 | head -5"
echo ""
podman exec gpu-cluster4-control-plane bash -c "ls /host-usr-lib 2>&1 | head -5"
echo ""

echo "3. Check mounts from inside the container:"
podman exec gpu-cluster4-control-plane mount | grep host
echo ""

echo "4. Try to find nvidia files differently:"
podman exec gpu-cluster4-control-plane bash -c "find /host-lib -name '*nvidia*' 2>&1 | head -5"
echo ""
podman exec gpu-cluster4-control-plane bash -c "find /host-usr-lib -name '*nvidia*' 2>&1 | head -5"
echo ""

echo "5. Check the actual host to see what SHOULD be there:"
echo "Host /lib/x86_64-linux-gnu:"
ls /lib/x86_64-linux-gnu/ | grep nvidia | head -5
echo ""
echo "Host /usr/lib/x86_64-linux-gnu:"
ls /usr/lib/x86_64-linux-gnu/ | grep nvidia | head -5
echo ""

echo "6. Check device plugin init container logs:"
kubectl logs -n kube-system -l name=nvidia-device-plugin-ds -c nvidia-driver-installer
echo ""
