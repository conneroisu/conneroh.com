#!/bin/bash

# Fix variable naming inconsistencies
if [[ -f internal/cache/procs.go ]]; then
  sed -i '' 's/pendingDbTasks/pendingDBTasks/g' internal/cache/procs.go
  sed -i '' 's/p\.stats\.DbTasksSubmitted/p\.stats\.DBTasksSubmitted/g' internal/cache/procs.go
  sed -i '' 's/p\.stats\.DbTasksCompleted/p\.stats\.DBTasksCompleted/g' internal/cache/procs.go
fi

# Reset and rebuild the database
echo "Removing existing database..."
rm -f master.db*

# Run the database update command
echo "Rebuilding database..."
go run ./cmd/update-db/main.go

echo "Database has been reset and rebuilt successfully!"