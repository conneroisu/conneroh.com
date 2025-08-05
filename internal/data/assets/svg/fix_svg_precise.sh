#!/bin/sh

# Precise SVG Internal Links Removal Script
# Carefully removes internal ID references from SVG files

set -eu

# Configuration
DRY_RUN=false
BACKUP=false
VERBOSE=false
DIRECTORY="."

# Target files that need processing (space-separated list)
TARGET_FILES="${TARGET_FILE:-apache-original.svg clion-original.svg d3js-original.svg eclipse-original.svg dropwizard-original.svg gentoo-original.svg gimp-original.svg goland-original.svg json-original.svg jira-original.svg moodle-original.svg nextjs-original.svg ocaml-original.svg poetry-original.svg prolog-original.svg rollup-original.svg ruby-original.svg webstorm-original.svg xcode-original.svg maven-original.svg renpy-original.svg vscode-original.svg}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Counters
FILES_PROCESSED=0
FILES_MODIFIED=0
TOTAL_CHANGES=0

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

EOF
}

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

has_internal_refs() {
    local file="$1"
    
    # Check for various internal reference patterns
    if grep -q 'xlink:href="#[^"]*"' "$file" 2>/dev/null; then
        return 0
    fi
    if grep -q 'href="#[^"]*"' "$file" 2>/dev/null; then
        return 0
    fi
    if grep -q 'url(#[^)]*)' "$file" 2>/dev/null; then
        return 0
    fi
    
    return 1
}

process_svg_file() {
    local file_path="$1"
    local filename=$(basename "$file_path")
    local changes_made=0
    
    log "DEBUG" "Processing: $filename"
    
    if ! has_internal_refs "$file_path"; then
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
    
    # Create temporary file for safe processing
    local temp_file=$(mktemp)
    
    # Read the file and process it line by line to avoid destroying it
    while IFS= read -r line; do
        original_line="$line"
        
        # Remove xlink:href="#..." attributes
        line=$(echo "$line" | sed 's/ xlink:href="#[^"]*"//g')
        if [ "$line" != "$original_line" ]; then
            changes_made=$((changes_made + 1))
            log "DEBUG" "Removed xlink:href from line"
        fi
        
        original_line="$line"
        # Remove href="#..." attributes (but preserve external links starting with http)
        line=$(echo "$line" | sed 's/ href="#[^"]*"//g')
        if [ "$line" != "$original_line" ]; then
            changes_made=$((changes_made + 1))
            log "DEBUG" "Removed href from line"
        fi
        
        original_line="$line"
        # Replace fill="url(#...)" with fallback color
        line=$(echo "$line" | sed 's/fill="url(#[^"]*)"/ fill="#666666"/g')
        if [ "$line" != "$original_line" ]; then
            changes_made=$((changes_made + 1))
            log "DEBUG" "Replaced fill url() with fallback color"
        fi
        
        original_line="$line"
        # Remove other url(#...) patterns from mask, filter, etc.
        line=$(echo "$line" | sed 's/mask="url(#[^"]*)"//g; s/filter="url(#[^"]*)"//g; s/clip-path="url(#[^"]*)"//g')
        if [ "$line" != "$original_line" ]; then
            changes_made=$((changes_made + 1))
            log "DEBUG" "Removed other url() references"
        fi
        
        echo "$line" >> "$temp_file"
    done < "$file_path"
    
    # Remove definition blocks that are commonly problematic
    # This creates a new temp file without these blocks
    local temp_file2=$(mktemp)
    local in_gradient=false
    local in_mask=false
    local in_filter=false
    local in_clippath=false
    local skip_line=false
    
    while IFS= read -r line; do
        skip_line=false
        
        # Check if we're entering a problematic definition block
        if echo "$line" | grep -q '<linearGradient id="[a-z]"'; then
            in_gradient=true
            skip_line=true
        elif echo "$line" | grep -q '<mask id="[a-z]"'; then
            in_mask=true
            skip_line=true
        elif echo "$line" | grep -q '<filter id="[a-z]"'; then
            in_filter=true  
            skip_line=true
        elif echo "$line" | grep -q '<clipPath id="[a-z]"'; then
            in_clippath=true
            skip_line=true
        fi
        
        # Check if we're exiting a problematic definition block
        if [ "$in_gradient" = true ] && echo "$line" | grep -q '</linearGradient>'; then
            in_gradient=false
            skip_line=true
        elif [ "$in_mask" = true ] && echo "$line" | grep -q '</mask>'; then
            in_mask=false
            skip_line=true
        elif [ "$in_filter" = true ] && echo "$line" | grep -q '</filter>'; then
            in_filter=false
            skip_line=true
        elif [ "$in_clippath" = true ] && echo "$line" | grep -q '</clipPath>'; then
            in_clippath=false
            skip_line=true
        fi
        
        # Skip lines if we're inside a problematic block
        if [ "$in_gradient" = true ] || [ "$in_mask" = true ] || [ "$in_filter" = true ] || [ "$in_clippath" = true ]; then
            skip_line=true
        fi
        
        if [ "$skip_line" = false ]; then
            echo "$line" >> "$temp_file2"
        fi
    done < "$temp_file"
    
    # Only update the file if changes were made
    if [ $changes_made -gt 0 ] || [ "$(wc -l < "$temp_file2")" != "$(wc -l < "$file_path")" ]; then
        mv "$temp_file2" "$file_path"
        rm -f "$temp_file"
        FILES_MODIFIED=$((FILES_MODIFIED + 1))
        TOTAL_CHANGES=$((TOTAL_CHANGES + changes_made))
        log "SUCCESS" "✓ $filename: Processed and cleaned internal references"
    else
        rm -f "$temp_file" "$temp_file2"
        log "DEBUG" "No changes needed for $filename"
    fi
    
    return 0
}

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

parse_args() {
    while [ $# -gt 0 ]; do
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

main() {
    parse_args "$@"
    
    if [ ! -d "$DIRECTORY" ]; then
        log "ERROR" "Directory does not exist: $DIRECTORY"
        exit 1
    fi
    
    log "INFO" "Precise SVG Internal Links Removal Script"
    log "INFO" "Directory: $DIRECTORY"
    
    if [ "$DRY_RUN" = true ]; then
        log "INFO" "DRY RUN MODE - No files will be modified"
    fi
    
    if [ "$BACKUP" = true ] && [ "$DRY_RUN" = false ]; then
        log "INFO" "Backup mode enabled - .bak files will be created"
    fi
    
    echo "----------------------------------------"
    
    if process_directory; then
        echo "----------------------------------------"
        log "INFO" "Processing complete!"
        log "INFO" "Files processed: $FILES_PROCESSED"
        
        if [ "$DRY_RUN" = false ]; then
            log "INFO" "Files modified: $FILES_MODIFIED"
            log "INFO" "Total changes made: $TOTAL_CHANGES"
        fi
    else
        exit 1
    fi
}

main "$@"