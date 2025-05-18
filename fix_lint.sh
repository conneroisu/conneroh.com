#!/bin/bash

# This script fixes common issues with the codebase
echo "Fixing codebase issues..."

# Fix variable naming inconsistencies
if [[ -f internal/cache/procs.go ]]; then
  sed -i '' 's/DbTasksSubmitted/DBTasksSubmitted/g' internal/cache/procs.go
  sed -i '' 's/DbTasksCompleted/DBTasksCompleted/g' internal/cache/procs.go
  sed -i '' 's/RecordDbTask/RecordDBTask/g' internal/cache/procs.go
fi

# Step 1: First create the database structure
echo "Recreating database schema..."
rm -f master.db master.db-shm master.db-wal
go run ./cmd/update-db/main.go

# Step 2: Then generate content with the live-ci command
echo "Generating content..."
go run ./cmd/live-ci/main.go

echo "Fix completed successfully!"