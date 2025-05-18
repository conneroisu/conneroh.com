#!/bin/bash

# Exit on any error
set -e

echo "=== Running Application with Doppler ==="

if command -v doppler >/dev/null 2>&1; then
  echo "Using Doppler for environment variables..."
  doppler run -- go run ./main.go
else
  echo "Doppler not found, running without environment variables..."
  go run ./main.go
fi