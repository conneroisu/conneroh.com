#!/bin/bash

# Exit on any error
set -e

echo "=== Running Update with Fixed Command ==="

if command -v doppler >/dev/null 2>&1; then
  echo "Using Doppler to run the update command..."
  doppler run -- go run ./cmd/update-fixed
else
  echo "Doppler not found, running without environment variables..."
  go run ./cmd/update-fixed
fi

echo "=== Update completed successfully ==="