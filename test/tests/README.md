# Test Files

This directory contains the main test files organized by functionality and scope. Tests are designed to validate CLI functionality, index operations, and various edge cases.

## Directory Structure

```
tests/
├── README.md                    # This file
├── setup_suite.bash            # Global test suite setup
├── _metatests/                 # Tests for test helpers themselves
├── global/                     # Global CLI functionality tests
└── index/                      # Index operation tests
    ├── basic_operations.bats
    ├── constraints_and_misusage.bats
    ├── deletion_protection_flow.bats
    ├── interactive_mode_flow.bats
    ├── serverless/             # Serverless index specific tests
    ├── pod/                    # Pod index tests (placeholder)
    └── integrated/             # Integrated index tests (placeholder)
```

## Test Categories

- **Global Tests (`global/`)** - verify basic CLI functionality and global features:
- **Index Tests (`index/`)** - generic index operation tests
- **Serverless Index Tests (`index/serverless/`)** - tests for serverless index configurations
- **Pod Index Tests (`index/pod/`)** - tests for pod index configurations
- **Integrated Index Tests (`index/pod/`)** - tests for integrated index configurations
- **Metatests (`_metatests/`)** - tests for testing infrastructure itself

## Test File Structure

### Standard Test File Layout

```bash
#!/usr/bin/env bats

# =============================================================================
# Load helpers
# -----------------------------------------------------------------------------
load "$BATS_ROOT/../helpers/custom/all.bash"
# =============================================================================

# =============================================================================
# Set tags for all tests in this file
# -----------------------------------------------------------------------------
# bats file_tags=scope:index, index-type:default, scenario:success
# =============================================================================

# -----------------------------------------------------------------------------
# Setup and teardown
# -----------------------------------------------------------------------------

setup_file() {
    # File-level setup (if needed)
}

teardown_file() {
    # File-level cleanup (if needed)
}

setup() {
    # Per-test setup
}

teardown() {
    # Per-test cleanup
}

# -----------------------------------------------------------------------------
# Tests
# -----------------------------------------------------------------------------

# bats test_tags=action:CRD, mode:params, index-type:default
@test "Test description" {
    # Test implementation
}
```

### Required Components

1. **Helper loading** - Load all necessary helper functions
2. **File tags** - Set common tags for all tests in the file
3. **Setup/teardown** - Handle test isolation and cleanup
4. **Test functions** - Individual test cases with descriptive names
5. **Test tags** - Specific tags for each test case

## Test Tagging Strategy

### File-Level Tags

Set common characteristics for all tests in a file:

```bash
# bats file_tags=scope:index, index-type:serverless, scenario:success
```

### Test-Level Tags

Set specific characteristics for individual tests:

```bash
# bats test_tags=action:create, mode:params, cloud:aws
@test "Create AWS serverless index" {
```

### Tag Categories

- **`scope:*`** - Test scope (global, index, collection)
- **`index-type:*`** - Index type (serverless, pod, integrated)
- **`vector-type:*`** - Vector type (dense, sparse)
- **`action:*`** - Operation type (create, read, update, delete, CRD)
- **`mode:*`** - Execution mode (params, interactive, api)
- **`cloud:*`** - Cloud provider (aws, gcp, azure, default)
- **`scenario:*`** - Test scenario (success, error, edge-case)

## Test Development Guidelines

### Adding New Tests

1. **Choose appropriate directory** - Place tests in the most relevant category
2. **Follow naming conventions** - Use descriptive, lowercase names with underscores
3. **Implement proper setup/teardown** - Ensure test isolation and cleanup
4. **Add appropriate tags** - Use both file-level and test-level tags
5. **Write descriptive test names** - Make test purpose clear from the name

### Test File Naming

- **`{feature}_operations.bats`** - Basic operations for a feature
- **`{feature}_constraints.bats`** - Validation and constraint testing
- **`{feature}_edge_cases.bats`** - Boundary conditions and unusual scenarios
- **`{feature}_flow.bats`** - Multi-step workflows and integration

### Test Function Naming

- **Use descriptive names** - `@test "Create serverless index with custom tags"`
- **Include key parameters** - `@test "Create GCP index with euclidean metric"`
- **Describe expected outcome** - `@test "Index creation fails with invalid dimension"`
