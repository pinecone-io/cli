#!/bin/bash

# Install development dependencies for Pinecone CLI
# Usage: ./install-deps.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    local level="$1"
    shift
    case "$level" in
        "INFO")
            echo -e "${BLUE}[INFO]${NC} $*"
            ;;
        "SUCCESS")
            echo -e "${GREEN}[SUCCESS]${NC} $*"
            ;;
        "WARNING")
            echo -e "${YELLOW}[WARNING]${NC} $*"
            ;;
        "ERROR")
            echo -e "${RED}[ERROR]${NC} $*"
            ;;
    esac
}

print_status "INFO" "Installing development dependencies..."

# Check which package manager is available
if command -v brew >/dev/null 2>&1; then
    print_status "INFO" "Using Homebrew..."
    brew install bats-core jq
    print_status "SUCCESS" "Dependencies installed via Homebrew"
elif command -v apt-get >/dev/null 2>&1; then
    print_status "INFO" "Using apt-get..."
    sudo apt-get update
    sudo apt-get install -y bats jq
    print_status "SUCCESS" "Dependencies installed via apt-get"
elif command -v yum >/dev/null 2>&1; then
    print_status "INFO" "Using yum..."
    sudo yum install -y bats jq
    print_status "SUCCESS" "Dependencies installed via yum"
else
    print_status "WARNING" "No supported package manager found"
    print_status "INFO" "Please install dependencies manually:"
    echo "  BATS: https://github.com/bats-core/bats-core"
    echo "  jq: https://stedolan.github.io/jq/"
    echo ""
    print_status "INFO" "Or install a package manager:"
    echo "  macOS: Install Homebrew (https://brew.sh/)"
    echo "  Ubuntu/Debian: apt-get is usually available"
    echo "  CentOS/RHEL: yum is usually available"
    exit 1
fi

# Verify installation
print_status "INFO" "Verifying installation..."
if command -v bats >/dev/null 2>&1; then
    print_status "SUCCESS" "BATS installed: $(bats --version)"
else
    print_status "ERROR" "BATS installation failed"
    exit 1
fi

if command -v jq >/dev/null 2>&1; then
    print_status "SUCCESS" "jq installed: $(jq --version)"
else
    print_status "ERROR" "jq installation failed"
    exit 1
fi

print_status "SUCCESS" "All dependencies installed successfully!"

# Check if goreleaser is available
if ! command -v goreleaser >/dev/null 2>&1; then
    print_status "WARNING" "goreleaser not found. Install it with:"
    echo "  go install github.com/goreleaser/goreleaser/v2/cmd/goreleaser@latest"
    echo "  or visit: https://goreleaser.com/install/"
fi

print_status "INFO" "You can now run: ./run-tests.sh" 