#!/bin/bash

# Integration test runner for Pinecone CLI
# This script sets up the environment and runs BATS tests

set -euo pipefail

# Signal handling for Ctrl+C
trap 'echo -e "\n[INFO] Received interrupt signal. Stopping tests..."; exit 130' INT TERM

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Default values
BATS_BINARY="bats"
TEST_FILE="index_flow_test.bats"
VERBOSE=false
PARALLEL=false
FILTER=""
SKIP_SETUP=false

# Function to print usage
print_usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Integration test runner for Pinecone CLI

OPTIONS:
    -h, --help              Show this help message
    -v, --verbose           Enable verbose output
    -p, --parallel          Run tests in parallel (if supported)
    -f, --filter PATTERN    Filter tests by pattern
    -s, --skip-setup        Skip environment setup
    -b, --bats PATH         Path to BATS binary (default: bats)
    -t, --test-file FILE    Test file to run (default: index_flow_test.bats)

ENVIRONMENT VARIABLES:
    PC_BINARY               Path to the Pinecone CLI binary (default: pcdev)
    BATS_BINARY            Path to BATS binary (default: bats)
    SKIP_LOGIN             Skip login check (default: false)

EXAMPLES:
    # Run all tests
    $0

    # Run with verbose output
    $0 -v

    # Run only tests matching "serverless"
    $0 -f "serverless"

    # Use custom BATS binary
    $0 -b /usr/local/bin/bats

    # Skip setup (for CI/CD)
    $0 -s

EOF
}

# Function to print colored output
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

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check dependencies
check_dependencies() {
    print_status "INFO" "Checking dependencies..."
    
    local missing_deps=()
    
    # Check for BATS
    if ! command_exists "$BATS_BINARY"; then
        missing_deps+=("BATS ($BATS_BINARY)")
    fi
    
    # Check for jq
    if ! command_exists "jq"; then
        missing_deps+=("jq")
    fi
    
    # Check for timeout
    if ! command_exists "timeout"; then
        missing_deps+=("timeout")
    fi
    
    if [ ${#missing_deps[@]} -gt 0 ]; then
        print_status "ERROR" "Missing dependencies: ${missing_deps[*]}"
        print_status "INFO" "Please install the missing dependencies:"
        echo "  - BATS: https://github.com/bats-core/bats-core"
        echo "  - jq: https://stedolan.github.io/jq/"
        echo "  - timeout: Usually included with coreutils"
        exit 1
    fi
    
    print_status "SUCCESS" "All dependencies found"
}

# Function to setup environment
setup_environment() {
    if [ "$SKIP_SETUP" = true ]; then
        print_status "INFO" "Skipping environment setup"
        return 0
    fi
    
    print_status "INFO" "Setting up test environment..."
    
    # Set default PC_BINARY if not set
    if [ -z "${PC_BINARY:-}" ]; then
        PC_BINARY="$PROJECT_ROOT/pcdev"
        export PC_BINARY
    fi
    
    # Check if PC_BINARY exists
    if [ ! -f "$PC_BINARY" ]; then
        print_status "ERROR" "PC_BINARY not found: $PC_BINARY"
        print_status "INFO" "Please build the CLI first: make build"
        exit 1
    fi
    
    # Make PC_BINARY executable
    chmod +x "$PC_BINARY"
    
    print_status "SUCCESS" "Using PC_BINARY: $PC_BINARY"
    
    # Check if logged in (unless SKIP_LOGIN is set)
    if [ "${SKIP_LOGIN:-}" != "true" ]; then
        print_status "INFO" "Checking login status..."
        if ! "$PC_BINARY" whoami >/dev/null 2>&1; then
            print_status "WARNING" "Not logged in to Pinecone"
            print_status "INFO" "Please run: $PC_BINARY login"
            print_status "INFO" "Or set SKIP_LOGIN=true to skip this check"
            exit 1
        fi
        print_status "SUCCESS" "Logged in to Pinecone"
    else
        print_status "INFO" "Skipping login check"
    fi
}

# Function to run tests
run_tests() {
    print_status "INFO" "Running tests..."
    
    local test_path="$SCRIPT_DIR/$TEST_FILE"
    
    if [ ! -f "$test_path" ]; then
        print_status "ERROR" "Test file not found: $test_path"
        exit 1
    fi
    
    # Build BATS command
    local bats_cmd=("$BATS_BINARY")
    
    if [ "$VERBOSE" = true ]; then
        bats_cmd+=("--verbose-run")
    fi
    
    if [ "$PARALLEL" = true ]; then
        bats_cmd+=("--jobs" "4")
    fi
    
    if [ -n "$FILTER" ]; then
        bats_cmd+=("--filter" "$FILTER")
    fi
    
    bats_cmd+=("$test_path")
    
    print_status "INFO" "Running: ${bats_cmd[*]}"
    
    # Run tests with proper signal handling
    # Only use timeout in non-interactive CI/CD environments
    if [ "${SKIP_LOGIN:-}" = "true" ] && [ -t 0 ]; then
        # Interactive mode - no timeout, allow Ctrl+C
        print_status "INFO" "Running in interactive mode - Ctrl+C will stop tests"
        if "${bats_cmd[@]}"; then
            print_status "SUCCESS" "All tests passed!"
            return 0
        else
            print_status "ERROR" "Some tests failed"
            return 1
        fi
    elif [ "${SKIP_LOGIN:-}" = "true" ]; then
        # Non-interactive CI/CD mode - use timeout
        print_status "INFO" "Using 10-minute timeout for CI/CD environment"
        if timeout 10m "${bats_cmd[@]}"; then
            print_status "SUCCESS" "All tests passed!"
            return 0
        else
            local exit_code=$?
            if [ $exit_code -eq 124 ]; then
                print_status "ERROR" "Tests timed out after 10 minutes"
            else
                print_status "ERROR" "Some tests failed"
            fi
            return 1
        fi
    else
        # Normal mode - no timeout
        if "${bats_cmd[@]}"; then
            print_status "SUCCESS" "All tests passed!"
            return 0
        else
            print_status "ERROR" "Some tests failed"
            return 1
        fi
    fi
}

# Function to cleanup
cleanup() {
    print_status "INFO" "Cleaning up..."
    
    # Cleanup any test indexes that might have been left behind
    if [ -n "${PC_BINARY:-}" ] && [ -f "$PC_BINARY" ]; then
        local test_indexes
        test_indexes=$("$PC_BINARY" index list --json 2>/dev/null | jq -r '.[] | select(.name | startswith("test-index-")) | .name' 2>/dev/null || true)
        
        if [ -n "$test_indexes" ]; then
            print_status "WARNING" "Found test indexes that weren't cleaned up:"
            echo "$test_indexes" | while read -r index_name; do
                if [ -n "$index_name" ]; then
                    print_status "INFO" "Deleting: $index_name"
                    "$PC_BINARY" index delete "$index_name" --yes >/dev/null 2>&1 || true
                fi
            done
        fi
    fi
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            print_usage
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -p|--parallel)
            PARALLEL=true
            shift
            ;;
        -f|--filter)
            FILTER="$2"
            shift 2
            ;;
        -s|--skip-setup)
            SKIP_SETUP=true
            shift
            ;;
        -b|--bats)
            BATS_BINARY="$2"
            shift 2
            ;;
        -t|--test-file)
            TEST_FILE="$2"
            shift 2
            ;;
        *)
            print_status "ERROR" "Unknown option: $1"
            print_usage
            exit 1
            ;;
    esac
done

# Set up trap for cleanup
trap cleanup EXIT

# Main execution
main() {
    print_status "INFO" "Starting integration tests..."
    print_status "INFO" "Project root: $PROJECT_ROOT"
    print_status "INFO" "Script directory: $SCRIPT_DIR"
    
    check_dependencies
    setup_environment
    run_tests
}

# Run main function
main "$@" 