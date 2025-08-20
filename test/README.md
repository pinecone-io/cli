# Test Suite Documentation

## Overview

This directory contains a comprehensive test suite for the Pinecone CLI using BATS (Bash Automated Testing System). The tests cover CLI functionality, index operations, and various edge cases.

## Installation

The test suite is using BATS testing framework. Both the framework and its helpers
are expected to be installed in the `test` directory as Git submodules.

The framework also depends on the following tools:

- `jq` for JSON processing
- `jd` for JSON validation
- `parallel` for parallel test execution

The BATS framework and the dependencies can be installed by running the following command:

```bash
# From project root
test/install-deps.sh

# From test directory
./install-deps.sh
```

## Test Framework Structure

```
test/
├── README.md                  # This file - main documentation
├── install-deps.sh            # Test dependencies installer
├── bats                       # BATS executable
├── bats-core/                 # BATS core library
├── helpers/                   # Test helper functions and libraries
│   ├── bats-assert/           # BATS assertion library
│   ├── bats-support/          # BATS support library
│   ├── custom/                # Custom helper functions
│   ├── templates/             # JSON templates for validation
│   └── bin/                   # Utility scripts
├── tdd/                       # Test-driven development tests
└── tests/                     # Main test files
    ├── _metatests/            # Tests for test helpers themselves
    ├── global/                # Global CLI functionality tests
    └── index/                 # Index operation tests
        ├── serverless/        # Serverless index specific tests
        ├── pod/               # Pod index tests (placeholder)
        └── integrated/        # Integrated index tests (placeholder)
```

## Running Tests

**From project root:**

```bash
test/bats test/tests/
test/bats test/tests/index/
test/bats test/tests/index/serverless/
```

**From test directory:**

```bash
cd test
./bats tests/
./bats tests/index/
./bats tests/index/serverless/
```

**With tags:**

```bash
# From project root
test/bats --filter-tags action:create test/tests/
test/bats --filter-tags index-type:serverless test/tests/

# From test directory
./bats --filter-tags action:create tests/
./bats --filter-tags index-type:serverless tests/
```

## Test Tagging System

We use a standardized tagging system to categorize and organize tests.

### Tag Categories

- **`action:*`** - Operation type (create, read, update, delete, CRD, CRUD)
- **`mode:*`** - Execution mode (params, interactive, api)
- **`index-type:*`** - Index type (default, serverless, pod, integrated)
- **`cloud:*`** - Cloud provider (aws, gcp, azure, default)
- **`validation:*`** - Validation aspect (flags, constraints, format)
- **`scenario:*`** - Test scenario (success, error, edge-case)

### Usage Examples

```bash
# Run all create operations
test/bats --filter-tags action:create test/tests/

# Run all serverless index tests
test/bats --filter-tags scope:index,index-type:serverless test/tests/

# Run serverless index tests on AWS
test/bats --filter-tags scope:index,index-type:serverless,cloud:aws test/tests/

# Ignore metatests
test/bats --filter-tags !scope:_metatests test/tests/

```

### Tagging Guidelines

1. **Use standard format**: `category:value`
2. **Use lowercase with hyphens** for multi-word values
3. **Apply multiple tags** when a test fits multiple categories
4. **Be consistent** with tag values across similar tests

## Test Organization

Tests are organized by functionality and scope:

- **`tests/global/`** - Basic CLI functionality (execution, authentication)
- **`tests/index/`** - Core index operations (CRD flows, validation, interactive mode)
- **`tests/index/serverless/`** - Serverless index specific tests (cloud providers, constraints, edge cases)
- **`tests/_metatests/`** - Tests for test helpers and infrastructure
- **`tests/index/pod/`** - Pod index tests (placeholder for future)
- **`tests/index/integrated/`** - Integrated index tests (placeholder for future)

Each directory contains relevant test files following the naming convention: `{feature}_{aspect}.bats`

## Custom Helper Functions

Helper functions are organized by functionality:

- **`helpers/custom/global.bash`** - CLI setup and global utilities
- **`helpers/custom/indexes.bash`** - Index-specific operations and validation
- **`helpers/bats-assert/`** - BATS assertion library
- **`helpers/bats-support/`** - BATS support library

Key functions include test data generation, index state management, JSON validation, and template handling. See `helpers/README.md` for detailed documentation.

## JSON Templates

Located in `helpers/templates/indexes/`, these templates define expected index response structures:

- **Cloud provider templates** - AWS, GCP, Azure configurations
- **Configuration variants** - Different metrics, vector types, dimensions
- **Feature templates** - Tags, deletion protection, source collections

See `helpers/templates/indexes/README.md` for detailed template documentation and usage examples.

## Adding New Tests

### Test File Structure

1. **Load helpers**: `load "$BATS_ROOT/../helpers/custom/all.bash"`
2. **Set file tags**: Use `bats file_tags` for common characteristics
3. **Setup/teardown**: Implement `setup()` and `teardown()` for test isolation
4. **Test functions**: Use `@test` with descriptive names and appropriate tags
5. **Cleanup**: Always clean up resources in `teardown()`
