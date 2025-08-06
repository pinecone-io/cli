# Pinecone CLI Integration Tests

This directory contains comprehensive integration tests for the Pinecone CLI, specifically focusing on the index flow: create → describe → delete.

## Overview

The integration tests use **BATS (Bash Automated Testing System)** to test all possible combinations of flags for the `index create` command, along with invalid and missing values. The tests are designed to be automated as much as possible using loops, variables, and helper functions.

## Test Coverage

### Index Types Tested

- **Serverless Indexes**: Default and all flag combinations
- **Pod Indexes**: All pod types, shards, replicas combinations
- **Integrated Indexes**: All embedding models and configurations

### Flag Combinations Tested

- All possible flag combinations for each index type
- Invalid flag combinations (e.g., serverless flags with pod indexes)
- Missing required values
- Invalid values (negative numbers, invalid strings, etc.)

### Test Categories

1. **Basic Functionality**: Minimal flag combinations
2. **Full Configuration**: All possible flags for each index type
3. **Parameter Variations**: Different values for each parameter
4. **Error Handling**: Invalid inputs and edge cases
5. **Lifecycle Testing**: Index states and transitions
6. **Concurrent Operations**: Multiple indexes created simultaneously
7. **Performance**: Large configurations and edge cases

## Prerequisites

### Required Dependencies

- **BATS**: Bash Automated Testing System
- **jq**: JSON processor for parsing CLI output
- **timeout**: For handling interactive mode tests
- **goreleaser**: For building the CLI (optional, will be installed if missing)
- **Pinecone CLI**: Built and available via `pcdev` script

### Installing Dependencies

#### Quick Install (Recommended)

```bash
# Install all dependencies automatically
./install-deps.sh
```

#### Manual Installation

##### macOS (using Homebrew)

```bash
# Install BATS
brew install bats-core

# Install jq (if not already installed)
brew install jq
```

##### Ubuntu/Debian

```bash
# Install BATS
sudo apt-get update
sudo apt-get install bats

# Install jq
sudo apt-get install jq
```

##### CentOS/RHEL

```bash
# Install BATS
sudo yum install bats

# Install jq
sudo yum install jq
```

##### Manual BATS Installation

```bash
# Clone BATS repository
git clone https://github.com/bats-core/bats-core.git
cd bats-core
sudo ./install.sh /usr/local
```

## Setup

### 1. Build the CLI

```bash
# From the project root
./pcdev
```

The project uses goreleaser for building. The `pcdev` script automatically:

- Builds the CLI using goreleaser
- Detects your OS and architecture
- Runs the correct binary from the `dist` directory

### 2. Login to Pinecone

```bash
./pcdev login
```

### 3. Verify Setup

```bash
./pcdev whoami
```

## Running Tests

### Quick Start

```bash
# Install dependencies (first time only)
./install-deps.sh

# Run all tests
./run-tests.sh

# Run with verbose output
./run-tests.sh -v

# Run only tests matching a pattern
./run-tests.sh -f "serverless"
```

### Advanced Usage

```bash
# Run tests in parallel
./run-tests.sh -p

# Skip environment setup (for CI/CD)
./run-tests.sh -s

# Use custom BATS binary
BATS_BINARY=/usr/local/bin/bats ./run-tests.sh

# Skip login check (for CI/CD)
SKIP_LOGIN=true ./run-tests.sh

**Note**: The `SKIP_LOGIN=true` environment variable is automatically passed through to the BATS tests, so you can use it with the main test runner to skip login checks in CI/CD environments.
```

### Direct BATS Usage

```bash
# Run specific test file
bats tests/integration/index_flow_test.bats

# Run with verbose output
bats --verbose tests/integration/index_flow_test.bats

# Run only tests matching pattern
bats --filter "serverless" tests/integration/index_flow_test.bats
```

## Test Structure

### Test File: `index_flow_test.bats`

The main test file contains:

#### Setup and Teardown

- **setup()**: Initializes test environment and checks login status
- **teardown()**: Cleans up any test indexes that weren't properly deleted

#### Helper Functions

- **generate_index_name()**: Creates unique index names with timestamps
- **wait_for_index_ready()**: Waits for index to reach "Ready" state
- **extract_json_field()**: Extracts specific fields from JSON output

#### Test Categories

##### 1. Basic Index Creation

```bash
@test "create serverless index with minimal flags"
@test "create pod index with minimal flags"
@test "create integrated index with minimal flags"
```

##### 2. Full Configuration Tests

```bash
@test "create serverless index with all flags"
@test "create pod index with all flags"
@test "create integrated index with all flags"
```

##### 3. Parameter Variation Tests

```bash
@test "create serverless index with different vector types"
@test "create indexes with different metrics"
@test "create serverless indexes with different cloud providers"
@test "create indexes with different dimensions"
@test "create pod indexes with different pod types"
@test "create pod indexes with different shards and replicas"
@test "create integrated indexes with different models"
```

##### 4. Feature Tests

```bash
@test "create indexes with deletion protection"
@test "create indexes with tags"
@test "create indexes with source collection"
```

##### 5. Error Handling Tests

```bash
@test "create index with invalid flag combinations"
@test "create index with invalid values"
@test "create index with missing required values"
@test "describe non-existent index"
@test "delete non-existent index"
```

##### 6. Advanced Tests

```bash
@test "create and describe index with JSON output"
@test "create index in interactive mode"
@test "create index without confirmation"
@test "create multiple indexes concurrently"
@test "test index lifecycle states"
@test "handle network errors gracefully"
@test "create index with large configuration"
@test "test edge cases"
```

## Test Automation Features

### Loops and Arrays

The tests use bash arrays and loops to test multiple values efficiently:

```bash
# Test all vector types
local vector_types=("dense" "sparse")
for vector_type in "${vector_types[@]}"; do
    # Test logic here
done

# Test all metrics
local metrics=("cosine" "euclidean" "dotproduct")
for metric in "${metrics[@]}"; do
    # Test logic here
done
```

### Dynamic Test Generation

Tests are generated dynamically based on parameter combinations:

```bash
# Test different shards and replicas combinations
local shards_replicas=("1:1" "2:1" "1:2" "2:2" "4:2" "2:4")
for combo in "${shards_replicas[@]}"; do
    IFS=':' read -r shards replicas <<< "$combo"
    # Test logic here
done
```

### Helper Functions

Reusable functions reduce code duplication:

```bash
# Generate unique index names
generate_index_name() {
    local suffix="${1:-}"
    echo "${TEST_INDEX_PREFIX}${suffix}"
}

# Wait for index to be ready
wait_for_index_ready() {
    local index_name="$1"
    # Wait logic here
}
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Integration Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Install dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y bats jq
      - name: Build CLI
        run: ./pcdev
      - name: Run integration tests
        run: |
          export SKIP_LOGIN=true
          ./tests/integration/run_tests.sh -s
```

### Local Development

```bash
# Run tests during development
./run-tests.sh

# Run specific test categories
./run-tests.sh -f "serverless"
./run-tests.sh -f "pod"
./run-tests.sh -f "integrated"
```

## Troubleshooting

### Common Issues

#### 1. BATS not found

```bash
# Install BATS manually
git clone https://github.com/bats-core/bats-core.git
cd bats-core
sudo ./install.sh /usr/local
```

#### 2. jq not found

```bash
# macOS
brew install jq

# Ubuntu/Debian
sudo apt-get install jq

# CentOS/RHEL
sudo yum install jq
```

#### 3. CLI binary not found

```bash
# Build the CLI first
./pcdev
```

#### 4. Not logged in

```bash
# Login to Pinecone
./pcdev login

# Or skip login check for CI/CD
export SKIP_LOGIN=true
```

#### 5. Tests timing out

```bash
# Increase timeout in wait_for_index_ready function
# Default is 30 attempts with 10-second intervals
```

### Debug Mode

```bash
# Run with verbose output
./run-tests.sh -v

# Run single test
bats --verbose tests/integration/index_flow_test.bats -f "create serverless index"

# Debug specific test
bats --verbose --tap tests/integration/index_flow_test.bats
```

## Test Results

### Expected Output

```
 ✓ create serverless index with minimal flags
 ✓ create serverless index with all flags
 ✓ create pod index with minimal flags
 ✓ create pod index with all flags
 ✓ create integrated index with minimal flags
 ✓ create integrated index with all flags
 ...

X tests, 0 failures
```

### Test Statistics

- **Total Tests**: 50+ test cases
- **Coverage**: All index types, all flag combinations, error cases
- **Automation**: 95% automated using loops and helper functions
- **Runtime**: ~30-60 minutes (depending on index creation time)

## Contributing

### Adding New Tests

1. Add new test functions to `index_flow_test.bats`
2. Follow the naming convention: `@test "descriptive test name"`
3. Use helper functions for common operations
4. Add appropriate cleanup in teardown

### Test Guidelines

- Use descriptive test names
- Test both success and failure cases
- Clean up resources in teardown
- Use helper functions for common operations
- Add comments for complex test logic

### Running Tests Locally

```bash
# Quick test run
./run-tests.sh -f "minimal"

# Full test suite
./run-tests.sh -v

# Debug specific test
bats --verbose tests/integration/index_flow_test.bats -f "test name"
```

## Performance Considerations

### Test Optimization

- Tests run in parallel where possible
- Index creation is the bottleneck (5-15 minutes per index)
- Use `--yes` flag to skip confirmations
- Clean up resources promptly

### Resource Management

- Tests create temporary indexes with unique names
- Automatic cleanup in teardown function
- Timeout handling for long-running operations
- Error handling for network issues

## Security Notes

### API Key Management

- Tests use the logged-in user's credentials
- No hardcoded API keys in test files
- Environment variables for sensitive data
- Secure cleanup of test resources

### Test Isolation

- Each test uses unique index names
- No interference between test runs
- Proper cleanup even on test failure
- Isolated test environment
