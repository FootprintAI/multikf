#!/bin/bash

echo "Checking if libraries were actually copied to nodes..."
echo ""

echo "1. Check control-plane node:"
echo "----------------------------"
podman exec gpu-cluster4-control-plane bash -c "ldconfig -p | grep nvidia-ml"
echo ""
podman exec gpu-cluster4-control-plane bash -c "ls -la /usr/lib/x86_64-linux-gnu/libnvidia-ml* 2>&1"
echo ""

echo "2. Check worker node:"
echo "----------------------------"
podman exec gpu-cluster4-worker bash -c "ldconfig -p | grep nvidia-ml"
echo ""
podman exec gpu-cluster4-worker bash -c "ls -la /usr/lib/x86_64-linux-gnu/libnvidia-ml* 2>&1"
echo ""

echo "3. Check what device plugin pod sees:"
echo "----------------------------"
POD=$(kubectl get pods -n kube-system -l name=nvidia-device-plugin-ds -o jsonpath='{.items[0].metadata.name}')
echo "Checking pod: $POD"
kubectl exec -n kube-system $POD -- bash -c "ldconfig -p | grep nvidia" 2>&1 || echo "Can't exec into pod"
echo ""

echo "4. Check LD_LIBRARY_PATH in device plugin pod:"
kubectl exec -n kube-system $POD -- bash -c "echo \$LD_LIBRARY_PATH" 2>&1 || echo "Can't exec into pod"
echo ""

echo "5. Try to find libnvidia-ml.so in device plugin pod:"
kubectl exec -n kube-system $POD -- bash -c "find /usr -name 'libnvidia-ml.so*' 2>&1" 2>&1 || echo "Can't exec into pod"
echo ""
