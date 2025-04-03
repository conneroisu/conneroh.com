#!/usr/bin/env bash
set -e
EXCLUDE=""
DIR_PATH=""
# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    -dir) DIR_PATH="$2" shift 2 ;;
    -exclude) EXCLUDE="$2" shift 2 ;;
    *) echo "Unknown option: $1" exit 1 ;;
  esac
done
if [[ -z "$DIR_PATH" ]]; then
  echo "Error: directory path is required. Use -dir flag." && exit 1
fi
if [[ ! -d "$DIR_PATH" ]]; then
  echo "Error: $DIR_PATH is not a directory" && exit 1
fi
if [[ -n "$EXCLUDE" ]]; then
  # Convert comma-separated exclude patterns to Git-compatible ignore pattern
  EXCLUDE_ARGS=()
  IFS=',' read -ra PATTERNS <<< "$EXCLUDE"
  for pattern in "${PATTERNS[@]}"; do
    EXCLUDE_ARGS+=(":(exclude)$pattern")
  done
fi
# --name-only shows only filenames
# --diff-filter=ACMRT filters for Added, Copied, Modified, Renamed, Type-changed files
DIFF_OUTPUT=$(git diff --name-only --diff-filter=ACMRT -- "$DIR_PATH" "${EXCLUDE_ARGS[@]}" 2>/dev/null)
# Also check for untracked files
UNTRACKED_OUTPUT=$(git ls-files --others --exclude-standard -- "$DIR_PATH" "${EXCLUDE_ARGS[@]}" 2>/dev/null)
# Combine outputs
STATUS_OUTPUT="${DIFF_OUTPUT}${UNTRACKED_OUTPUT:+\n'$UNTRACKED_OUTPUT}"
# Check if there are any changes
if [[ -z "$STATUS_OUTPUT" ]]; then
  echo "No changes detected in $DIR_PATH"
  exit 0
else
  echo "Changes detected in $DIR_PATH"
  if [[ -n "$EXCLUDE" ]]; then
    # Add files with exclude patterns
    git add "$DIR_PATH" "${EXCLUDE_ARGS[@]}"
  else
    # Add all files in the directory
    git add "$DIR_PATH"
  fi
  exit 1
fi
