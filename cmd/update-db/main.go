// Package main is a simplified version of the update command
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/conneroisu/conneroh.com/internal/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	_ "modernc.org/sqlite"
)

// Define simplified models to avoid dependency issues
type (
	// Base model for all entities
	BaseEntity struct {
		bun.BaseModel
		ID          int64 `bun:"id,pk,autoincrement"`
		Title       string
		Slug        string `bun:"slug,unique"`
		Description string
		Content     string
		BannerPath  string
	}

	// Post model
	Post struct {
		BaseEntity `bun:"table:posts"`
	}

	// Tag model
	Tag struct {
		BaseEntity `bun:"table:tags"`
		Icon       string
	}

	// Project model  
	Project struct {
		BaseEntity `bun:"table:projects"`
	}

	// Cache model
	Cache struct {
		bun.BaseModel `bun:"table:caches"`
		ID            int64  `bun:"id,pk,autoincrement"`
		Path          string `bun:"path,unique"`
		Hash          string `bun:"hashed,unique"`
	}

	// Relationship models
	PostToTag struct {
		bun.BaseModel `bun:"table:post_to_tags"`
		PostID        int64 `bun:"post_id,pk"`
		TagID         int64 `bun:"tag_id,pk"`
	}

	PostToPost struct {
		bun.BaseModel `bun:"table:post_to_posts"`
		SourcePostID  int64 `bun:"source_post_id,pk"`
		TargetPostID  int64 `bun:"target_post_id,pk"`
	}

	PostToProject struct {
		bun.BaseModel `bun:"table:post_to_projects"`
		PostID        int64 `bun:"post_id,pk"`
		ProjectID     int64 `bun:"project_id,pk"`
	}

	ProjectToTag struct {
		bun.BaseModel `bun:"table:project_to_tags"`
		ProjectID     int64 `bun:"project_id,pk"`
		TagID         int64 `bun:"tag_id,pk"`
	}

	ProjectToProject struct {
		bun.BaseModel `bun:"table:project_to_projects"`
		SourceProjectID int64 `bun:"source_project_id,pk"`
		TargetProjectID int64 `bun:"target_project_id,pk"`
	}

	TagToTag struct {
		bun.BaseModel `bun:"table:tag_to_tags"`
		SourceTagID   int64 `bun:"source_tag_id,pk"`
		TargetTagID   int64 `bun:"target_tag_id,pk"`
	}
)

const dbName = "file:master.db?_pragma=busy_timeout=5000&_pragma=journal_mode=WAL&_pragma=mmap_size=30000000000&_pragma=page_size=32768"

func main() {
	slog.SetDefault(logger.DefaultProdLogger)

	// Create context that will be canceled on interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Error("Error running update", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Connect to database
	sqldb, err := sql.Open("sqlite", dbName)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer sqldb.Close()

	// Create bun DB
	db := bun.NewDB(sqldb, sqlitedialect.New())
	defer db.Close()

	if os.Getenv("DEBUG") == "true" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	// Enable foreign keys
	_, err = db.ExecContext(ctx, "PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Create entity tables
	models := []interface{}{
		(*Post)(nil),
		(*Tag)(nil),
		(*Project)(nil),
		(*Cache)(nil),
	}

	for _, model := range models {
		_, err = db.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create table for %T: %w", model, err)
		}
	}

	// Create relationship tables with foreign keys
	relationships := []struct {
		model interface{}
		fks   []string
	}{
		{
			model: (*PostToTag)(nil),
			fks: []string{
				"(post_id) REFERENCES posts(id) ON DELETE CASCADE",
				"(tag_id) REFERENCES tags(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*PostToPost)(nil),
			fks: []string{
				"(source_post_id) REFERENCES posts(id) ON DELETE CASCADE",
				"(target_post_id) REFERENCES posts(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*PostToProject)(nil),
			fks: []string{
				"(post_id) REFERENCES posts(id) ON DELETE CASCADE",
				"(project_id) REFERENCES projects(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*ProjectToTag)(nil),
			fks: []string{
				"(project_id) REFERENCES projects(id) ON DELETE CASCADE",
				"(tag_id) REFERENCES tags(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*ProjectToProject)(nil),
			fks: []string{
				"(source_project_id) REFERENCES projects(id) ON DELETE CASCADE",
				"(target_project_id) REFERENCES projects(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*TagToTag)(nil),
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
		
		_, err = query.Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create relationship table for %T: %w", rel.model, err)
		}
	}

	slog.Info("Database tables created successfully")
	return nil
}