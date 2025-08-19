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
# bats file_tags=scope:index, index-type:default, scenario:success
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

# bats test_tags=action:CRD, mode:params, index-type:default, cloud:default
@test "Create, read and delete an index with default behavior (serverless)" {
    $CLI index create ${TEST_INDEX_NAME} -y
    
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local host_value=$(echo "$index_json" | jq -r '.host')
    local placeholders_values="
        __NAME__:${TEST_INDEX_NAME} 
        __HOST__:$host_value
    "
    assert_index_json_matches_template "$index_json" "serverless_default" "$placeholders_values"
    
}


# bats test_tags=action:CRD, mode:interactive, index-type:default, cloud:default
@test "Interactive mode works with no parameters and defaults to serverless" {
    # Test interactive mode by providing all inputs programmatically
    # The interactive flow is:
    # 1. Index name
    # 2. Index type (serverless/pod/integrated) - default: serverless
    # 3. Vector type (dense/sparse) - default: dense
    # 4. Metric (cosine/euclidean/dotproduct) - default: cosine
    # 5. Cloud provider (aws/gcp/azure) - default: aws
    # 6. Region (depends on cloud) - default: us-east-1
    # 7. Dimension (if dense vector) - default: 1536
    # 8. Confirmation prompt
    
    # All the answers: name + defaults + final yes
    answers="${TEST_INDEX_NAME}\n\n\n\n\n\n1536\ny"
    # Use bash -c to pipe the answers to the CLI
    run bash -c "echo -e '${answers}' | $CLI index create"
    
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local host_value=$(echo "$index_json" | jq -r '.host')
    local placeholders_values="
        __NAME__:${TEST_INDEX_NAME} 
        __HOST__:$host_value
    "
    assert_index_json_matches_template "$index_json" "serverless_default" "$placeholders_values"

}
