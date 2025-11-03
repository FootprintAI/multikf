#!/bin/bash
# Apply the fixed device plugin with library mounts

set -e

echo "=========================================="
echo "Applying fixed NVIDIA device plugin"
echo "=========================================="
echo ""

echo "Step 1: Deleting old device plugin..."
kubectl delete daemonset nvidia-device-plugin-daemonset -n kube-system || echo "No existing daemonset found"
echo ""

echo "Step 2: Waiting for old pods to terminate..."
sleep 5
echo ""

echo "Step 3: Applying updated device plugin..."
kubectl apply -f "$(dirname "$0")/nvidia-device-plugin-simple.yaml"
echo ""

echo "Step 4: Waiting for device plugin to be ready..."
sleep 5
kubectl wait --for=condition=Ready pods -n kube-system -l name=nvidia-device-plugin-ds --timeout=60s
echo ""

echo "Step 5: Checking device plugin logs..."
kubectl logs -n kube-system -l name=nvidia-device-plugin-ds --tail=20
echo ""

echo "Step 6: Verifying GPU capacity on nodes..."
kubectl get nodes -o=jsonpath='{range .items[*]}{.metadata.name}{": "}{.status.capacity.nvidia\.com/gpu}{"\n"}{end}'
echo ""

echo "=========================================="
echo "Device plugin update complete!"
echo "=========================================="
echo ""
echo "If you see GPU capacity above (e.g., 'gpu-cluster4-control-plane: 1'), you're ready to test!"
echo ""
echo "Test with:"
echo "  kubectl delete pod gpu-test-cuda --ignore-not-found"
echo "  kubectl apply -f $(dirname "$0")/testpod.yaml"
echo "  kubectl logs -f gpu-test-cuda"
echo ""
