#!/usr/bin/env bats

# =============================================================================
# Load helpers
#
# -----------------------------------------------------------------------------
load "$BATS_ROOT/../helpers/custom/all.bash"
# =============================================================================


# =============================================================================
# Set tags for all tests in this file
#
# -----------------------------------------------------------------------------
# bats file_tags=scope:index, scenario:error
# =============================================================================


# -----------------------------------------------------------------------------
# Setup and teardown
# -----------------------------------------------------------------------------

# If something needs to be done before all tests in this file, do it here.
setup_file() {
    :
}

# If something needs to be done after all tests in this file, do it here.
teardown_file() {
    :
}

# If something needs to be done before each test in this file, do it here.
setup() {
    # Generate a unique index name for this test
    export TEST_INDEX_NAME=$(generate_index_name)
}

# If something needs to be done after each test in this file, do it here.
teardown() {
    # Attempt to clean up the index created by this test
    # This may fail if the test didn't actually create an index, which is fine
    if [ -n "$TEST_INDEX_NAME" ]; then
        $CLI index delete "$TEST_INDEX_NAME" 2>/dev/null || true
    fi
}

# -----------------------------------------------------------------------------
# Tests
# -----------------------------------------------------------------------------

# bats test_tags=action:create, validation:constraints
@test "Only one index type flag is allowed" {
    # Test all possible combinations of index type flags
    # Each should fail with the same error message
    
    # Test --pod and --integrated together
    run $CLI index create ${TEST_INDEX_NAME} --pod --integrated
    assert_failure
    assert_output --partial "only one index type can be specified"
    
    # Test --pod and --serverless together
    run $CLI index create ${TEST_INDEX_NAME} --pod --serverless
    assert_failure
    assert_output --partial "only one index type can be specified"
    
    # Test --integrated and --serverless together
    run $CLI index create ${TEST_INDEX_NAME} --integrated --serverless
    assert_failure
    assert_output --partial "only one index type can be specified"
    
    # Test all three flags together
    run $CLI index create ${TEST_INDEX_NAME} --pod --integrated --serverless
    assert_failure
    assert_output --partial "only one index type can be specified"
} 

# bats test_tags=action:create, validation:constraints
@test "Index name validation" {
    # Test very long index name
    local long_name=$(printf 'a%.0s' {1..100})
    run $CLI index create "${long_name}" --yes
    assert_failure
    assert_output --partial "Name too long, please use names shorter than 45 characters"
    
    # Test special characters in index name
    run $CLI index create "test@index" --yes
    assert_failure
    assert_output --partial "Name must consist of lower case alphanumeric characters or '-'"
    
    run $CLI index create "test#index" --yes
    assert_failure
    assert_output --partial "Name must consist of lower case alphanumeric characters or '-'"
    
    run $CLI index create "test\$index" --yes
    assert_failure
    assert_output --partial "Name must consist of lower case alphanumeric characters or '-'"
    
    # Test invalid name that should fail
    run $CLI index create "Test" --yes
    assert_failure
    assert_output --partial "Name must consist of lower case alphanumeric characters or '-'"
}