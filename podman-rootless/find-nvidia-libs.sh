#!/bin/bash
# Find where NVIDIA libraries are actually located on the host

echo "=========================================="
echo "Finding NVIDIA libraries on HOST system"
echo "=========================================="
echo ""

echo "Checking /lib/x86_64-linux-gnu (what's mounted as /host-lib):"
ls -la /lib/x86_64-linux-gnu/libnvidia*.so* 2>&1 | head -10
echo ""

echo "Checking /usr/lib/x86_64-linux-gnu:"
ls -la /usr/lib/x86_64-linux-gnu/libnvidia*.so* 2>&1 | head -10
echo ""

echo "Searching for libnvidia-ml.so across the system:"
find /lib /usr/lib -name "libnvidia-ml.so*" 2>/dev/null
echo ""

echo "Checking what the kind node sees in /host-lib:"
podman exec gpu-cluster4-control-plane ls -la /host-lib/ | grep nvidia | head -10
echo ""

echo "Checking what the kind node sees in its own /lib and /usr/lib:"
podman exec gpu-cluster4-control-plane find /lib /usr/lib -name "libnvidia*.so*" 2>/dev/null | head -10
echo ""
