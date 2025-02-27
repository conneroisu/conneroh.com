// Package master contains the master schema for the database.
package master

import (
	_ "embed"
)

//go:generate sqlcquash combine
//go:generate sqlc generate
//go:generate rm -f db.go
//go:generate gomarkdoc -o README.md -e .
//go:generate rm -f static.db
//go:generate sh -c "cat combined/schema.sql | sqlite3 -batch static.db"

// Schema is the schema for the database.
//
//go:embed combined/schema.sql
var Schema string

// Seed is the seed for the database.
//
//go:embed combined/seeds.sql
var Seed string
