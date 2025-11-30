#!/usr/bin/env bash

# Usage: ./setup-huggingface.sh [TOKEN]
# Example: ./setup-huggingface.sh hf_xxxxxxxxxxxxx
# If no token provided, will prompt for interactive login

set -e

VENV_DIR="${HOME}/.huggingface-venv"

echo "Setting up Hugging Face CLI..."

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo "Error: Python3 is not installed. Please install Python3 first."
    exit 1
fi

# Create virtual environment if it doesn't exist
if [[ ! -d "$VENV_DIR" ]]; then
    echo "Creating virtual environment at $VENV_DIR..."
    python3 -m venv "$VENV_DIR"
else
    echo "Virtual environment already exists at $VENV_DIR"
fi

# Activate virtual environment
echo "Activating virtual environment..."
source "$VENV_DIR/bin/activate"

# Upgrade pip
echo "Upgrading pip..."
pip install --upgrade pip

# Install huggingface-hub
echo "Installing huggingface-hub..."
pip install --upgrade huggingface-hub

# Verify installation
if ! command -v hf &> /dev/null; then
    echo "Error: hf CLI was not installed successfully."
    exit 1
fi

echo "huggingface-hub installed successfully!"

# Login to Hugging Face
TOKEN=${1:-""}

if [[ -n "$TOKEN" ]]; then
    echo "Logging in with provided token..."
    hf auth login --token "$TOKEN"
else
    echo "No token provided. Starting interactive login..."
    echo "You can get your token from: https://huggingface.co/settings/tokens"
    hf auth login
fi

echo ""
echo "Setup complete!"
echo ""
echo "To use Hugging Face CLI in the future, activate the virtual environment first:"
echo "  source $VENV_DIR/bin/activate"
echo ""
echo "Then verify your login with:"
echo "  hf auth whoami"
echo ""
echo "To add this to your shell profile (~/.bashrc or ~/.zshrc):"
echo "  alias hf-activate='source $VENV_DIR/bin/activate'"
