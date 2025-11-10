#!/bin/bash
# Diagnostic script for rootless podman + kind issues on RHEL 9.4

echo "=========================================="
echo "Rootless Podman + Kind Diagnostics"
echo "=========================================="
echo ""

echo "1. OS Version:"
cat /etc/os-release | grep -E "^(NAME|VERSION)="
echo ""

echo "2. Cgroup Version:"
stat -fc %T /sys/fs/cgroup/
echo ""

echo "3. Podman Version:"
podman --version
echo ""

echo "4. Kind Version:"
kind --version
echo ""

echo "5. User Cgroup Controllers (should include: cpu cpuset io memory pids):"
cat /sys/fs/cgroup/user.slice/user-$(id -u).slice/user@$(id -u).service/cgroup.controllers
echo ""

echo "6. Systemd Delegation Config:"
if [ -f /etc/systemd/system/user@.service.d/delegate.conf ]; then
    echo "✓ Delegate config exists:"
    cat /etc/systemd/system/user@.service.d/delegate.conf
else
    echo "✗ Delegate config NOT found at /etc/systemd/system/user@.service.d/delegate.conf"
fi
echo ""

echo "7. Podman Info (cgroup info):"
podman info | grep -A 5 "cgroupManager"
echo ""

echo "8. Test systemd in container:"
echo "Testing if systemd can run in rootless podman..."
TEST_RESULT=$(podman run --rm --systemd=always --name test-systemd kindest/node:v1.33.2 /sbin/init & sleep 5; podman logs test-systemd 2>&1 | grep -i "systemd\|multi-user\|failed" || echo "No systemd logs found"; podman stop test-systemd 2>/dev/null; podman rm -f test-systemd 2>/dev/null)
echo "$TEST_RESULT"
echo ""

echo "=========================================="
echo "Diagnostics Complete"
echo "=========================================="
