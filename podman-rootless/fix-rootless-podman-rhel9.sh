#!/bin/bash
# Fix script for rootless podman + kind on RHEL 9.4
# Run this script to configure your system for proper systemd support in rootless containers

set -e

echo "=========================================="
echo "Fixing Rootless Podman for Kind on RHEL 9.4"
echo "=========================================="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "This script needs to be run with sudo for system configuration"
    echo "Please run: sudo bash $0"
    exit 1
fi

# Get the actual user (not root)
ACTUAL_USER="${SUDO_USER:-$USER}"
ACTUAL_UID=$(id -u "$ACTUAL_USER")

echo "Configuring for user: $ACTUAL_USER (UID: $ACTUAL_UID)"
echo ""

# Step 1: Configure cgroup v2 delegation
echo "Step 1: Configuring cgroup v2 delegation..."
mkdir -p /etc/systemd/system/user@.service.d

cat > /etc/systemd/system/user@.service.d/delegate.conf << 'EOF'
[Service]
Delegate=cpu cpuset io memory pids
EOF

echo "✓ Created /etc/systemd/system/user@.service.d/delegate.conf"
echo ""

# Step 2: Enable lingering for the user
echo "Step 2: Enabling user lingering (keeps user services running)..."
loginctl enable-linger "$ACTUAL_USER"
echo "✓ Lingering enabled for $ACTUAL_USER"
echo ""

# Step 3: Reload systemd
echo "Step 3: Reloading systemd configuration..."
systemctl daemon-reexec
systemctl daemon-reload
echo "✓ Systemd reloaded"
echo ""

# Step 4: Restart user service
echo "Step 4: Restarting user service..."
systemctl restart "user@${ACTUAL_UID}.service"
echo "✓ User service restarted"
echo ""

# Step 5: Configure podman for systemd
echo "Step 5: Configuring podman for systemd support..."

# Create containers.conf for the user if it doesn't exist
USER_CONTAINERS_CONF="/home/$ACTUAL_USER/.config/containers/containers.conf"
mkdir -p "$(dirname "$USER_CONTAINERS_CONF")"

# Check if file exists and has systemd config
if ! grep -q "cgroup_manager" "$USER_CONTAINERS_CONF" 2>/dev/null; then
    cat >> "$USER_CONTAINERS_CONF" << 'EOF'

[engine]
cgroup_manager = "systemd"
events_logger = "journald"
EOF
    chown "$ACTUAL_USER:$ACTUAL_USER" "$USER_CONTAINERS_CONF"
    echo "✓ Updated podman configuration at $USER_CONTAINERS_CONF"
else
    echo "✓ Podman configuration already exists"
fi
echo ""

# Step 6: Verify delegation
echo "Step 6: Verifying cgroup delegation..."
sleep 2  # Wait for systemd to apply changes
CONTROLLERS=$(cat /sys/fs/cgroup/user.slice/user-${ACTUAL_UID}.slice/user@${ACTUAL_UID}.service/cgroup.controllers 2>/dev/null || echo "ERROR: Cannot read controllers")
echo "Available controllers: $CONTROLLERS"

if echo "$CONTROLLERS" | grep -q "cpu" && echo "$CONTROLLERS" | grep -q "memory" && echo "$CONTROLLERS" | grep -q "pids"; then
    echo "✓ Required controllers are available"
else
    echo "⚠ Warning: Not all required controllers are available"
    echo "  You may need to reboot the system for changes to take effect"
fi
echo ""

echo "=========================================="
echo "Configuration Complete!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Switch back to your user account: exit"
echo "2. Verify the fix: ./podman-rootless/diagnose-rootless-podman.sh"
echo "3. If controllers are not available, reboot the system: sudo reboot"
echo "4. After reboot, try creating the cluster again:"
echo "   KIND_EXPERIMENTAL_PROVIDER=podman kind create cluster --name gpu-cluster4 --config ./podman-rootless/gpu-kind-config.yaml --image $KIND_NODE_IMAGE"
echo ""
