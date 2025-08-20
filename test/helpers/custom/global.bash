#!/usr/bin/env bash
add_executable_to_path() {
    PCDEV_DIR="$(cd "$(dirname "${BATS_ROOT[0]}")/.." && pwd)"
    PATH="$PCDEV_DIR:$PATH"
    
    # Set CLI to use the pcdev wrapper
    export CLI="pcdev"
}