#!/bin/bash

# Script to generate and update the database

# Exit on any error
set -e

echo "=== Database Generation Process ==="

echo "Creating database with schema and content..."
go run ./cmd/simple-import/main.go

echo "=== Database generation complete! ==="