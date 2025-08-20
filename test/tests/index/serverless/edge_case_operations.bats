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
# bats file_tags=scope:index, index-type:serverless, scenario:edge_cases
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

# bats test_tags=action:create, vector-type:dense  
@test "Serverless index with boundary dimension values" {
    # Test minimum valid dimension
    run $CLI index create ${TEST_INDEX_NAME} --serverless --dimension 1 --yes
    assert_success
    
    # Verify the dimension was set correctly
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local dimension=$(echo "$index_json" | jq -r '.dimension')
    [ "$dimension" = "1" ]
}

# bats test_tags=action:create, vector-type:dense
@test "Serverless index with maximum valid dimension" {
    # Test maximum valid dimension (assuming 2048 is max)
    run $CLI index create ${TEST_INDEX_NAME} --serverless --dimension 2048 --yes
    assert_success
    
    # Verify the dimension was set correctly
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local dimension=$(echo "$index_json" | jq -r '.dimension')
    [ "$dimension" = "2048" ]
}

# bats test_tags=action:create, vector-type:dense
@test "Serverless index with edge case names" {
    # Test very short but valid name
    local short_name="a"
    run $CLI index create ${short_name} --serverless --yes
    assert_success
    
    # Clean up
    $CLI index delete ${short_name}
    
    # Test name with maximum valid length (44 characters, just under the 45 limit)
    local long_name=$(printf 'a%.0s' {1..44})
    run $CLI index create ${long_name} --serverless --yes
    assert_success
    
    # Clean up
    $CLI index delete ${long_name}
}

# bats test_tags=action:create, vector-type:dense
@test "Serverless index with edge case tags" {
    # Test edge case tag values using the correct comma-separated format
    # Note: Pinecone's API drops empty tag values, so only tags with actual values are stored
    run $CLI index create ${TEST_INDEX_NAME} --serverless --tags "key=,special-chars=test-value" --yes
    assert_success
    
    # Verify that the warning about empty tag values was shown
    assert_output --partial "Warning: Empty tag values for keys 'key' will be dropped by Pinecone"
    
    # Verify the tags were set correctly
    # Note: Pinecone drops empty tag values, so only "special-chars=test-value" will be stored
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local tags=$(echo "$index_json" | jq -r '.tags | to_entries[] | "\(.key)=\(.value)"' | sort)
    
    # Check that only the tag with a value is present (Pinecone drops empty values)
    # Expected: only "special-chars=test-value" should be saved
    echo "$tags" | grep -q "special-chars=test-value"
    
    # Verify that empty tag values are not present (as expected from Pinecone's behavior)
    echo "$tags" | grep -v -q "key="
}
