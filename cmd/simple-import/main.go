// Package main provides a simple importer for markdown content
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/logger"
	"github.com/yuin/goldmark"
	"github.com/spf13/afero"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	_ "modernc.org/sqlite"
)

func main() {
	slog.SetDefault(logger.DefaultProdLogger)
	ctx := context.Background()

	if err := run(ctx); err != nil {
		slog.Error("Error running import", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	slog.Info("Starting simple database import")

	// Remove existing database
	dbPath := assets.DBName()
	slog.Info("removing existing database", "path", dbPath)
	
	// Extract actual file path from SQLite connection string
	dbFilePath := dbPath
	if len(dbPath) > 5 && dbPath[0:5] == "file:" {
		// Extract file path from "file:path?options"
		for i := 5; i < len(dbPath); i++ {
			if dbPath[i] == '?' {
				dbFilePath = dbPath[5:i]
				break
			}
		}
		if dbFilePath == dbPath[5:] {
			dbFilePath = dbPath[5:]
		}
	}
	
	// Remove database files
	os.Remove(dbFilePath)
	os.Remove(dbFilePath + "-shm")
	os.Remove(dbFilePath + "-wal")

	// Connect to database
	sqldb, err := sql.Open("sqlite", assets.DBName())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer sqldb.Close()

	// Initialize Bun DB
	db := bun.NewDB(sqldb, sqlitedialect.New())
	defer db.Close()

	// Enable debug mode if needed
	if os.Getenv("DEBUG") == "true" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	// Enable foreign keys
	_, err = db.ExecContext(ctx, "PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Create tables using SQL directly to avoid model issues
	slog.Info("Creating database tables")
	
	sqlStatements := []string{
		// Create main entity tables
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			description TEXT,
			content TEXT,
			banner_path TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			x REAL,
			y REAL,
			z REAL
		)`,
		`CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			description TEXT,
			content TEXT,
			banner_path TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			icon TEXT,
			x REAL,
			y REAL,
			z REAL
		)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			description TEXT,
			content TEXT,
			banner_path TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			x REAL,
			y REAL,
			z REAL
		)`,
		`CREATE TABLE IF NOT EXISTS caches (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL UNIQUE,
			hashed TEXT NOT NULL,
			x REAL,
			y REAL,
			z REAL
		)`,
		// Create relationship tables
		`CREATE TABLE IF NOT EXISTS post_to_tags (
			post_id INTEGER,
			tag_id INTEGER,
			PRIMARY KEY (post_id, tag_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS post_to_posts (
			source_post_id INTEGER,
			target_post_id INTEGER,
			PRIMARY KEY (source_post_id, target_post_id),
			FOREIGN KEY (source_post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (target_post_id) REFERENCES posts(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS post_to_projects (
			post_id INTEGER,
			project_id INTEGER,
			PRIMARY KEY (post_id, project_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS project_to_tags (
			project_id INTEGER,
			tag_id INTEGER,
			PRIMARY KEY (project_id, tag_id),
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS project_to_projects (
			source_project_id INTEGER,
			target_project_id INTEGER,
			PRIMARY KEY (source_project_id, target_project_id),
			FOREIGN KEY (source_project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY (target_project_id) REFERENCES projects(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS tag_to_tags (
			source_tag_id INTEGER,
			target_tag_id INTEGER,
			PRIMARY KEY (source_tag_id, target_tag_id),
			FOREIGN KEY (source_tag_id) REFERENCES tags(id) ON DELETE CASCADE,
			FOREIGN KEY (target_tag_id) REFERENCES tags(id) ON DELETE CASCADE
		)`,
	}

	for _, stmt := range sqlStatements {
		_, err := db.ExecContext(ctx, stmt)
		if err != nil {
			return fmt.Errorf("failed to execute SQL: %s: %w", stmt, err)
		}
	}

	slog.Info("Database tables created successfully")

	// Process markdown files
	fs := afero.NewBasePathFs(afero.NewOsFs(), assets.VaultLoc)
	md := goldmark.New()

	// Process posts
	if err := processMarkdownDir(ctx, db, fs, md, assets.PostsLoc, "post"); err != nil {
		return fmt.Errorf("failed to process posts: %w", err)
	}

	// Process projects
	if err := processMarkdownDir(ctx, db, fs, md, assets.ProjectsLoc, "project"); err != nil {
		return fmt.Errorf("failed to process projects: %w", err)
	}

	// Process tags
	if err := processMarkdownDir(ctx, db, fs, md, assets.TagsLoc, "tag"); err != nil {
		return fmt.Errorf("failed to process tags: %w", err)
	}

	// Process assets by just creating cache entries for main assets directory
	assetsDir := "assets"
	if err := processAssets(ctx, db, fs, assetsDir); err != nil {
		return fmt.Errorf("failed to process assets: %w", err)
	}

	slog.Info("Database import completed successfully")
	return nil
}

// processMarkdownDir processes all markdown files in a directory
func processMarkdownDir(ctx context.Context, db *bun.DB, fs afero.Fs, md goldmark.Markdown, dirPath, entityType string) error {
	// Read the directory
	files, err := afero.ReadDir(fs, dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	// Process each markdown file
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		// Only process .md files
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		
		filePath := dirPath + "/" + file.Name()
		slog.Info("Processing file", "path", filePath)
		
		// Read file content
		content, err := afero.ReadFile(fs, filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}
		
		// Create a simplified document
		fileName := filepath.Base(filePath)
		fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		slug := assets.Slugify(filePath)
		
		// Use a SQL statement to insert the entity directly
		var sqlStmt string
		var args []interface{}
		
		switch entityType {
		case "post":
			sqlStmt = `INSERT INTO posts (title, slug, description, content, banner_path, created_at)
				VALUES (?, ?, ?, ?, ?, ?)
				ON CONFLICT (slug) DO UPDATE SET
				title = excluded.title,
				description = excluded.description,
				content = excluded.content,
				banner_path = excluded.banner_path,
				created_at = excluded.created_at`
			args = []interface{}{
				fileNameWithoutExt,
				slug,
				"Description for " + fileNameWithoutExt,
				string(content),
				"", // banner_path
				time.Now(),
			}
		case "project":
			sqlStmt = `INSERT INTO projects (title, slug, description, content, banner_path, created_at)
				VALUES (?, ?, ?, ?, ?, ?)
				ON CONFLICT (slug) DO UPDATE SET
				title = excluded.title,
				description = excluded.description,
				content = excluded.content,
				banner_path = excluded.banner_path,
				created_at = excluded.created_at`
			args = []interface{}{
				fileNameWithoutExt,
				slug,
				"Description for " + fileNameWithoutExt,
				string(content),
				"", // banner_path
				time.Now(),
			}
		case "tag":
			sqlStmt = `INSERT INTO tags (title, slug, description, content, icon, created_at)
				VALUES (?, ?, ?, ?, ?, ?)
				ON CONFLICT (slug) DO UPDATE SET
				title = excluded.title,
				description = excluded.description,
				content = excluded.content,
				icon = excluded.icon,
				created_at = excluded.created_at`
			args = []interface{}{
				fileNameWithoutExt,
				slug,
				"Description for " + fileNameWithoutExt,
				string(content),
				"tag", // icon
				time.Now(),
			}
		default:
			return fmt.Errorf("invalid entity type: %s", entityType)
		}
		
		// Execute the SQL statement
		_, err = db.ExecContext(ctx, sqlStmt, args...)
		if err != nil {
			return fmt.Errorf("failed to insert %s %s: %w", entityType, slug, err)
		}
		
		slog.Info("Successfully processed file", "path", filePath, "slug", slug)
	}

	return nil
}

// processAssets adds cache entries for assets
func processAssets(ctx context.Context, db *bun.DB, fs afero.Fs, assetDir string) error {
	slog.Info("Processing assets", "dir", assetDir)
	err := afero.Walk(fs, assetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		// Only process allowed asset types
		if !assets.IsAllowedMediaType(path) {
			return nil
		}
		
		slog.Info("Processing asset", "path", path)
		
		// Hash the file contents
		content, err := afero.ReadFile(fs, path)
		if err != nil {
			return fmt.Errorf("failed to read asset %s: %w", path, err)
		}
		
		hash := assets.Hash(content)
		
		// Add to cache
		sqlStmt := `INSERT INTO caches (path, hashed)
			VALUES (?, ?)
			ON CONFLICT (path) DO UPDATE SET
			hashed = excluded.hashed`
		
		_, err = db.ExecContext(ctx, sqlStmt, path, hash)
		if err != nil {
			return fmt.Errorf("failed to cache asset %s: %w", path, err)
		}
		
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to process assets: %w", err)
	}
	
	slog.Info("Assets processing completed")
	return nil
}