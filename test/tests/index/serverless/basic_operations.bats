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
# bats file_tags=scope:index, index-type:serverless, scenario:success
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

# bats test_tags=action:CRD, mode:params, cloud:default, vector-type:dense
@test "Create, read and delete a serverless index with minimal flags" {
    local template_name="serverless_aws"
    local cli_params=$(extract_cli_params_from_template "$template_name")

    $CLI index create ${TEST_INDEX_NAME} $cli_params -y
    
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local host_value=$(echo "$index_json" | jq -r '.host')
    local placeholders_values="
        __NAME__:${TEST_INDEX_NAME} 
        __HOST__:$host_value
    "
    assert_index_json_matches_template_file "$index_json" "$template_name" "$placeholders_values"
    
}

# bats test_tags=action:CRD, mode:params, cloud:aws, vector-type:dense
@test "Create serverless index with AWS cloud and custom region" {
    local template_name="serverless_aws_west"
    local cli_params=$(extract_cli_params_from_template "$template_name")

    $CLI index create ${TEST_INDEX_NAME} $cli_params -y
    
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local host_value=$(echo "$index_json" | jq -r '.host')
    local placeholders_values="
        __NAME__:${TEST_INDEX_NAME} 
        __HOST__:$host_value
    "
    assert_index_json_matches_template_file "$index_json" "$template_name" "$placeholders_values"
    
}

# bats test_tags=action:CRD, mode:params, cloud:gcp, vector-type:dense
@test "Create serverless index with GCP cloud and custom region" {
    local template_name="serverless_gcp"
    local cli_params=$(extract_cli_params_from_template "$template_name")

    $CLI index create ${TEST_INDEX_NAME} $cli_params -y
    
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local host_value=$(echo "$index_json" | jq -r '.host')
    local placeholders_values="
        __NAME__:${TEST_INDEX_NAME} 
        __HOST__:$host_value
    "
    assert_index_json_matches_template_file "$index_json" "$template_name" "$placeholders_values"
    
}

# bats test_tags=action:CRD, mode:params, cloud:azure, vector-type:dense
@test "Create serverless index with Azure cloud and custom region" {
    local template_name="serverless_azure"
    local cli_params=$(extract_cli_params_from_template "$template_name")

    $CLI index create ${TEST_INDEX_NAME} $cli_params -y
    
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local host_value=$(echo "$index_json" | jq -r '.host')
    local placeholders_values="
        __NAME__:${TEST_INDEX_NAME} 
        __HOST__:$host_value
    "
    assert_index_json_matches_template_file "$index_json" "$template_name" "$placeholders_values"
    
}

# bats test_tags=action:CRD, mode:params, validation:format, vector-type:dense
@test "Create serverless index with custom dimension" {
    local template_name="serverless_euclidean_768"
    local cli_params=$(extract_cli_params_from_template "$template_name")

    $CLI index create ${TEST_INDEX_NAME} $cli_params -y
    
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local host_value=$(echo "$index_json" | jq -r '.host')
    local placeholders_values="
        __NAME__:${TEST_INDEX_NAME} 
        __HOST__:$host_value
    "
    assert_index_json_matches_template_file "$index_json" "$template_name" "$placeholders_values"
    
}

# bats test_tags=action:CRD, mode:params, validation:format, vector-type:dense
@test "Create serverless index with custom metric" {
    local template_name="serverless_euclidean"
    local cli_params=$(extract_cli_params_from_template "$template_name")

    $CLI index create ${TEST_INDEX_NAME} $cli_params -y
    
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local host_value=$(echo "$index_json" | jq -r '.host')
    local placeholders_values="
        __NAME__:${TEST_INDEX_NAME} 
        __HOST__:$host_value
    "
    assert_index_json_matches_template_file "$index_json" "$template_name" "$placeholders_values"
    
}

# bats test_tags=action:CRD, mode:params, validation:format, vector-type:dense
@test "Create serverless index with dense vector type" {
    local template_name="serverless_aws"
    local cli_params=$(extract_cli_params_from_template "$template_name")

    $CLI index create ${TEST_INDEX_NAME} $cli_params -y
    
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local host_value=$(echo "$index_json" | jq -r '.host')
    local placeholders_values="
        __NAME__:${TEST_INDEX_NAME} 
        __HOST__:$host_value
    "
    assert_index_json_matches_template_file "$index_json" "$template_name" "$placeholders_values"
    
}

# bats test_tags=action:CRD, mode:params, validation:format, vector-type:sparse
@test "Create serverless index with sparse vector type" {
    local template_name="serverless_sparse"
    local cli_params=$(extract_cli_params_from_template "$template_name")

    $CLI index create ${TEST_INDEX_NAME} $cli_params -y
    
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local host_value=$(echo "$index_json" | jq -r '.host')
    local placeholders_values="
        __NAME__:${TEST_INDEX_NAME} 
        __HOST__:$host_value
    "
    assert_index_json_matches_template_file "$index_json" "$template_name" "$placeholders_values"
    
}



# bats test_tags=action:CRD, mode:params, validation:format, vector-type:dense
@test "Create serverless index with tags" {
    local template_name="serverless_with_tags"
    local cli_params=$(extract_cli_params_from_template "$template_name")

    $CLI index create ${TEST_INDEX_NAME} $cli_params -y
    
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local host_value=$(echo "$index_json" | jq -r '.host')
    local placeholders_values="
        __NAME__:${TEST_INDEX_NAME} 
        __HOST__:$host_value
    "
    assert_index_json_matches_template_file "$index_json" "$template_name" "$placeholders_values"
    
}
