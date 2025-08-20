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
# bats file_tags=scope:index, mode:interactive, scenario:success
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

# bats test_tags=action:CRD, cloud:default
@test "Interactive mode works for serverless index creation" {

    #####
    # TODO: Verify index type is not shown in the menu
    #####
    
    # All the answers: name + serverless + defaults + final yes
    answers="${TEST_INDEX_NAME}\nserverless\n\n\n\n\n1536\ny"


    # Use bash -c to pipe the answers to the CLI
    run bash -c "echo -e '${answers}' | $CLI index create --serverless"
    
    
    # Wait for index to be ready
    local index_json=$(index_describe_wait_for_ready ${TEST_INDEX_NAME})
    local host_value=$(echo "$index_json" | jq -r '.host')
    local placeholders_values="
        __NAME__:${TEST_INDEX_NAME} 
        __HOST__:$host_value
    "
    local template_name="serverless_aws"
    assert_index_json_matches_template_file "$index_json" "$template_name" "$placeholders_values"

}