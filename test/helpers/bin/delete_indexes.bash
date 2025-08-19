#!/bin/bash

# Source the indexes.bash file to get the get_machine_id function
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../custom/indexes.bash"

# Determine the pcdev binary location based on script location
# The script is in test/helpers/bin/, so pcdev should be in the project root (../../)
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
PCDEV_BIN="$PROJECT_ROOT/pcdev"

# Check if pcdev binary exists
if [ ! -f "$PCDEV_BIN" ]; then
    echo "Error: pcdev binary not found at $PCDEV_BIN"
    echo "Make sure you're running this script from the correct location and the binary is built."
    exit 1
fi

# Default behavior: only delete indexes created from this machine
DELETE_ALL=false
REMOVE_PROTECTION=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --all|-a)
            DELETE_ALL=true
            shift
            ;;
        --remove-protection|-r)
            REMOVE_PROTECTION=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --all, -a               Delete all test indexes (not just created from this machine)"
            echo "  --remove-protection, -r Remove deletion protection if deletion fails"
            echo "  --help, -h              Show this help message"
            echo ""
            echo "By default, only deletes indexes created from this machine."
            echo ""
            echo "Using pcdev binary at: $PCDEV_BIN"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Get the current machine ID
MACHINE_ID=$(get_machine_id)

if [ "$DELETE_ALL" = true ]; then
    echo "Deleting ALL test indexes..."
    indexes=$("$PCDEV_BIN" index list | grep "t-" | awk '{print $1}')
else
    echo "Deleting test indexes created from this machine (ID: $MACHINE_ID)..."
    indexes=$("$PCDEV_BIN" index list | grep "t-$MACHINE_ID" | awk '{print $1}')
fi

# Check if there are any indexes to delete
if [ -z "$indexes" ]; then
    if [ "$DELETE_ALL" = true ]; then
        echo "No test indexes found to delete."
    else
        echo "No test indexes found created from this machine (ID: $MACHINE_ID)."
    fi
    exit 0
fi

# Delete the indexes
echo "Found $(echo "$indexes" | wc -l) index(es) to delete:"
echo "$indexes" | sed 's/^/  /'

echo ""
echo "Proceeding with deletion..."

# Function to delete an index with protection handling
delete_index_with_protection_handling() {
    local index_name="$1"
    local max_attempts=2
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        echo "Deleting $index_name... (attempt $attempt/$max_attempts)"
        
        # Try to delete the index
        local delete_output
        local delete_exit_code
        
        delete_output=$("$PCDEV_BIN" index delete "$index_name" 2>&1)
        delete_exit_code=$?
        
        if [ $delete_exit_code -eq 0 ]; then
            echo "Successfully deleted $index_name"
            return 0
        fi
        
        # Check if deletion failed due to protection
        if echo "$delete_output" | grep -qi "deletion protection\|protected\|cannot delete.*protection"; then
            if [ "$REMOVE_PROTECTION" = true ] && [ $attempt -eq 1 ]; then
                echo "Deletion protection detected. Attempting to remove protection for $index_name..."
                
                # Try to disable deletion protection
                local disable_output
                local disable_exit_code
                
                disable_output=$("$PCDEV_BIN" index configure --name "$index_name" --deletion_protection disabled 2>&1)
                disable_exit_code=$?
                
                if [ $disable_exit_code -eq 0 ]; then
                    echo "Successfully disabled deletion protection for $index_name"
                    # Continue to next attempt (will try to delete again)
                else
                    echo "Failed to disable deletion protection for $index_name: $disable_output"
                    echo "Skipping $index_name"
                    return 1
                fi
            else
                echo "Deletion protection enabled for $index_name and --remove-protection not specified"
                echo "Skipping $index_name"
                return 1
            fi
        else
            # Some other error occurred
            echo "Failed to delete $index_name: $delete_output"
            return 1
        fi
        
        attempt=$((attempt + 1))
    done
    
    echo "Failed to delete $index_name after $max_attempts attempts"
    return 1
}

# Process each index
success_count=0
failed_count=0

while IFS= read -r index; do
    if [ -n "$index" ]; then
        if delete_index_with_protection_handling "$index"; then
            success_count=$((success_count + 1))
        else
            failed_count=$((failed_count + 1))
        fi
        echo ""  # Add spacing between operations
    fi
done <<< "$indexes"

echo "Deletion summary:"
echo "  Successfully deleted: $success_count"
echo "  Failed to delete: $failed_count"

if [ $failed_count -gt 0 ]; then
    echo ""
    echo "Note: Some indexes may have failed due to deletion protection."
    echo "Use --remove-protection flag to automatically handle protected indexes."
    exit 1
else
    echo "All indexes processed successfully."
    exit 0
fi