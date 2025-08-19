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
# bats file_tags=scope:_metatests
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
    :
}

# If something needs to be done after each test in this file, do it here.
teardown() {
    :
}

# -----------------------------------------------------------------------------
# Tests
# -----------------------------------------------------------------------------


@test "load template from file" {
    local template_json=$(load_index_template "serverless_default")
    [ -n "$template_json" ]
    
    # Verify it's valid JSON
    echo "$template_json" | jq . >/dev/null
}

@test "fail when template file is missing" {
    run load_index_template "nonexistent_template" 2>&1
    assert_failure
    [[ "$output" == *"Template file not found"* ]]
}

@test "fail when template file has invalid JSON" {
    # Create a temporary invalid JSON file
    local temp_dir=$(mktemp -d)
    local invalid_template="$temp_dir/invalid.json"
    echo '{"invalid": json}' > "$invalid_template"
    
    run load_index_template "$temp_dir/invalid" 2>&1
    assert_failure
    [[ "$output" == *"Invalid JSON in template file"* ]]
    
    rm -rf "$temp_dir"
}

@test "fail validation when template cannot be loaded" {
    local valid_json='{"name": "test", "metric": "cosine"}'
    
    # Use a non-existent template name - should fail to load
    run assert_index_json_matches_template "$valid_json" "nonexistent_template"
    assert_failure
    [[ "$output" == *"Failed to load template 'nonexistent_template' from file"* ]]
}
