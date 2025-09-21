echo "=== Removing any existing CUDA / NVIDIA drivers ==="

# Stop running services that might hold GPU
sudo systemctl stop nvidia-persistenced || true

# Purge existing NVIDIA & CUDA packages
sudo apt-get remove --purge -y \
    'cuda*' \
    'nvidia*' \
    'libnvidia*' \
    'nsight*' \
    'libcudnn*' || true

# Clean up apt cache
sudo apt-get autoremove -y
sudo apt-get clean

# Remove any leftover directories and repos
sudo rm -rf /usr/local/cuda*
sudo rm -rf /var/lib/cuda-repo*
sudo rm -f /etc/apt/sources.list.d/cuda*.list
sudo rm -f /etc/apt/sources.list.d/nvidia*.list
sudo rm -f /usr/share/keyrings/cuda-*.gpg

echo "=== Old CUDA / NVIDIA installation removed ==="
