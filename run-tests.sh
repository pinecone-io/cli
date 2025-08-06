#!/bin/bash

# Simple test runner for Pinecone CLI
# Usage: ./run-tests.sh [options]

set -e

# Signal handling for Ctrl+C
trap 'echo -e "\n[INFO] Received interrupt signal. Stopping tests..."; exit 130' INT TERM

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

# Check if CLI exists, build if needed
if [ ! -f "./pcdev" ]; then
    print_status "ERROR" "pcdev script not found. Please ensure you're in the project root."
    exit 1
fi

# Check if integration tests exist
if [ ! -f "./tests/integration/run_tests.sh" ]; then
    print_status "ERROR" "Integration tests not found. Please run from project root."
    exit 1
fi

# Make test runner executable
chmod +x ./tests/integration/run_tests.sh

# Set PC_BINARY to use the pcdev wrapper
export PC_BINARY="./pcdev"

# Pass through SKIP_LOGIN environment variable
if [ "${SKIP_LOGIN:-}" = "true" ]; then
    export SKIP_LOGIN=true
    print_status "INFO" "SKIP_LOGIN=true - skipping login checks"
fi

# Run the tests with all arguments passed through
print_status "INFO" "Running integration tests..."
./tests/integration/run_tests.sh "$@" 