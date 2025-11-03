#!/bin/bash
# Check what's actually in /host-lib in the kind node

echo "=========================================="
echo "Checking /host-lib content in kind node"
echo "=========================================="
echo ""

echo "Total files in /host-lib:"
podman exec gpu-cluster4-control-plane ls /host-lib | wc -l
echo ""

echo "Looking for nvidia files in /host-lib:"
podman exec gpu-cluster4-control-plane ls /host-lib | grep nvidia
echo ""

echo "Looking for cuda files in /host-lib:"
podman exec gpu-cluster4-control-plane ls /host-lib | grep cuda
echo ""

echo "Checking if libraries exist (first 20 files):"
podman exec gpu-cluster4-control-plane ls -la /host-lib | head -20
echo ""

echo "=========================================="
echo "What /host-lib is actually mounted to:"
echo "=========================================="
podman exec gpu-cluster4-control-plane mount | grep host-lib
echo ""

echo "=========================================="
echo "Check if we need a different mount path"
echo "=========================================="
echo "On HOST - where are NVIDIA driver libraries?"
ldconfig -p | grep nvidia | head -10
echo ""
