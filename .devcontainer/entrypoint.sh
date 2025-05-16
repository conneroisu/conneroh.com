#!/usr/bin/env bash

# Enable verbose debugging
set -x                  # Print each command before execution
set -v                  # Print script lines as they are read
set -euo pipefail       # Strict error handling

# Define debug function with timestamps
debug() {
    local timestamp=$(date +"%Y-%m-%d %H:%M:%S.%3N")
    echo "[$timestamp] [DEBUG] [Line ${BASH_LINENO[0]}] $*" >&2
}

# Define error handler
error_handler() {
    local line_num=$1
    local command=$2
    local error_code=$3
    echo "----------------------------------------" >&2
    echo "❌ ERROR: Command '$command' failed with exit code $error_code at line $line_num" >&2
    echo "----------------------------------------" >&2
    env | sort >&2                  # Dump all environment variables
    echo "----------------------------------------" >&2
    caller 0 >&2                    # Show call stack
    echo "----------------------------------------" >&2
}

# Set trap to catch errors
trap 'error_handler ${LINENO} "$BASH_COMMAND" $?' ERR

# Set trap for script exit
trap 'debug "Script is exiting with status $?"' EXIT

# Print script start information
debug "Starting script $(basename "$0")"
debug "Running as user: $(whoami)"
debug "Current directory: $(pwd)"
debug "Bash version: $BASH_VERSION"

# Print original environment before loading new variables
debug "Original environment variables (before loading):"
env | sort | grep -v "^_" | while IFS= read -r line; do
    debug "  $line"
done

# 1) Load all of the exports that nix-shell generated:
debug "Attempting to load environment variables from /env-vars"
if [[ -f /env-vars ]]; then
    debug "File /env-vars exists. Loading..."
    # Source with verbose output
    (set -x; . /env-vars)
    debug "Environment variables loaded successfully"
else
    debug "❌ WARNING: File /env-vars does not exist!"
    ls -la / >&2  # List root directory to help debug
fi

# Check if important variables were loaded
debug "Checking for essential Nix-related variables:"
for var in NIX_PATH NIX_SSL_CERT_FILE NIX_PROFILES; do
    if [[ -v $var ]]; then
        debug "  ✓ $var is set to '${!var}'"
    else
        debug "  ✗ $var is NOT set!"
    fi
done

# 2) Print a message so you know it ran:
echo "✅ Nix env loaded from /env-vars"
debug "Nix environment loading completed"

# Print final environment after loading variables
debug "Final environment variables (after loading):"
env | sort | grep -v "^_" | while IFS= read -r line; do
    debug "  $line"
done

# 3) Keep the container alive by handing off to sleep
debug "Starting infinite sleep to keep container alive"
echo "💤 Container will now sleep indefinitely (press Ctrl+C to interrupt)"
set -x  # Make the exec command visible
exec sleep infinity
