#!/usr/bin/env bats

# Integration tests for index flow: create -> describe -> delete
# Tests all possible combinations of flags for index create command
# Also tests invalid and missing values

setup() {
    # Setup test environment
    export TEST_INDEX_PREFIX="test-index-$(date +%s)"
    
    # Initialize array to track created indexes for this test
    TEST_CREATED_INDEXES=()
    
    # Use PC_BINARY from environment or default to relative path
    if [ -z "${PC_BINARY:-}" ]; then
        export PC_BINARY="../../pcdev"
    fi
    
    # Check if we should skip login check
    if [ "${SKIP_LOGIN:-}" = "true" ]; then
        echo "Skipping login check (SKIP_LOGIN=true)"
        return 0
    fi
    
    # Ensure we're logged in (this should be done in CI/CD environment)
    # For local testing, you may need to run: pc login
    if ! $PC_BINARY whoami >/dev/null 2>&1; then
        skip "Not logged in to Pinecone. Run 'pc login' first."
    fi
}

teardown() {
    echo "Cleaning up test indexes..."
    
    # First, try to delete indexes tracked by this test
    if [ ${#TEST_CREATED_INDEXES[@]} -gt 0 ]; then
        echo "Deleting tracked indexes: ${TEST_CREATED_INDEXES[*]}"
        for index_name in "${TEST_CREATED_INDEXES[@]}"; do
            if [ -n "$index_name" ]; then
                echo "Deleting tracked index: $index_name"
                $PC_BINARY index delete "$index_name" >/dev/null 2>&1 || echo "Failed to delete tracked index: $index_name"
            fi
        done
    fi
    
    # Then, cleanup any remaining test indexes (fallback)
    local test_indexes
    test_indexes=$($PC_BINARY index list --json 2>/dev/null | jq -r '.[] | select(.name | startswith("test-index-")) | .name' 2>/dev/null || true)
    
    if [ -n "$test_indexes" ]; then
        echo "Found remaining test indexes, cleaning up..."
        echo "$test_indexes" | while read -r index_name; do
            if [ -n "$index_name" ]; then
                echo "Deleting remaining test index: $index_name"
                $PC_BINARY index delete "$index_name" >/dev/null 2>&1 || echo "Failed to delete remaining index: $index_name"
            fi
        done
    fi
    
    echo "Cleanup completed"
}

# Helper function to generate unique index names
generate_index_name() {
    local suffix="${1:-}"
    # Ensure name is under 45 characters (Pinecone limit) and ends with alphanumeric
    local timestamp=$(date +%s)
    local max_suffix_length=$((45 - ${#TEST_INDEX_PREFIX} - ${#timestamp} - 2))
    if [ ${#suffix} -gt $max_suffix_length ]; then
        suffix="${suffix:0:$max_suffix_length}"
    fi
    # Remove leading and trailing hyphens and ensure it ends with alphanumeric
    suffix=$(echo "$suffix" | sed 's/^-*//' | sed 's/-*$//')
    if [ -z "$suffix" ]; then
        suffix="test"
    fi
    # Ensure we don't have double hyphens
    echo "${TEST_INDEX_PREFIX}${timestamp}-${suffix}"
}

# Helper function to track created indexes
track_index() {
    local index_name="$1"
    TEST_CREATED_INDEXES+=("$index_name")
    echo "Tracking index for cleanup: $index_name"
}

# Helper function to wait for index to be ready
wait_for_index_ready() {
    local index_name="$1"
    local max_attempts=12  # 1 minute (12 attempts * 5 seconds)
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        local status
        local raw_output
        local json_output
        raw_output=$($PC_BINARY index describe "$index_name" --json 2>/dev/null || echo "")
        json_output=$(extract_json_from_output "$raw_output" 2>/dev/null || echo "")
        status=$(echo "$json_output" | jq -r '.status.state' 2>/dev/null || echo "unknown")
        
        echo "Index $index_name status: $status (attempt $attempt/$max_attempts)" >&2
        
        if [ "$status" = "Ready" ]; then
            echo "Index $index_name is ready!" >&2
            return 0
        elif [ "$status" = "Failed" ]; then
            echo "Index $index_name failed to initialize" >&2
            return 1
        elif [ "$status" = "Initializing" ]; then
            echo "Index $index_name is still initializing, waiting..." >&2
        elif [ "$status" = "unknown" ] || [ "$status" = "null" ]; then
            echo "Index $index_name status unknown, retrying... (attempt $attempt/$max_attempts)" >&2
        else
            echo "Index $index_name has status: $status, waiting..." >&2
        fi
        
        sleep 5
        attempt=$((attempt + 1))
    done
    
    echo "Index $index_name did not reach Ready state within 1 minute timeout" >&2
    return 1
}

# Helper function to wait for collection to be ready
wait_for_collection_ready() {
    local collection_name="$1"
    local max_attempts=6  # 1 minute (6 attempts * 10 seconds)
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        local status
        local raw_output
        local json_output
        raw_output=$($PC_BINARY collection describe "$collection_name" --json 2>/dev/null || echo "")
        json_output=$(extract_json_from_output "$raw_output" 2>/dev/null || echo "")
        status=$(echo "$json_output" | jq -r '.status.state' 2>/dev/null || echo "unknown")
        
        echo "Collection $collection_name status: $status (attempt $attempt/$max_attempts)" >&2
        
        if [ "$status" = "Ready" ]; then
            echo "Collection $collection_name is ready!" >&2
            return 0
        elif [ "$status" = "Failed" ]; then
            echo "Collection $collection_name failed to initialize" >&2
            return 1
        elif [ "$status" = "Initializing" ]; then
            echo "Collection $collection_name is still initializing, waiting..." >&2
        elif [ "$status" = "unknown" ] || [ "$status" = "null" ]; then
            echo "Collection $collection_name status unknown, retrying... (attempt $attempt/$max_attempts)" >&2
        else
            echo "Collection $collection_name has status: $status, waiting..." >&2
        fi
        
        sleep 10
        attempt=$((attempt + 1))
    done
    
    echo "Collection $collection_name did not reach Ready state within 1 minute timeout" >&2
    return 1
}

# Helper function to extract JSON from CLI output
extract_json_from_output() {
    local output="$1"
    # Find the line that starts with { and extract from there
    echo "$output" | awk '/^{/,0'
}

# Helper function to extract field from JSON output
extract_json_field() {
    local json="$1"
    local field="$2"
    # Remove any leading/trailing whitespace and newlines
    local result
    result=$(echo "$json" | jq -r ".$field" 2>/dev/null | tr -d '\n\r' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    echo "$result"
}

# Test basic serverless index creation with minimal flags
@test "create serverless index with minimal flags" {
    local index_name
    index_name=$(generate_index_name "-minimal")
    
    # Track the index for cleanup
    track_index "$index_name"
    
    # Create index
    run $PC_BINARY index create "$index_name" --yes
    [ "$status" -eq 0 ]
    
    # Wait for index to be ready
    wait_for_index_ready "$index_name"
    
    # Describe index and verify fields
    run $PC_BINARY index describe "$index_name" --json
    [ "$status" -eq 0 ]
    
    # Extract just the JSON part from the output
    local json_output
    json_output=$(extract_json_from_output "$output")
    [ "$(extract_json_field "$json_output" "name")" = "$index_name" ]
    [ "$(extract_json_field "$json_output" "spec.serverless.cloud")" = "aws" ]
    [ "$(extract_json_field "$json_output" "spec.serverless.region")" = "us-east-1" ]
    [ "$(extract_json_field "$json_output" "dimension")" = "1536" ]
    [ "$(extract_json_field "$json_output" "metric")" = "cosine" ]
    
    # Delete index
    run $PC_BINARY index delete "$index_name"
    [ "$status" -eq 0 ]
}

# Test serverless index with all possible flags
@test "create serverless index with all flags" {
    local index_name
    index_name=$(generate_index_name "-serverless-full")
    
    # Create index with all serverless flags
    run $PC_BINARY index create "$index_name" \
        --serverless \
        --dimension 768 \
        --metric euclidean \
        --cloud gcp \
        --region us-central1 \
        --vector_type dense \
        --deletion_protection enabled \
        --tags "env=test,type=serverless" \
        --yes
    [ "$status" -eq 0 ]
    
    # Wait for index to be ready
    wait_for_index_ready "$index_name"
    
    # Describe index and verify all fields
    run $PC_BINARY index describe "$index_name" --json
    [ "$status" -eq 0 ]
    
    local json_output="$output"
    [ "$(extract_json_field "$json_output" "name")" = "$index_name" ]
    [ "$(extract_json_field "$json_output" "spec.serverless.cloud")" = "gcp" ]
    [ "$(extract_json_field "$json_output" "spec.serverless.region")" = "us-central1" ]
    [ "$(extract_json_field "$json_output" "spec.dimension")" = "768" ]
    [ "$(extract_json_field "$json_output" "spec.metric")" = "euclidean" ]
    [ "$(extract_json_field "$json_output" "spec.serverless.vectorType")" = "dense" ]
    [ "$(extract_json_field "$json_output" "spec.deletionProtection")" = "enabled" ]
    
    # Verify tags
    local tags
    tags=$(extract_json_field "$json_output" "tags.env")
    [ "$tags" = "test" ]
    
    # Delete index
    run $PC_BINARY index delete "$index_name"
    [ "$status" -eq 0 ]
}

# Test pod index with minimal flags
@test "create pod index with minimal flags" {
    local index_name
    index_name=$(generate_index_name "-pod-minimal")
    
    # Track the index for cleanup
    track_index "$index_name"
    
    # Create pod index
    run $PC_BINARY index create "$index_name" \
        --pod \
        --yes
    [ "$status" -eq 0 ]
    
    # Wait for index to be ready
    wait_for_index_ready "$index_name"
    
    # Describe index and verify fields
    run $PC_BINARY index describe "$index_name" --json
    [ "$status" -eq 0 ]
    
    # Extract just the JSON part from the output
    local json_output
    json_output=$(extract_json_from_output "$output")
    [ "$(extract_json_field "$json_output" "name")" = "$index_name" ]
    [ "$(extract_json_field "$json_output" "spec.pod.environment")" = "us-east-1-aws" ]
    [ "$(extract_json_field "$json_output" "spec.pod.pod_type")" = "p1.x1" ]
    [ "$(extract_json_field "$json_output" "dimension")" = "1536" ]
    [ "$(extract_json_field "$json_output" "metric")" = "cosine" ]
    [ "$(extract_json_field "$json_output" "spec.pod.shard_count")" = "1" ]
    [ "$(extract_json_field "$json_output" "spec.pod.replicas")" = "1" ]
    
    # Delete index
    run $PC_BINARY index delete "$index_name"
    [ "$status" -eq 0 ]
}

# Test pod index with all possible flags
@test "create pod index with all flags" {
    local index_name
    index_name=$(generate_index_name "-pod-full")
    
    # Create pod index with all flags
    run $PC_BINARY index create "$index_name" \
        --pod \
        --dimension 1024 \
        --metric dotproduct \
        --environment us-west1-gcp \
        --pod_type p1.x2 \
        --shards 2 \
        --replicas 2 \
        --metadata_config "field1" \
        --metadata_config "field2" \
        --deletion_protection enabled \
        --tags "env=test,type=pod" \
        --yes
    [ "$status" -eq 0 ]
    
    # Wait for index to be ready
    wait_for_index_ready "$index_name"
    
    # Describe index and verify all fields
    run $PC_BINARY index describe "$index_name" --json
    [ "$status" -eq 0 ]
    
    local json_output="$output"
    [ "$(extract_json_field "$json_output" "name")" = "$index_name" ]
    [ "$(extract_json_field "$json_output" "spec.pod.environment")" = "us-west1-gcp" ]
    [ "$(extract_json_field "$json_output" "spec.pod.podType")" = "p1.x2" ]
    [ "$(extract_json_field "$json_output" "spec.dimension")" = "1024" ]
    [ "$(extract_json_field "$json_output" "spec.metric")" = "dotproduct" ]
    [ "$(extract_json_field "$json_output" "spec.pod.shards")" = "2" ]
    [ "$(extract_json_field "$json_output" "spec.pod.replicas")" = "2" ]
    [ "$(extract_json_field "$json_output" "spec.deletionProtection")" = "enabled" ]
    
    # Verify metadata config
    local metadata_config
    metadata_config=$(extract_json_field "$json_output" "spec.pod.metadataConfig.indexed[0]")
    [ "$metadata_config" = "field1" ]
    
    # Delete index
    run $PC_BINARY index delete "$index_name"
    [ "$status" -eq 0 ]
}

# Test integrated index with minimal flags
@test "create integrated index with minimal flags" {
    local index_name
    index_name=$(generate_index_name "-integrated-minimal")
    
    # Track the index for cleanup
    track_index "$index_name"
    
    # Create integrated index
    run $PC_BINARY index create "$index_name" \
        --integrated \
        --vector_type dense \
        --yes
    [ "$status" -eq 0 ]
    
    # Wait for index to be ready
    wait_for_index_ready "$index_name"
    
    # Describe index and verify fields
    run $PC_BINARY index describe "$index_name" --json
    [ "$status" -eq 0 ]
    
    # Extract just the JSON part from the output
    local json_output
    json_output=$(extract_json_from_output "$output")
    
    local actual_name
    actual_name=$(extract_json_field "$json_output" "name")
    echo "Expected name: $index_name"
    echo "Actual name: '$actual_name'"
    [ "$actual_name" = "$index_name" ]
    [ "$(extract_json_field "$json_output" "spec.serverless.cloud")" = "aws" ]
    [ "$(extract_json_field "$json_output" "spec.serverless.region")" = "us-east-1" ]
    # For integrated indexes, check embed dimension
    local actual_dimension
    actual_dimension=$(extract_json_field "$json_output" "embed.dimension")
    if [ -n "$actual_dimension" ]; then
        [ "$actual_dimension" = "1024" ]
    else
        [ "$(extract_json_field "$json_output" "dimension")" = "1536" ]
    fi
    [ "$(extract_json_field "$json_output" "metric")" = "cosine" ]
    
    # Verify embed configuration
    local model
    model=$(extract_json_field "$json_output" "embed.model")
    [ "$model" = "multilingual-e5-large" ]
    
    # Delete index
    run $PC_BINARY index delete "$index_name"
    [ "$status" -eq 0 ]
}

# Test integrated index with all possible flags
@test "create integrated index with all flags" {
    local index_name
    index_name=$(generate_index_name "-integrated-full")
    
    # Create integrated index with all flags
    run $PC_BINARY index create "$index_name" \
        --integrated \
        --vector_type dense \
        --dimension 768 \
        --metric euclidean \
        --cloud azure \
        --region eastus2 \
        --model multilingual-e5-large \
        --field_map "text=chunk_text" \
        --field_map "title=chunk_title" \
        --read_parameters "input_type=query" \
        --write_parameters "input_type=passage" \
        --deletion_protection enabled \
        --tags "env=test,type=integrated" \
        --yes
    [ "$status" -eq 0 ]
    
    # Wait for index to be ready
    wait_for_index_ready "$index_name"
    
    # Describe index and verify all fields
    run $PC_BINARY index describe "$index_name" --json
    [ "$status" -eq 0 ]
    
    local json_output="$output"
    [ "$(extract_json_field "$json_output" "name")" = "$index_name" ]
    [ "$(extract_json_field "$json_output" "spec.serverless.cloud")" = "azure" ]
    [ "$(extract_json_field "$json_output" "spec.serverless.region")" = "eastus2" ]
    [ "$(extract_json_field "$json_output" "spec.dimension")" = "768" ]
    [ "$(extract_json_field "$json_output" "spec.metric")" = "euclidean" ]
    [ "$(extract_json_field "$json_output" "spec.deletionProtection")" = "enabled" ]
    
    # Verify embed configuration
    local model
    model=$(extract_json_field "$json_output" "spec.embed.model")
    [ "$model" = "multilingual-e5-large" ]
    
    # Verify field map
    local field_map
    field_map=$(extract_json_field "$json_output" "spec.embed.fieldMap.text")
    [ "$field_map" = "chunk_text" ]
    
    # Delete index
    run $PC_BINARY index delete "$index_name"
    [ "$status" -eq 0 ]
}

# Test all vector types for serverless
@test "create serverless index with different vector types" {
    local vector_types=("dense" "sparse")
    
    for vector_type in "${vector_types[@]}"; do
        local index_name
        index_name=$(generate_index_name "-serverless-${vector_type}")
        
        # Create index with specific vector type
        run $PC_BINARY index create "$index_name" \
            --serverless \
            --vector_type "$vector_type" \
            --yes
        [ "$status" -eq 0 ]
        
        # Wait for index to be ready
        wait_for_index_ready "$index_name"
        
        # Describe index and verify vector type
        run $PC_BINARY index describe "$index_name" --json
        [ "$status" -eq 0 ]
        
        local json_output="$output"
        local actual_vector_type
        actual_vector_type=$(extract_json_field "$json_output" "spec.serverless.vectorType")
        [ "$actual_vector_type" = "$vector_type" ]
        
        # Delete index
        run $PC_BINARY index delete "$index_name"
        [ "$status" -eq 0 ]
    done
}

# Test all metrics
@test "create indexes with different metrics" {
    local metrics=("cosine" "euclidean" "dotproduct")
    
    for metric in "${metrics[@]}"; do
        local index_name
        index_name=$(generate_index_name "-metric-${metric}")
        
        # Create index with specific metric
        run $PC_BINARY index create "$index_name" \
            --serverless \
            --metric "$metric" \
            --yes
        [ "$status" -eq 0 ]
        
        # Wait for index to be ready
        wait_for_index_ready "$index_name"
        
        # Describe index and verify metric
        run $PC_BINARY index describe "$index_name" --json
        [ "$status" -eq 0 ]
        
        local json_output="$output"
        local actual_metric
        actual_metric=$(extract_json_field "$json_output" "spec.metric")
        [ "$actual_metric" = "$metric" ]
        
        # Delete index
        run $PC_BINARY index delete "$index_name"
        [ "$status" -eq 0 ]
    done
}

# Test all cloud providers
@test "create serverless indexes with different cloud providers" {
    local clouds=("aws" "gcp" "azure")
    
    for cloud in "${clouds[@]}"; do
        local index_name
        index_name=$(generate_index_name "-cloud-${cloud}")
        
        # Create index with specific cloud
        run $PC_BINARY index create "$index_name" \
            --serverless \
            --cloud "$cloud" \
            --yes
        [ "$status" -eq 0 ]
        
        # Wait for index to be ready
        wait_for_index_ready "$index_name"
        
        # Describe index and verify cloud
        run $PC_BINARY index describe "$index_name" --json
        [ "$status" -eq 0 ]
        
        local json_output="$output"
        local actual_cloud
        actual_cloud=$(extract_json_field "$json_output" "spec.serverless.cloud")
        [ "$actual_cloud" = "$cloud" ]
        
        # Delete index
        run $PC_BINARY index delete "$index_name"
        [ "$status" -eq 0 ]
    done
}

# Test different dimensions
@test "create indexes with different dimensions" {
    local dimensions=(512 768 1024 1536 2048)
    
    for dimension in "${dimensions[@]}"; do
        local index_name
        index_name=$(generate_index_name "-dimension-${dimension}")
        
        # Create index with specific dimension
        run $PC_BINARY index create "$index_name" \
            --serverless \
            --dimension "$dimension" \
            --yes
        [ "$status" -eq 0 ]
        
        # Wait for index to be ready
        wait_for_index_ready "$index_name"
        
        # Describe index and verify dimension
        run $PC_BINARY index describe "$index_name" --json
        [ "$status" -eq 0 ]
        
        local json_output="$output"
        local actual_dimension
        actual_dimension=$(extract_json_field "$json_output" "spec.dimension")
        [ "$actual_dimension" = "$dimension" ]
        
        # Delete index
        run $PC_BINARY index delete "$index_name"
        [ "$status" -eq 0 ]
    done
}

# Test different pod types
@test "create pod indexes with different pod types" {
    local pod_types=("p1.x1" "p1.x2" "p1.x4" "p1.x8" "s1.x1" "s1.x2" "s1.x4" "s1.x8")
    
    for pod_type in "${pod_types[@]}"; do
        local index_name
        index_name=$(generate_index_name "-podtype-${pod_type}")
        
        # Create index with specific pod type
        run $PC_BINARY index create "$index_name" \
            --pod \
            --pod_type "$pod_type" \
            --yes
        [ "$status" -eq 0 ]
        
        # Wait for index to be ready
        wait_for_index_ready "$index_name"
        
        # Describe index and verify pod type
        run $PC_BINARY index describe "$index_name" --json
        [ "$status" -eq 0 ]
        
        local json_output="$output"
        local actual_pod_type
        actual_pod_type=$(extract_json_field "$json_output" "spec.pod.podType")
        [ "$actual_pod_type" = "$pod_type" ]
        
        # Delete index
        run $PC_BINARY index delete "$index_name"
        [ "$status" -eq 0 ]
    done
}

# Test different shards and replicas combinations
@test "create pod indexes with different shards and replicas" {
    local shards_replicas=("1:1" "2:1" "1:2" "2:2" "4:2" "2:4")
    
    for combo in "${shards_replicas[@]}"; do
        local shards replicas
        IFS=':' read -r shards replicas <<< "$combo"
        
        local index_name
        index_name=$(generate_index_name "-shards-${shards}-replicas-${replicas}")
        
        # Create index with specific shards and replicas
        run $PC_BINARY index create "$index_name" \
            --pod \
            --shards "$shards" \
            --replicas "$replicas" \
            --yes
        [ "$status" -eq 0 ]
        
        # Wait for index to be ready
        wait_for_index_ready "$index_name"
        
        # Describe index and verify shards and replicas
        run $PC_BINARY index describe "$index_name" --json
        [ "$status" -eq 0 ]
        
        local json_output="$output"
        local actual_shards actual_replicas
        actual_shards=$(extract_json_field "$json_output" "spec.pod.shards")
        actual_replicas=$(extract_json_field "$json_output" "spec.pod.replicas")
        [ "$actual_shards" = "$shards" ]
        [ "$actual_replicas" = "$replicas" ]
        
        # Delete index
        run $PC_BINARY index delete "$index_name"
        [ "$status" -eq 0 ]
    done
}

# Test different embedding models for integrated indexes
@test "create integrated indexes with different models" {
    local models=("multilingual-e5-large" "llama-text-embed-v2")
    
    for model in "${models[@]}"; do
        local index_name
        index_name=$(generate_index_name "-model-${model}")
        
        # Create index with specific model
        run $PC_BINARY index create "$index_name" \
            --integrated \
            --vector_type dense \
            --model "$model" \
            --yes
        [ "$status" -eq 0 ]
        
        # Wait for index to be ready
        wait_for_index_ready "$index_name"
        
        # Describe index and verify model
        run $PC_BINARY index describe "$index_name" --json
        [ "$status" -eq 0 ]
        
        local json_output="$output"
        local actual_model
        actual_model=$(extract_json_field "$json_output" "embed.model")
        [ "$actual_model" = "$model" ]
        
        # Delete index
        run $PC_BINARY index delete "$index_name"
        [ "$status" -eq 0 ]
    done
}

# Test deletion protection
@test "create indexes with deletion protection" {
    local protection_values=("enabled" "disabled")
    
    for protection in "${protection_values[@]}"; do
        local index_name
        index_name=$(generate_index_name "-protection-${protection}")
        
        # Create index with specific deletion protection
        run $PC_BINARY index create "$index_name" \
            --serverless \
            --deletion_protection "$protection" \
            --yes
        [ "$status" -eq 0 ]
        
        # Wait for index to be ready
        wait_for_index_ready "$index_name"
        
        # Describe index and verify deletion protection
        run $PC_BINARY index describe "$index_name" --json
        [ "$status" -eq 0 ]
        
        local json_output="$output"
        local actual_protection
        actual_protection=$(extract_json_field "$json_output" "spec.deletionProtection")
        [ "$actual_protection" = "$protection" ]
        
        # Delete index
        run $PC_BINARY index delete "$index_name"
        [ "$status" -eq 0 ]
    done
}

# Test tags functionality
@test "create indexes with tags" {
    local index_name
    index_name=$(generate_index_name "-with-tags")
    
    # Create index with tags
    run $PC_BINARY index create "$index_name" \
        --serverless \
        --tags "env=test" \
        --tags "team=qa" \
        --tags "version=1.0" \
        --yes
    [ "$status" -eq 0 ]
    
    # Wait for index to be ready
    wait_for_index_ready "$index_name"
    
    # Describe index and verify tags
    run $PC_BINARY index describe "$index_name" --json
    [ "$status" -eq 0 ]
    
    local json_output="$output"
    local env_tag team_tag version_tag
    env_tag=$(extract_json_field "$json_output" "tags.env")
    team_tag=$(extract_json_field "$json_output" "tags.team")
    version_tag=$(extract_json_field "$json_output" "tags.version")
    [ "$env_tag" = "test" ]
    [ "$team_tag" = "qa" ]
    [ "$version_tag" = "1.0" ]
    
    # Delete index
    run $PC_BINARY index delete "$index_name"
    [ "$status" -eq 0 ]
}

# Test source collection functionality
@test "create indexes with source collection" {
    # First create a collection to use as source
    local collection_name
    collection_name=$(generate_index_name "-collection")
    
    # Create collection
    run $PC_BINARY collection create "$collection_name" \
        --dimension 1536 \
        --metric cosine \
        --yes
    [ "$status" -eq 0 ]
    
    # Wait for collection to be ready
    wait_for_collection_ready "$collection_name"
    
    local index_name
    index_name=$(generate_index_name "-from-collection")
    
    # Create index from collection
    run $PC_BINARY index create "$index_name" \
        --serverless \
        --source_collection "$collection_name" \
        --yes
    [ "$status" -eq 0 ]
    
    # Wait for index to be ready
    wait_for_index_ready "$index_name"
    
    # Describe index and verify source collection
    run $PC_BINARY index describe "$index_name" --json
    [ "$status" -eq 0 ]
    
    local json_output="$output"
    local source_collection
    source_collection=$(extract_json_field "$json_output" "spec.serverless.sourceCollection")
    [ "$source_collection" = "$collection_name" ]
    
    # Delete index and collection
    run $PC_BINARY index delete "$index_name"
    [ "$status" -eq 0 ]
    
    run $PC_BINARY collection delete "$collection_name" --yes
    [ "$status" -eq 0 ]
}

# Test invalid flag combinations
@test "create index with invalid flag combinations" {
    # Test serverless with pod-specific flags
    run $PC_BINARY index create "test-invalid" \
        --serverless \
        --environment "us-east-1-aws" \
        --pod_type "p1.x1" \
        --shards 2 \
        --replicas 2 \
        --yes
    [ "$status" -ne 0 ]
    [[ "$output" == *"cannot be used with serverless indexes"* ]]
    
    # Test pod with serverless-specific flags
    run $PC_BINARY index create "test-invalid" \
        --pod \
        --cloud "aws" \
        --region "us-east-1" \
        --vector_type "dense" \
        --yes
    [ "$status" -ne 0 ]
    [[ "$output" == *"cannot be used with pod indexes"* ]]
    
    # Test integrated with pod-specific flags
    run $PC_BINARY index create "test-invalid" \
        --integrated \
        --environment "us-east-1-aws" \
        --pod_type "p1.x1" \
        --shards 2 \
        --replicas 2 \
        --yes
    [ "$status" -ne 0 ]
    [[ "$output" == *"cannot be used with integrated indexes"* ]]
}

# Test invalid values
@test "create index with invalid values" {
    # Test invalid metric
    run $PC_BINARY index create "test-invalid" \
        --serverless \
        --metric "invalid-metric" \
        --yes
    [ "$status" -ne 0 ]
    
    # Test invalid cloud
    run $PC_BINARY index create "test-invalid" \
        --serverless \
        --cloud "invalid-cloud" \
        --yes
    [ "$status" -ne 0 ]
    
    # Test invalid dimension (negative)
    run $PC_BINARY index create "test-invalid" \
        --serverless \
        --dimension -1 \
        --yes
    [ "$status" -ne 0 ]
    
    # Test invalid dimension (too large)
    run $PC_BINARY index create "test-invalid" \
        --serverless \
        --dimension 100000 \
        --yes
    [ "$status" -ne 0 ]
    
    # Test invalid pod type
    run $PC_BINARY index create "test-invalid" \
        --pod \
        --pod_type "invalid-pod-type" \
        --yes
    [ "$status" -ne 0 ]
    
    # Test invalid shards (negative)
    run $PC_BINARY index create "test-invalid" \
        --pod \
        --shards -1 \
        --yes
    [ "$status" -ne 0 ]
    
    # Test invalid replicas (negative)
    run $PC_BINARY index create "test-invalid" \
        --pod \
        --replicas -1 \
        --yes
    [ "$status" -ne 0 ]
}

# Test missing required values
@test "create index with missing required values" {
    # Test missing name
    run $PC_BINARY index create
    [ "$status" -ne 0 ]
    
    # Test missing cloud for serverless
    run $PC_BINARY index create "test-missing" \
        --serverless \
        --region "us-east-1" \
        --yes
    [ "$status" -ne 0 ]
    
    # Test missing region for serverless
    run $PC_BINARY index create "test-missing" \
        --serverless \
        --cloud "aws" \
        --yes
    [ "$status" -ne 0 ]
    
    # Test missing environment for pod
    run $PC_BINARY index create "test-missing" \
        --pod \
        --pod_type "p1.x1" \
        --yes
    [ "$status" -ne 0 ]
    
    # Test missing pod type for pod
    run $PC_BINARY index create "test-missing" \
        --pod \
        --environment "us-east-1-aws" \
        --yes
    [ "$status" -ne 0 ]
    
    # Test missing model for integrated
    run $PC_BINARY index create "test-missing" \
        --integrated \
        --cloud "aws" \
        --region "us-east-1" \
        --yes
    [ "$status" -ne 0 ]
}

# Test describe non-existent index
@test "describe non-existent index" {
    run $PC_BINARY index describe "non-existent-index-$(date +%s)"
    [ "$status" -ne 0 ]
    [[ "$output" == *"does not exist"* ]]
}

# Test delete non-existent index
@test "delete non-existent index" {
    run $PC_BINARY index delete "non-existent-index-$(date +%s)" --yes
    [ "$status" -ne 0 ]
    [[ "$output" == *"does not exist"* ]]
}

# Test JSON output format
@test "create and describe index with JSON output" {
    local index_name
    index_name=$(generate_index_name "-json-test")
    
    # Create index
    run $PC_BINARY index create "$index_name" \
        --serverless \
        --json \
        --yes
    [ "$status" -eq 0 ]
    
    # Verify JSON output contains expected fields
    echo "$output" | jq -e '.name' >/dev/null
    echo "$output" | jq -e '.spec.serverless.cloud' >/dev/null
    echo "$output" | jq -e '.spec.serverless.region' >/dev/null
    
    # Wait for index to be ready
    wait_for_index_ready "$index_name"
    
    # Describe index with JSON output
    run $PC_BINARY index describe "$index_name" --json
    [ "$status" -eq 0 ]
    
    # Verify JSON output contains expected fields
    echo "$output" | jq -e '.name' >/dev/null
    echo "$output" | jq -e '.status.state' >/dev/null
    echo "$output" | jq -e '.spec.dimension' >/dev/null
    echo "$output" | jq -e '.spec.metric' >/dev/null
    
    # Delete index
    run $PC_BINARY index delete "$index_name"
    [ "$status" -eq 0 ]
}

# Test interactive mode (basic test)
@test "create index in interactive mode" {
    # This test is limited since we can't easily simulate user input
    # We'll just test that the command doesn't fail when no name is provided
    run timeout 5s $PC_BINARY index create
    # Should either timeout (interactive mode waiting for input) or fail gracefully
    [ "$status" -eq 124 ] || [ "$status" -ne 0 ]
}

# Test confirmation prompt (when --yes is not used)
@test "create index without confirmation" {
    # This test is limited since we can't easily simulate user input
    # We'll test that the command works with --yes flag
    local index_name
    index_name=$(generate_index_name "-no-confirm")
    
    run $PC_BINARY index create "$index_name" --serverless --yes
    [ "$status" -eq 0 ]
    
    # Wait for index to be ready
    wait_for_index_ready "$index_name"
    
    # Delete index
    run $PC_BINARY index delete "$index_name"
    [ "$status" -eq 0 ]
}

# Test concurrent index operations
@test "create multiple indexes concurrently" {
    local index_names=()
    local pids=()
    
    # Create 3 indexes concurrently
    for i in {1..3}; do
        local index_name
        index_name=$(generate_index_name "-concurrent-${i}")
        index_names+=("$index_name")
        
        $PC_BINARY index create "$index_name" --serverless --yes &
        pids+=($!)
    done
    
    # Wait for all processes to complete
    for pid in "${pids[@]}"; do
        wait "$pid"
    done
    
    # Verify all indexes were created
    for index_name in "${index_names[@]}"; do
        run $PC_BINARY index describe "$index_name" --json
        [ "$status" -eq 0 ]
        
        # Delete index
        run $PC_BINARY index delete "$index_name"
        [ "$status" -eq 0 ]
    done
}

# Test index lifecycle with different states
@test "test index lifecycle states" {
    local index_name
    index_name=$(generate_index_name "-lifecycle")
    
    # Create index
    run $PC_BINARY index create "$index_name" --serverless --yes
    [ "$status" -eq 0 ]
    
    # Check initial state (should be Initializing)
    run $PC_BINARY index describe "$index_name" --json
    [ "$status" -eq 0 ]
    
    local json_output="$output"
    local initial_state
    initial_state=$(extract_json_field "$json_output" "status.state")
    [ "$initial_state" = "Initializing" ] || [ "$initial_state" = "ScalingUp" ] || [ "$initial_state" = "Ready" ]
    
    # Wait for index to be ready
    wait_for_index_ready "$index_name"
    
    # Check final state (should be Ready)
    run $PC_BINARY index describe "$index_name" --json
    [ "$status" -eq 0 ]
    
    json_output="$output"
    local final_state
    final_state=$(extract_json_field "$json_output" "status.state")
    [ "$final_state" = "Ready" ]
    
    # Delete index
    run $PC_BINARY index delete "$index_name"
    [ "$status" -eq 0 ]
}

# Test error handling for network issues
@test "handle network errors gracefully" {
    # This test would require mocking network failures
    # For now, we'll test that the CLI handles invalid API responses gracefully
    
    # Test with invalid API key (if possible)
    local original_api_key
    original_api_key=$($PC_BINARY config get api-key 2>/dev/null || echo "")
    
    if [ -n "$original_api_key" ]; then
        # Temporarily set invalid API key
        $PC_BINARY config set api-key "invalid-key" >/dev/null 2>&1 || true
        
        # Try to create index (should fail gracefully)
        run $PC_BINARY index create "test-network-error" --serverless --yes
        [ "$status" -ne 0 ]
        
        # Restore original API key
        $PC_BINARY config set api-key "$original_api_key" >/dev/null 2>&1 || true
    fi
}

# Test performance with large configurations
@test "create index with large configuration" {
    local index_name
    index_name=$(generate_index_name "-large-config")
    
    # Create index with many metadata config fields
    local metadata_configs=()
    for i in {1..10}; do
        metadata_configs+=("field${i}")
    done
    
    run $PC_BINARY index create "$index_name" \
        --pod \
        --metadata_config "${metadata_configs[@]}" \
        --tags "field1=value1" \
        --tags "field2=value2" \
        --tags "field3=value3" \
        --tags "field4=value4" \
        --tags "field5=value5" \
        --yes
    [ "$status" -eq 0 ]
    
    # Wait for index to be ready
    wait_for_index_ready "$index_name"
    
    # Verify the configuration was applied correctly
    run $PC_BINARY index describe "$index_name" --json
    [ "$status" -eq 0 ]
    
    local json_output="$output"
    # Verify metadata config fields
    for i in {1..10}; do
        local field_value
        field_value=$(extract_json_field "$json_output" "spec.pod.metadataConfig.indexed[$((i-1))]")
        [ "$field_value" = "field${i}" ]
    done
    
    # Delete index
    run $PC_BINARY index delete "$index_name"
    [ "$status" -eq 0 ]
}

# Test edge cases
@test "test edge cases" {
    # Test very long index name
    local long_name
    long_name="$(printf 'a%.0s' {1..50})-$(date +%s)"
    
    run $PC_BINARY index create "$long_name" --serverless --yes
    [ "$status" -eq 0 ]
    
    # Wait for index to be ready
    wait_for_index_ready "$long_name"
    
    # Verify index was created
    run $PC_BINARY index describe "$long_name" --json
    [ "$status" -eq 0 ]
    
    # Delete index
    run $PC_BINARY index delete "$long_name" --yes
    [ "$status" -eq 0 ]
    
    # Test special characters in tags
    local index_name
    index_name=$(generate_index_name "-special-chars")
    
    run $PC_BINARY index create "$index_name" \
        --serverless \
        --tags "special=value-with-dashes" \
        --tags "underscore=value_with_underscores" \
        --tags "numbers=value123" \
        --yes
    [ "$status" -eq 0 ]
    
    # Wait for index to be ready
    wait_for_index_ready "$index_name"
    
    # Verify tags were applied
    run $PC_BINARY index describe "$index_name" --json
    [ "$status" -eq 0 ]
    
    local json_output="$output"
    local special_tag underscore_tag numbers_tag
    special_tag=$(extract_json_field "$json_output" "tags.special")
    underscore_tag=$(extract_json_field "$json_output" "tags.underscore")
    numbers_tag=$(extract_json_field "$json_output" "tags.numbers")
    [ "$special_tag" = "value-with-dashes" ]
    [ "$underscore_tag" = "value_with_underscores" ]
    [ "$numbers_tag" = "value123" ]
    
    # Delete index
    run $PC_BINARY index delete "$index_name"
    [ "$status" -eq 0 ]
} 