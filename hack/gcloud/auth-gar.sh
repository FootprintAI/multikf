#!/bin/bash
# Copyright 2024 FootprintAI
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

# Configuration
GAR_REGION="${GAR_REGION:-asia-east1}"
GAR_PROJECT_ID="${GAR_PROJECT_ID:-footprintai-dev}"
SERVICE_ACCOUNT_FILE="${SERVICE_ACCOUNT_FILE:-$HOME/.config/gcloud/service-account-key.json}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if gcloud is installed
check_gcloud_installed() {
    if command -v gcloud &> /dev/null; then
        return 0
    else
        return 1
    fi
}

# Function to install gcloud CLI
install_gcloud() {
    print_info "Installing Google Cloud SDK..."

    # Detect OS
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux installation
        print_info "Detected Linux OS"

        # Install prerequisites
        if command -v apt-get &> /dev/null; then
            print_info "Installing prerequisites..."
            sudo apt-get update
            sudo apt-get install -y apt-transport-https ca-certificates gnupg curl

            # Add Google Cloud SDK repo
            print_info "Adding Google Cloud SDK repository..."
            echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | \
                sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list

            # Import Google Cloud public key
            curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | \
                sudo gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg

            # Install gcloud
            print_info "Installing gcloud CLI..."
            sudo apt-get update
            sudo apt-get install -y google-cloud-cli

        elif command -v yum &> /dev/null; then
            print_info "Installing via yum..."
            sudo tee -a /etc/yum.repos.d/google-cloud-sdk.repo << EOM
[google-cloud-cli]
name=Google Cloud CLI
baseurl=https://packages.cloud.google.com/yum/repos/cloud-sdk-el9-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=0
gpgkey=https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOM
            sudo yum install -y google-cloud-cli
        else
            print_error "Package manager not supported. Please install gcloud manually from:"
            print_error "https://cloud.google.com/sdk/docs/install"
            exit 1
        fi

    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS installation
        print_info "Detected macOS"

        if command -v brew &> /dev/null; then
            print_info "Installing via Homebrew..."
            brew install --cask google-cloud-sdk
        else
            print_error "Homebrew not found. Please install Homebrew or install gcloud manually from:"
            print_error "https://cloud.google.com/sdk/docs/install"
            exit 1
        fi
    else
        print_error "OS not supported. Please install gcloud manually from:"
        print_error "https://cloud.google.com/sdk/docs/install"
        exit 1
    fi

    print_info "Google Cloud SDK installed successfully"
}

# Function to authenticate with service account
authenticate_with_service_account() {
    local sa_file="$1"

    print_info "Authenticating with service account: $sa_file"

    if [[ ! -f "$sa_file" ]]; then
        print_error "Service account file not found: $sa_file"
        print_error "Please provide the path to your service account JSON file"
        print_error "Example: export SERVICE_ACCOUNT_FILE=/path/to/service-account-key.json"
        exit 1
    fi

    # Activate service account
    gcloud auth activate-service-account --key-file="$sa_file"

    # Set project (extract from service account file if possible)
    if [[ -n "$GAR_PROJECT_ID" ]]; then
        print_info "Setting project to: $GAR_PROJECT_ID"
        gcloud config set project "$GAR_PROJECT_ID"
    fi

    print_info "Service account authenticated successfully"
}

# Function to configure Docker for GAR
configure_docker_gar() {
    print_info "Configuring Docker authentication for Google Artifact Registry..."

    local gar_endpoint="${GAR_REGION}-docker.pkg.dev"

    gcloud auth configure-docker "$gar_endpoint" --quiet

    print_info "Docker configured for GAR endpoint: $gar_endpoint"
}

# Main script
main() {
    print_info "Starting Google Cloud SDK setup and GAR authentication..."
    echo ""

    # Check if gcloud is already installed
    if check_gcloud_installed; then
        print_info "Google Cloud SDK is already installed"
        gcloud version
    else
        print_warn "Google Cloud SDK not found"
        install_gcloud
    fi

    echo ""

    # Authenticate with service account
    authenticate_with_service_account "$SERVICE_ACCOUNT_FILE"

    echo ""

    # Configure Docker for GAR
    configure_docker_gar

    echo ""
    print_info "✓ Setup completed successfully!"
    print_info "You can now use Docker with Google Artifact Registry"
    print_info ""
    print_info "Example usage:"
    print_info "  docker pull ${GAR_REGION}-docker.pkg.dev/${GAR_PROJECT_ID}/cockburn/your-image:tag"
    print_info "  docker push ${GAR_REGION}-docker.pkg.dev/${GAR_PROJECT_ID}/cockburn/your-image:tag"
}

# Run main function
main
