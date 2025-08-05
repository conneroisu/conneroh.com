#!/bin/sh

# SVG Internal Links Removal Script (Shell version)
# Removes internal ID references from SVG files to prevent conflicts

set -eu

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
DRY_RUN=false
BACKUP=false
VERBOSE=false
DIRECTORY="."

# Target files that need processing (space-separated list)
TARGET_FILES="${TARGET_FILE:-apache-original.svg clion-original.svg d3js-original.svg eclipse-original.svg dropwizard-original.svg gentoo-original.svg gimp-original.svg goland-original.svg json-original.svg jira-original.svg moodle-original.svg nextjs-original.svg ocaml-original.svg poetry-original.svg prolog-original.svg rollup-original.svg ruby-original.svg webstorm-original.svg xcode-original.svg maven-original.svg renpy-original.svg vscode-original.svg}"

# Counters
FILES_PROCESSED=0
FILES_MODIFIED=0
TOTAL_CHANGES=0

# Usage information
usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Remove internal links from SVG files to prevent ID conflicts.

OPTIONS:
    -d, --directory DIR     Directory containing SVG files (default: current)
    -n, --dry-run          Show what would be changed without modifying files
    -b, --backup           Create .bak files before making changes
    -v, --verbose          Show detailed processing information
    -h, --help             Show this help message

EXAMPLES:
    $0                     # Process current directory
    $0 --dry-run           # Show what would be changed
    $0 --backup --verbose  # Create backups and show details

EOF
}

# Logging function
log() {
    local level=$1
    shift
    local message="$*"
    
    if [ "$VERBOSE" = true ] || [ "$level" = "INFO" ]; then
        case $level in
            "ERROR")   echo -e "${RED}[ERROR]${NC} $message" >&2 ;;
            "WARNING") echo -e "${YELLOW}[WARNING]${NC} $message" ;;
            "SUCCESS") echo -e "${GREEN}[SUCCESS]${NC} $message" ;;
            "INFO")    echo -e "${BLUE}[INFO]${NC} $message" ;;
            *)         echo -e "$message" ;;
        esac
    fi
}

# Check if a file has internal references
has_internal_references() {
    local file="$1"
    
    # Check for xlink:href="#..." patterns
    if grep -q 'xlink:href="#[^"]*"' "$file"; then
        return 0
    fi
    
    # Check for href="#..." patterns
    if grep -q 'href="#[^"]*"' "$file"; then
        return 0
    fi
    
    # Check for url(#...) patterns
    if grep -q 'url(#[^)]*)' "$file"; then
        return 0
    fi
    
    return 1
}

# Process a single SVG file
process_svg_file() {
    local file_path="$1"
    local filename=$(basename "$file_path")
    local changes_made=0
    
    log "DEBUG" "Processing: $filename"
    
    # Check if file has internal references
    if ! has_internal_references "$file_path"; then
        log "DEBUG" "No internal references found in $filename"
        return 0
    fi
    
    if [ "$DRY_RUN" = true ]; then
        log "INFO" "Would process $filename (contains internal references)"
        return 0
    fi
    
    # Create backup if requested
    if [ "$BACKUP" = true ]; then
        cp "$file_path" "$file_path.bak"
        log "DEBUG" "Created backup: $file_path.bak"
    fi
    
    # Create temporary file for processing
    local temp_file=$(mktemp)
    cp "$file_path" "$temp_file"
    
    # Remove xlink:href="#..." attributes
    if sed -i 's/xlink:href="#[^"]*"//g' "$temp_file"; then
        local xlink_changes=$(diff -u "$file_path" "$temp_file" | grep -c '^-.*xlink:href=' || true)
        if [ $xlink_changes -gt 0 ]; then
            changes_made=$((changes_made + xlink_changes))
            log "DEBUG" "Removed $xlink_changes xlink:href references from $filename"
        fi
    fi
    
    # Remove href="#..." attributes (but preserve external links)
    if sed -i 's/href="#[^"]*"//g' "$temp_file"; then
        local href_changes=$(diff -u "$file_path" "$temp_file" | grep -c '^-.*href="#' || true)
        if [ $href_changes -gt 0 ]; then
            changes_made=$((changes_made + href_changes))
            log "DEBUG" "Removed $href_changes href references from $filename"
        fi
    fi
    
    # Replace url(#...) with fallback colors in fill attributes
    if sed -i 's/fill="url(#[^"]*)"/ fill="#666666"/g' "$temp_file"; then
        local fill_changes=$(diff -u "$file_path" "$temp_file" | grep -c '^-.*fill="url(' || true)
        if [ $fill_changes -gt 0 ]; then
            changes_made=$((changes_made + fill_changes))
            log "DEBUG" "Replaced $fill_changes fill url() references with fallback color in $filename"
        fi
    fi
    
    # Remove url(#...) from other attributes (mask, filter, etc.)
    if sed -i 's/mask="url(#[^"]*)"//g; s/filter="url(#[^"]*)"//g; s/clip-path="url(#[^"]*)"//g' "$temp_file"; then
        local other_changes=$(diff -u "$file_path" "$temp_file" | grep -c '^-.*="url(' || true)
        if [ $other_changes -gt 0 ]; then
            changes_made=$((changes_made + other_changes))
            log "DEBUG" "Removed $other_changes other url() references from $filename"
        fi
    fi
    
    # Remove empty definition elements that are no longer referenced
    # This is a simplified approach - remove common definition elements with single-letter IDs
    sed -i '/<linearGradient id="[a-z]"[^>]*>/,/<\/linearGradient>/d' "$temp_file"
    sed -i '/<mask id="[a-z]"[^>]*>/,/<\/mask>/d' "$temp_file"
    sed -i '/<filter id="[a-z]"[^>]*>/,/<\/filter>/d' "$temp_file"
    sed -i '/<clipPath id="[a-z]"[^>]*>/,/<\/clipPath>/d' "$temp_file"
    
    # Only update the file if changes were made
    if [ $changes_made -gt 0 ]; then
        mv "$temp_file" "$file_path"
        FILES_MODIFIED=$((FILES_MODIFIED + 1))
        TOTAL_CHANGES=$((TOTAL_CHANGES + changes_made))
        log "SUCCESS" "✓ $filename: Removed $changes_made internal references"
    else
        rm "$temp_file"
        log "DEBUG" "No changes needed for $filename"
    fi
    
    return 0
}

# Process all target files in directory
process_directory() {
    local found_files=0
    
    log "INFO" "Searching for target SVG files in: $DIRECTORY"
    
    for target_file in $TARGET_FILES; do
        local file_path="$DIRECTORY/$target_file"
        
        if [ -f "$file_path" ]; then
            found_files=$((found_files + 1))
            FILES_PROCESSED=$((FILES_PROCESSED + 1))
            
            if process_svg_file "$file_path"; then
                log "DEBUG" "Successfully processed $target_file"
            else
                log "ERROR" "Failed to process $target_file"
            fi
        else
            log "DEBUG" "Target file not found: $target_file"
        fi
    done
    
    if [ $found_files -eq 0 ]; then
        log "WARNING" "No target SVG files found in directory: $DIRECTORY"
        return 1
    fi
    
    return 0
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--directory)
                DIRECTORY="$2"
                shift 2
                ;;
            -n|--dry-run)
                DRY_RUN=true
                shift
                ;;
            -b|--backup)
                BACKUP=true
                shift
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                echo "Unknown option: $1" >&2
                usage >&2
                exit 1
                ;;
        esac
    done
}

# Main function
main() {
    parse_args "$@"
    
    # Verify directory exists
    if [ ! -d "$DIRECTORY" ]; then
        log "ERROR" "Directory does not exist: $DIRECTORY"
        exit 1
    fi
    
    log "INFO" "SVG Internal Links Removal Script"
    log "INFO" "Directory: $DIRECTORY"
    
    if [ "$DRY_RUN" = true ]; then
        log "INFO" "DRY RUN MODE - No files will be modified"
    fi
    
    if [ "$BACKUP" = true ] && [ "$DRY_RUN" = false ]; then
        log "INFO" "Backup mode enabled - .bak files will be created"
    fi
    
    echo "----------------------------------------"
    
    # Process the directory
    if process_directory; then
        echo "----------------------------------------"
        log "INFO" "Processing complete!"
        log "INFO" "Files processed: $FILES_PROCESSED"
        
        if [[ "$DRY_RUN" == false ]]; then
            log "INFO" "Files modified: $FILES_MODIFIED"
            log "INFO" "Total changes made: $TOTAL_CHANGES"
        fi
    else
        exit 1
    fi
}

# Run main function with all arguments
main "$@"