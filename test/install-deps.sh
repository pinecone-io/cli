#!/bin/bash

# Install test dependencies for Pinecone CLI
# Usage: ./install-deps.sh (from test directory)
#        test/install-deps.sh (from project root)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print status messages
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

# Function to check if jq is available
check_jq() {
    if command -v jq >/dev/null 2>&1; then
        print_status "SUCCESS" "jq is already available: $(jq --version)"
        return 0
    else
        return 1
    fi
}

# Function to check if GNU parallel is available
check_parallel() {
    if command -v parallel >/dev/null 2>&1 && parallel --version >/dev/null 2>&1; then
        print_status "SUCCESS" "GNU parallel is already available: $(parallel --version | head -n1)"
        return 0
    else
        return 1
    fi
}

# Function to check if jd is available
check_jd() {
    if command -v jd >/dev/null 2>&1; then
        print_status "SUCCESS" "jd is already available: $(jd --version)"
        return 0
    else
        return 1
    fi
}

# Function to check if BATS submodules are initialized
check_bats_submodules() {
    # BATS submodules are always relative to the script location
    local bats_core_dir="$SCRIPT_DIR/bats-core"
    local bats_assert_dir="$SCRIPT_DIR/helpers/bats-assert"
    local bats_support_dir="$SCRIPT_DIR/helpers/bats-support"
    
    if [ -d "$bats_core_dir" ] && [ -f "$bats_core_dir/bin/bats" ] && \
       [ -d "$bats_assert_dir" ] && [ -f "$bats_assert_dir/load.bash" ] && \
       [ -d "$bats_support_dir" ] && [ -f "$bats_support_dir/load.bash" ]; then
        print_status "SUCCESS" "BATS submodules are already initialized"
        return 0
    else
        print_status "WARNING" "BATS submodules are not initialized"
        print_status "INFO" "Run 'git submodule update --init --recursive' to initialize them"
        return 1
    fi
}

print_status "INFO" "Installing test dependencies..."

# Determine the correct paths based on execution context (for use throughout the script)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Detect package manager once
PACKAGE_MANAGER=""
if command -v brew >/dev/null 2>&1; then
    PACKAGE_MANAGER="brew"
    print_status "INFO" "Detected Homebrew package manager"
elif command -v apt-get >/dev/null 2>&1; then
    PACKAGE_MANAGER="apt"
    print_status "INFO" "Detected apt-get package manager"
elif command -v yum >/dev/null 2>&1; then
    PACKAGE_MANAGER="yum"
    print_status "INFO" "Detected yum package manager"
else
    print_status "ERROR" "No supported package manager found"
    print_status "INFO" "Please install a package manager:"
    echo "  macOS: Install Homebrew (https://brew.sh/)"
    echo "  Ubuntu/Debian: apt-get is usually available"
    echo "  CentOS/RHEL: yum is usually available"
    exit 1
fi

# Install dependencies using the detected package manager
case "$PACKAGE_MANAGER" in
    "brew")
        print_status "INFO" "Installing dependencies via Homebrew..."
        
        # Check dependencies and install only what's missing
        if ! check_jq; then
            print_status "INFO" "Installing jq CLI..."
            brew install jq
            print_status "SUCCESS" "jq installed via Homebrew"
        fi
        
        if ! check_parallel; then
            print_status "INFO" "Installing GNU parallel..."
            brew install parallel
            print_status "SUCCESS" "GNU parallel installed via Homebrew"
        fi
        
        if ! check_jd; then
            print_status "INFO" "Installing jd..."
            brew install jd
            print_status "SUCCESS" "jd installed via Homebrew"
        fi
        ;;
    "apt")
        print_status "INFO" "Installing dependencies via apt-get..."
        sudo apt-get update
        
        # Check dependencies and install only what's missing
        if ! check_jq; then
            sudo apt-get install -y jq
            print_status "SUCCESS" "jq installed via apt-get"
        fi
        
        if ! check_parallel; then
            sudo apt-get install -y parallel
            print_status "SUCCESS" "GNU parallel installed via apt-get"
        fi
        
        if ! check_jd; then
            sudo apt-get install -y jd
            print_status "SUCCESS" "jd installed via apt-get"
        fi
        ;;
    "yum")
        print_status "INFO" "Installing dependencies via yum..."
        
        # Check dependencies and install only what's missing
        if ! check_jq; then
            sudo yum install -y jq
            print_status "SUCCESS" "jq installed via yum"
        fi
        
        if ! check_parallel; then
            sudo yum install -y parallel
            print_status "SUCCESS" "GNU parallel installed via yum"
        fi
        
        if ! check_jd; then
            sudo yum install -y jd
            print_status "SUCCESS" "jd installed via yum"
        fi
        ;;
esac

# Check and initialize BATS submodules (same for all package managers)
print_status "INFO" "Checking BATS submodules..."
if ! check_bats_submodules; then
    print_status "INFO" "Initializing BATS submodules..."
    if git submodule update --init --recursive; then
        print_status "SUCCESS" "BATS submodules initialized successfully"
    else
        print_status "ERROR" "Failed to initialize BATS submodules"
        print_status "INFO" "Please run manually: git submodule update --init --recursive"
        exit 1
    fi
fi

# Verify installation
print_status "INFO" "Verifying installation..."

# Note: Dependencies were already checked during installation
# Only verify that they're actually working
if ! command -v jq >/dev/null 2>&1; then
    print_status "ERROR" "jq installation failed"
    exit 1
fi

if ! command -v parallel >/dev/null 2>&1; then
    print_status "ERROR" "GNU parallel installation failed"
    exit 1
fi

if ! command -v jd >/dev/null 2>&1; then
    print_status "ERROR" "jd installation failed"
    exit 1
fi

print_status "SUCCESS" "All test dependencies and BATS submodules are ready!"

# Show test commands relative to current working directory
if [ "$(pwd)" = "$SCRIPT_DIR" ]; then
    # User is in the test directory
    print_status "INFO" "You can now run tests: ./bats tests/"
    print_status "INFO" "For parallel test execution, use: ./bats --jobs N tests/"
else
    # User is in the project root or elsewhere
    print_status "INFO" "You can now run tests: test/bats test/tests/"
    print_status "INFO" "For parallel test execution, use: test/bats --jobs N test/tests/"
fi 