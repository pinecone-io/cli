# Get machine identifier for unique naming
get_machine_id() {
    echo "$(head -c 12 /etc/machine-id 2>/dev/null || hostname | sha1sum | cut -c1-12)"
}

# Generate unique test index name: t-{machine-id}-{timestamp}-{uuid}
generate_index_name() {
    echo "t-$(get_machine_id)-$(date +%s)-$(uuidgen | xxd -r -p | base64 | tr -d '=/+' | cut -c1-10 | tr '[:upper:]' '[:lower:]')"
}

# Extract JSON from CLI output (finds line starting with {)
extract_json_from_output() {
    local output="$1"
    # Find the line that starts with { and extract from there
    echo "$output" | awk '/^{/,0'
}


# Wait for index to reach Ready state
# Polls index status until Ready, Failed, or timeout (1 minute)
# Handles binary call failures with retry logic
#
# Usage: index_json=$(index_describe_wait_for_ready "index-name")
# Returns: 0 on success (JSON to stdout), 1 on failure
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
# JSON Template Validation Functions
# =============================================================================

# Validate JSON against template from file
# Usage: assert_index_json_matches_template_file "$actual_json" "serverless_default"
# Returns: 0 on match, 1 on mismatch
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


# Load JSON template from file or full path
# Usage: template_json=$(load_index_template "serverless_default")
# Returns: JSON string on success, error message on failure
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


# Core JSON validation using jd tool
# Handles placeholder replacement and template matching
# Usage: assert_json_matches_template_jd "$actual_json" "$template_json" "$placeholders_values"
# Returns: 0 on match, 1 on mismatch
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

# Convert JSON template to CLI flags
# Usage: cli_params=$(extract_cli_params_from_template "serverless_aws")
# Returns: CLI flags string on success, error message on failure
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

