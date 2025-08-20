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
# bats file_tags=scope:global
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

# bats test_tags=rationale:feature-request
@test "Use flag --version to show version" {
    run $CLI --version
    assert_success
    assert_output --partial "version"
}

# bats test_tags=rationale:feature-request
@test "Use short flag -v to show version" {
    run $CLI --version
    assert_success
    assert_output --partial "version"
}