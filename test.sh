#!/bin/bash

# Usage: ./interpolate.sh input_file start_marker end_marker replacement_text

input_file="$1"
start_marker="$2"
end_marker="$3"
replacement="$4"

# Create a temporary file
temp_file=$(mktemp)

# Process the file in three parts:
# 1. Copy everything before the start marker
# 2. Add the replacement text
# 3. Copy everything after the end marker
awk -v start="$start_marker" -v end="$end_marker" -v repl="$replacement" '
  # Flag to track if we are in the section to replace
  BEGIN { in_section = 0 }
  
  # If we find the start marker and are not yet in the section
  $0 ~ start && !in_section {
    print $0    # Print the start marker line
    print repl  # Print the replacement text
    in_section = 1
    next
  }
  
  # Print lines outside the section
  !in_section { print $0 }
' "$input_file" > "$temp_file"

# Replace the original file with the modified content
mv "$temp_file" "$input_file"
