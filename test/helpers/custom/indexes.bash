get_machine_id() {
    echo "$(head -c 12 /etc/machine-id 2>/dev/null || hostname | sha1sum | cut -c1-12)"
}

# Generate a unique index name with format: t-{timestamp}-{base64-uuid}
# The name only contains lowercase alphanumeric characters and dashes
# - t: prefix for test
# - timestamp: Unix timestamp for chronological ordering
# - base64-uuid: base64-encoded UUID with unsafe characters removed (for uniqueness)
generate_index_name() {
    echo "t-$(get_machine_id)-$(date +%s)-$(uuidgen | xxd -r -p | base64 | tr -d '=/+' | cut -c1-10 | tr '[:upper:]' '[:lower:]')"
}

# Helper function to extract JSON from CLI output
# This function finds the line that starts with { and extracts from there
extract_json_from_output() {
    local output="$1"
    # Find the line that starts with { and extract from there
    echo "$output" | awk '/^{/,0'
}


# Helper function to wait for index to be ready
# 
# This function polls an index until it reaches a "Ready" state or encounters a terminal error.
# It handles various index statuses appropriately:
# - Success: Returns 0 and outputs the index JSON to stdout when status is "Ready"
# - Error: Returns 1 for terminal states like "Failed", "InitializationFailed", "Terminating", "Disabled"
# - Wait: Continues polling for transient states like "Initializing", "ScalingUp", "ScalingDown", etc.
# 
# The function also handles binary call failures gracefully:
# - Fails fast for permanent errors (not found, unauthorized)
# - Retries transient errors up to 3 times before giving up
# - Times out after 1 minute (12 attempts × 5 seconds)
#
# Usage:
#   index_json=$(index_describe_wait_for_ready "index-name")
#   if [ $? -eq 0 ]; then
#       # Process the JSON output
#       index_id=$(echo "$index_json" | jq -r '.id')
#   fi
#
# Parameters:
#   $1: index_name - The name of the index to monitor
#
# Returns:
#   0: Index is ready (JSON output to stdout)
#   1: Index failed, terminated, or other error
#
# Output:
#   - JSON representation of the index to stdout on success
#   - Status messages and errors to stderr
index_describe_wait_for_ready() {
    local index_name="$1"
    local max_attempts=12  # 1 minute (12 attempts * 5 seconds)
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        local status
        local raw_output
        local json_output
        local exit_code
        
        # Capture both output and exit code
        raw_output=$($CLI index describe "$index_name" --json 2>&1)
        exit_code=$?
        
        # Check if the binary call itself failed
        if [ $exit_code -ne 0 ]; then
            echo "Binary call failed with exit code $exit_code: $raw_output (attempt $attempt/$max_attempts)" >&2
            
            # If it's a permanent error (like index not found), fail fast
            if echo "$raw_output" | grep -q "not found\|404\|not authorized\|unauthorized" 2>/dev/null; then
                echo "Permanent error detected, failing fast" >&2
                return 1
            fi
            
            # For other errors, retry a few times then fail
            if [ $attempt -ge 3 ]; then
                echo "Binary call failed repeatedly, giving up" >&2
                return 1
            fi
        fi
        
        json_output=$(extract_json_from_output "$raw_output" 2>/dev/null || echo "")
        status=$(echo "$json_output" | jq -r '.status.state' 2>/dev/null || echo "unknown")
        
        case "$status" in
            "Ready")
                echo "$json_output"
                return 0
                ;;
            "InitializationFailed"|"Failed"|"Terminating"|"Disabled")
                return 1
                ;;
            "Initializing"|"ScalingUp"|"ScalingDown"|"ScalingUpPodSize"|"ScalingDownPodSize")
                # Continue waiting
                ;;
            *)
                # Unknown status, continue waiting
                ;;
        esac
        
        sleep 5
        attempt=$((attempt + 1))
    done
    
    echo "Index $index_name did not reach Ready state within 1 minute timeout" >&2
    return 1
}

# =============================================================================
# JSON Comparison and Validation Functions
# =============================================================================

# Compare an actual index JSON response with an expected template
# This function validates the structure and key values while allowing for
# dynamic values like names, hosts, and timestamps
#
# Usage:
#   assert_index_json_matches_template "$actual_json" "$expected_template"
#   assert_index_json_matches_template "$actual_json" "serverless_default"
#   assert_index_json_matches_template "$actual_json" "$(load_index_template 'serverless_default')"
#
# Parameters:
#   $1: actual_json - The actual JSON response from the CLI
#   $2: expected_template - Either a template name (to load from file) or a JSON string
#
# Returns:
#   0: JSON matches expected template
#   1: JSON does not match (with detailed error message)
#
# Note: The function automatically detects whether the input is a template name
# (to load from file) or a JSON string. Template names that don't exist will
# cause the function to fail.
#
# This function now uses jd-based validation for efficient template matching.
# Expected values can be passed as a third parameter to replace placeholders.
assert_index_json_matches_template() {
    local actual_json="$1"
    local expected_template="$2"
    local placeholders_values="$3"
    
    # Check if actual_json is valid JSON
    if ! echo "$actual_json" | jq . >/dev/null 2>&1; then
        echo "Error: actual_json is not valid JSON" >&2
        return 1
    fi
    
    # Get the expected template JSON
    local expected_json
    
    # First, try to load as a template file
    expected_json=$(load_index_template "$expected_template" 2>/dev/null)
    if [ $? -eq 0 ] && [ -n "$expected_json" ]; then
        # Successfully loaded template from file
        :
    else
        # Check if it looks like a JSON string (starts with { or [)
        if [[ "$expected_template" =~ ^[[:space:]]*[{\[] ]]; then
            # Assume it's a direct JSON string
            expected_json="$expected_template"
        else
            # Failed to load template and doesn't look like JSON
            echo "Error: Failed to load template '$expected_template' from file" >&2
            return 1
        fi
    fi
    
    # Use jd-based template matching for simple and efficient validation
    if assert_json_matches_template_jd "$actual_json" "$expected_json" "$placeholders_values"; then
        return 0
    else
        return 1
    fi
}

# Load a JSON template from a file
# This allows for more complex templates and easier maintenance
#
# Usage:
#   template_json=$(load_index_template "serverless_default")
#   assert_index_json_matches_template "$actual_json" "$template_json"
#
# Parameters:
#   $1: template_name - The name of the template file (without .json extension)
#
# Returns:
#   JSON template string to stdout on success
#   Exits with error message if template file not found or invalid
load_index_template() {
    local template_name="$1"
    # Use relative path from the test directory
    local template_file="$BATS_ROOT/../helpers/templates/indexes/${template_name}.json"
    
    if [ ! -f "$template_file" ]; then
        echo "Error: Template file not found: $template_file" >&2
        return 1
    fi
    
    local template_content
    template_content=$(cat "$template_file")
    
    # Validate that the template is valid JSON
    if ! echo "$template_content" | jq . >/dev/null 2>&1; then
        echo "Error: Invalid JSON in template file: $template_file" >&2
        return 1
    fi
    
    echo "$template_content"
}


# Recursively validate JSON against a template
# This function treats every field in the template as required and checks if values match
#
# Usage:
#   assert_json_matches_template_recursive "$actual_json" "$template_json"
#
# Parameters:
#   $1: actual_json - The JSON data to validate
#   $2: template_json - The template JSON to validate against
#
# Returns:
#   0: JSON matches template structure and values
#   1: JSON does not match template (with detailed error message)
assert_json_matches_template_jd() {
    local actual_json="$1"
    local template_json="$2"
    local placeholders_values="$3"

    # Check if inputs are valid JSON
    if ! echo "$actual_json" | jq . >/dev/null 2>&1; then
        echo "Error: actual_json is not valid JSON" >&2
        return 1
    fi
    if ! echo "$template_json" | jq . >/dev/null 2>&1; then
        echo "Error: template_json is not valid JSON" >&2
        return 1
    fi

    # Replace placeholders in the template with expected values using simple string replacement
    local processed_template="$template_json"
    
    # Simple key:value replacement for placeholders
    if [ -n "$placeholders_values" ]; then
        # Parse space or newline-separated key:value pairs and do simple string replacement
        # Convert newlines to spaces for consistent processing
        local normalized_values=$(echo "$placeholders_values" | tr '\n' ' ')
        
        for pair in $normalized_values; do
            if [[ "$pair" =~ ^([^:]+):(.+)$ ]]; then
                local key="${BASH_REMATCH[1]}"
                local value="${BASH_REMATCH[2]}"
                processed_template=$(echo "$processed_template" | sed "s/$key/$value/g")
            fi
        done
    fi

    # Use jd -set to compare the processed template with actual JSON
    local diff_output
    diff_output=$(jd -set <(echo "$processed_template") <(echo "$actual_json") 2>/dev/null)
    
    if [ $? -eq 0 ] && [ -z "$diff_output" ]; then
        # No differences found - validation passed
        return 0
    else
        # Differences found - validation failed
        echo "Error: JSON does not match template. Differences:" >&2
        echo "$diff_output" >&2
        return 1
    fi
}

# Extract CLI parameters from a JSON template
# This function reads a template file and converts it to CLI command flags
#
# Usage:
#   cli_params=$(extract_cli_params_from_template "serverless_aws")
#   $CLI index create ${NAME} $cli_params -y
#
# Parameters:
#   $1: template_name - The name of the template file (without .json extension)
#
# Returns:
#   CLI flags string to stdout on success
#   Exits with error message if template file not found or invalid
extract_cli_params_from_template() {
    local template_name="$1"
    local template_json
    
    # Load the template
    template_json=$(load_index_template "$template_name")
    if [ $? -ne 0 ]; then
        return 1
    fi
    
    local params=""
    
    # Extract serverless flag if this is a serverless template
    if echo "$template_json" | jq -e '.spec.serverless' >/dev/null 2>&1; then
        params="$params --serverless"
        
        # Extract cloud and region
        local cloud=$(echo "$template_json" | jq -r '.spec.serverless.cloud')
        if [ "$cloud" != "null" ] && [ "$cloud" != "__CLOUD__" ]; then
            params="$params --cloud $cloud"
        fi
        
        local region=$(echo "$template_json" | jq -r '.spec.serverless.region')
        if [ "$region" != "null" ] && [ "$region" != "__REGION__" ]; then
            params="$params --region $region"
        fi
        

        
        local source_collection=$(echo "$template_json" | jq -r '.spec.serverless.sourceCollection')
        if [ "$source_collection" != "null" ] && [ "$source_collection" != "__SOURCE_COLLECTION__" ]; then
            params="$params --source_collection $source_collection"
        fi
    fi
    
    # Extract metric
    local metric=$(echo "$template_json" | jq -r '.metric')
    if [ "$metric" != "null" ] && [ "$metric" != "__METRIC__" ]; then
        params="$params --metric $metric"
    fi
    
    # Extract vector type first (needed to determine if dimension should be extracted)
    local vector_type=$(echo "$template_json" | jq -r '.vector_type')
    if [ "$vector_type" != "null" ] && [ "$vector_type" != "__VECTOR_TYPE__" ]; then
        params="$params --vector_type $vector_type"
    fi
    
    # Extract dimension (only for dense vectors)
    if [ "$vector_type" != "sparse" ]; then
        local dimension=$(echo "$template_json" | jq -r '.dimension')
        if [ "$dimension" != "null" ] && [ "$dimension" != "__DIMENSION__" ]; then
            params="$params --dimension $dimension"
        fi
    fi
    
    # Extract deletion protection
    local deletion_protection=$(echo "$template_json" | jq -r '.deletion_protection')
    if [ "$deletion_protection" != "null" ] && [ "$deletion_protection" != "__DELETION_PROTECTION__" ]; then
        if [ "$deletion_protection" = "enabled" ]; then
            params="$params --deletion_protection enabled"
        fi
    fi
    
    # Extract tags (if present)
    local tags=$(echo "$template_json" | jq -r '.tags | to_entries[] | "\(.key)=\(.value)"' 2>/dev/null)
    if [ $? -eq 0 ] && [ -n "$tags" ]; then
        while IFS= read -r tag; do
            if [ "$tag" != "null" ] && [[ "$tag" != *"__"* ]]; then
                params="$params --tags \"$tag\""
            fi
        done <<< "$tags"
    fi
    
    echo "$params"
}

