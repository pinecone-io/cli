#!/usr/bin/env bash

# Load helper functions
load "$BATS_ROOT/../helpers/custom/all.bash"

# =============================================================================
# Deletion Protection Tests
# =============================================================================
# Set tags for all tests in this file
#
# -----------------------------------------------------------------------------
# bats file_tags=scope:index, scenario:protection_flow
# =============================================================================


# -----------------------------------------------------------------------------
# Setup and teardown
# -----------------------------------------------------------------------------

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

# bats test_tags=action:CRUD, mode:flags
@test "Index with deletion protection can be deleted after removing protection" {
    local template_name="serverless_with_deletion_protection"
    local cli_params=$(extract_cli_params_from_template "$template_name")

    # Create index with deletion protection enabled
    $CLI index create ${TEST_INDEX_NAME} $cli_params -y
    
    # Wait for index to be ready
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    
    # Verify deletion protection is enabled
    local deletion_protection=$(echo "$index_json" | jq -r '.deletion_protection')
    [ "$deletion_protection" = "enabled" ]
    

    
    # Attempt deletion while protection is enabled - should fail
    run $CLI index delete ${TEST_INDEX_NAME}
    assert_failure
    
    # Verify the index still exists
    local index_json_still_exists=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    [ -n "$index_json_still_exists" ]
    
    # Remove deletion protection
    $CLI index configure --name ${TEST_INDEX_NAME} --deletion_protection disabled
    
    # Wait for configuration update to complete
    local index_json_after=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local deletion_protection_after=$(echo "$index_json_after" | jq -r '.deletion_protection')
    [ "$deletion_protection_after" = "disabled" ]
    
    # Now delete the index - this should succeed
    $CLI index delete ${TEST_INDEX_NAME}
        
    # Verify the index was actually deleted
    run $CLI index describe ${TEST_INDEX_NAME}
    assert_failure
}

# bats test_tags=action:CRUD, mode:flags
@test "Index deletion protection can be enabled after creation" {
    local template_name="serverless_aws"
    local cli_params=$(extract_cli_params_from_template "$template_name")

    # Create index without deletion protection
    $CLI index create ${TEST_INDEX_NAME} $cli_params -y
    
    # Wait for index to be ready
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    
    # Verify deletion protection is disabled initially
    local deletion_protection=$(echo "$index_json" | jq -r '.deletion_protection')
    [ "$deletion_protection" = "disabled" ]
    
    # Enable deletion protection
    $CLI index configure --name ${TEST_INDEX_NAME} --deletion_protection enabled
    
    # Wait for configuration update to complete
    local index_json_after=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local deletion_protection_after=$(echo "$index_json_after" | jq -r '.deletion_protection')
    [ "$deletion_protection_after" = "enabled" ]
    
    # Try to delete - should fail
    run $CLI index delete ${TEST_INDEX_NAME}
    assert_failure
    
    # Clean up by removing deletion protection first
    $CLI index configure --name ${TEST_INDEX_NAME} --deletion_protection disabled
    
    # Wait for configuration update to complete
    local index_json_final=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    
    # Now delete the index
    $CLI index delete ${TEST_INDEX_NAME}
    
    # Verify the index was actually deleted
    run $CLI index describe ${TEST_INDEX_NAME}
    assert_failure
}
