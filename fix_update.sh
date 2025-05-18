#!/bin/bash

# Create a version of the update command that fixes the registration issue
mkdir -p tmp

cat > tmp/main.go << 'EOF'
// Package main is a temporary fix for the update command
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/logger"
	"github.com/rotisserie/eris"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	_ "modernc.org/sqlite"
)

func main() {
	slog.SetDefault(logger.DefaultProdLogger)

	// Create context that will be canceled on interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP)
	defer stop()

	// Open database connection
	sqldb, err := sql.Open("sqlite", assets.DBName())
	if err != nil {
		panic(eris.Wrap(err, "failed to open database"))
	}
	defer sqldb.Close()

	// Initialize BUN DB
	db := bun.NewDB(sqldb, sqlitedialect.New())
	defer db.Close()

	if os.Getenv("DEBUG") == "true" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	// First, register the models
	registerModels(db)

	// Enable foreign keys
	_, err = db.ExecContext(ctx, "PRAGMA foreign_keys = ON;")
	if err != nil {
		panic(eris.Wrap(err, "failed to enable foreign keys"))
	}

	// Create a new clean database
	err = createTables(ctx, db)
	if err != nil {
		panic(eris.Wrap(err, "failed to create tables"))
	}

	fmt.Println("Database has been successfully recreated!")
}

func registerModels(db *bun.DB) {
	// Register entity models
	db.RegisterModel(
		(*assets.Post)(nil),
		(*assets.Tag)(nil),
		(*assets.Project)(nil),
		(*assets.Cache)(nil),
	)
	
	// Register relationship models
	db.RegisterModel(
		(*assets.PostToTag)(nil),
		(*assets.PostToPost)(nil),
		(*assets.PostToProject)(nil),
		(*assets.ProjectToTag)(nil),
		(*assets.ProjectToProject)(nil),
		(*assets.TagToTag)(nil),
	)
}

func createTables(ctx context.Context, db *bun.DB) error {
	// Create main entity tables
	models := []any{
		(*assets.Post)(nil),
		(*assets.Tag)(nil),
		(*assets.Project)(nil),
		(*assets.Cache)(nil),
	}

	for _, model := range models {
		_, err := db.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			return eris.Wrap(err, fmt.Sprintf("failed to create table for %T", model))
		}
	}

	// Create relationship tables
	relationships := []struct {
		model any
		fks   []string
	}{
		{
			model: (*assets.PostToTag)(nil),
			fks: []string{
				"(post_id) REFERENCES posts(id) ON DELETE CASCADE",
				"(tag_id) REFERENCES tags(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*assets.PostToPost)(nil),
			fks: []string{
				"(source_post_id) REFERENCES posts(id) ON DELETE CASCADE",
				"(target_post_id) REFERENCES posts(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*assets.PostToProject)(nil),
			fks: []string{
				"(post_id) REFERENCES posts(id) ON DELETE CASCADE",
				"(project_id) REFERENCES projects(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*assets.ProjectToTag)(nil),
			fks: []string{
				"(project_id) REFERENCES projects(id) ON DELETE CASCADE",
				"(tag_id) REFERENCES tags(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*assets.ProjectToProject)(nil),
			fks: []string{
				"(source_project_id) REFERENCES projects(id) ON DELETE CASCADE",
				"(target_project_id) REFERENCES projects(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*assets.TagToTag)(nil),
			fks: []string{
				"(source_tag_id) REFERENCES tags(id) ON DELETE CASCADE",
				"(target_tag_id) REFERENCES tags(id) ON DELETE CASCADE",
			},
		},
	}

	for _, rel := range relationships {
		query := db.NewCreateTable().
			Model(rel.model).
			IfNotExists()
		
		for _, fk := range rel.fks {
			query = query.ForeignKey(fk)
		}
		
		_, err := query.Exec(ctx)
		if err != nil {
			return eris.Wrap(err, fmt.Sprintf("failed to create table for %T", rel.model))
		}
	}

	return nil
}
EOF

# Remove existing database files
rm -f master.db master.db-shm master.db-wal

# Run the fix
go run tmp/main.go

# Clean up
rm -rf tmp

echo "Fix applied successfully!"