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
# bats file_tags=scope:index, index-type:serverless, scenario:error
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
@test "Serverless index cannot use pod-specific flags" {
    # Test that using serverless with pod-specific flags fails
    run $CLI index create ${TEST_INDEX_NAME} --serverless --environment "us-east-1-aws" --yes
    assert_failure
    assert_output --partial "cannot be used with serverless indexes"
    
    run $CLI index create ${TEST_INDEX_NAME} --serverless --pod_type "p1.x1" --yes
    assert_failure
    assert_output --partial "cannot be used with serverless indexes"
    
    run $CLI index create ${TEST_INDEX_NAME} --serverless --shards 2 --yes
    assert_failure
    assert_output --partial "cannot be used with serverless indexes"
    
    run $CLI index create ${TEST_INDEX_NAME} --serverless --replicas 2 --yes
    assert_failure
    assert_output --partial "cannot be used with serverless indexes"
}

# bats test_tags=action:create, validation:constraints
@test "Serverless index cannot use integrated-specific flags" {
    # Test that using serverless with integrated-specific flags fails
    run $CLI index create ${TEST_INDEX_NAME} --serverless --model "multilingual-e5-large" --yes
    assert_failure
    assert_output --partial "cannot be used with serverless indexes"
    
    run $CLI index create ${TEST_INDEX_NAME} --serverless --field_map "text=chunk_text" --yes
    assert_failure
    assert_output --partial "cannot be used with serverless indexes"
    
    run $CLI index create ${TEST_INDEX_NAME} --serverless --read_parameters "input_type=query" --yes
    assert_failure
    assert_output --partial "cannot be used with serverless indexes"
    
    run $CLI index create ${TEST_INDEX_NAME} --serverless --write_parameters "input_type=passage" --yes
    assert_failure
    assert_output --partial "cannot be used with serverless indexes"
}

# bats test_tags=action:create, validation:flags
@test "Serverless index with invalid cloud provider" {
    run $CLI index create ${TEST_INDEX_NAME} --serverless --cloud "invalid-cloud" --yes
    assert_failure
    assert_output --partial "Resource cloud: invalid-cloud region: us-east-1 not found"
}

# bats test_tags=action:create, validation:flags
@test "Serverless index with invalid region for cloud" {
    # Test invalid region for AWS
    run $CLI index create ${TEST_INDEX_NAME} --serverless --cloud aws --region "invalid-region" --yes
    assert_failure
    
    # Test invalid region for GCP
    run $CLI index create ${TEST_INDEX_NAME} --serverless --cloud gcp --region "invalid-region" --yes
    assert_failure
    
    # Test invalid region for Azure
    run $CLI index create ${TEST_INDEX_NAME} --serverless --cloud azure --region "invalid-region" --yes
    assert_failure
}

# bats test_tags=action:create, validation:flags
@test "Serverless index with invalid dimension values" {
    # Test negative dimension
    run $CLI index create ${TEST_INDEX_NAME} --serverless --dimension -1 --yes
    assert_failure
    assert_output --partial "dimension: invalid value: integer \`-1\`, expected u32"
    
    # Test zero dimension
    run $CLI index create ${TEST_INDEX_NAME} --serverless --dimension 0 --yes
    assert_failure
    assert_output --partial "Dimension is required for dense vector index"
    
    # Test dimension too large
    run $CLI index create ${TEST_INDEX_NAME} --serverless --dimension 100000 --yes
    assert_failure
    assert_output --partial "Must be greater than 0 and less than 20,000"
    
    # Test non-numeric dimension
    run $CLI index create ${TEST_INDEX_NAME} --serverless --dimension "invalid" --yes
    assert_failure
    assert_output --partial "strconv.ParseInt: parsing \"invalid\": invalid syntax"
}

# bats test_tags=action:create, validation:flags
@test "Serverless index with invalid metric" {
    run $CLI index create ${TEST_INDEX_NAME} --serverless --metric "invalid-metric" --yes
    assert_failure
    assert_output --partial "Invalid metric"
}

# bats test_tags=action:create, validation:flags
@test "Serverless index with invalid vector type" {
    run $CLI index create ${TEST_INDEX_NAME} --serverless --vector_type "invalid-type" --yes
    assert_failure
    assert_output --partial "unsupported VectorType: invalid-type"
}

# bats test_tags=action:create, validation:flags
@test "Serverless index with invalid deletion protection value" {
    run $CLI index create ${TEST_INDEX_NAME} --serverless --deletion_protection "invalid" --yes
    assert_failure
    assert_output --partial "Invalid deletion_protection, value should be either enabled or disabled"
}

# bats test_tags=action:create, validation:flags
@test "Serverless index with invalid tags format" {
    # Test invalid tag format (missing equals sign)
    run $CLI index create ${TEST_INDEX_NAME} --serverless --tags "invalid-tag" --yes
    assert_failure
    assert_output --partial "invalid-tag must be formatted as key=value"
    
    # Test invalid tag format (empty key)
    run $CLI index create ${TEST_INDEX_NAME} --serverless --tags "=value" --yes
    assert_failure
    assert_output --partial "Keys must not be empty"
    
    # Test multiple separate tag arguments (should fail - CLI expects comma-separated format)
    run $CLI index create ${TEST_INDEX_NAME} --serverless --tags "key1=value1" "key2=value2" --yes
    assert_failure
    assert_output --partial "unknown flag"
}

# bats test_tags=action:create, validation:flags
@test "Serverless index with invalid tags argument format" {
    # Test that --tags with separate arguments fails (CLI expects comma-separated format)
    # Note: Current CLI behavior is to accept only the first tag and ignore subsequent arguments
    # This test documents the expected behavior that should reject multiple separate arguments
    # The current behavior could be confusing for users and should be fixed
    run $CLI index create ${TEST_INDEX_NAME} --serverless --tags "key1=value1" "key2=value2" --yes
    assert_failure
    assert_output --partial "unknown flag"
    
    # Test that --tags with mixed format fails
    run $CLI index create ${TEST_INDEX_NAME} --serverless --tags "key1=value1,key2=value2" "key3=value3" --yes
    assert_failure
    assert_output --partial "unknown flag"
}

# bats test_tags=action:create, validation:flags
@test "Serverless index with boundary dimension violations" {
    # Test dimension below minimum (0 is already tested in invalid dimension values)
    run $CLI index create ${TEST_INDEX_NAME} --serverless --dimension -1 --yes
    assert_failure
    assert_output --partial "dimension: invalid value: integer \`-1\`, expected u32"
    
    # Test dimension at zero (should fail)
    run $CLI index create ${TEST_INDEX_NAME} --serverless --dimension 0 --yes
    assert_failure
    assert_output --partial "Dimension is required for dense vector index"
    
    # Test dimension above maximum (assuming 20,000 is max based on previous error)
    run $CLI index create ${TEST_INDEX_NAME} --serverless --dimension 25000 --yes
    assert_failure
    assert_output --partial "Must be greater than 0 and less than 20,000"
}

# bats test_tags=action:create, validation:constraints
@test "Serverless index with sparse vector type must use dotproduct metric" {
    # Test that sparse vectors with non-dotproduct metrics fail
    run $CLI index create ${TEST_INDEX_NAME} --serverless --vector_type sparse --metric cosine --yes
    assert_failure
    assert_output --partial "sparse vector type requires dotproduct metric"
    
    run $CLI index create ${TEST_INDEX_NAME} --serverless --vector_type sparse --metric euclidean --yes
    assert_failure
    assert_output --partial "sparse vector type requires dotproduct metric"
    
    # Test that sparse vectors with dotproduct metric succeed
    run $CLI index create ${TEST_INDEX_NAME} --serverless --vector_type sparse --metric dotproduct --yes
    assert_success
    
    # Clean up
    $CLI index delete ${TEST_INDEX_NAME}
}
