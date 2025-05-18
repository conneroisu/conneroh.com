// Package main provides a simplified update command for the database
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/logger"
	"github.com/rotisserie/eris"
	"github.com/spf13/afero"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	_ "modernc.org/sqlite"
)

// SimplePost represents a post in the database
type SimplePost struct {
	bun.BaseModel `bun:"table:posts,alias:p"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Title         string    `bun:"title,notnull"`
	Slug          string    `bun:"slug,notnull,unique"`
	Description   string    `bun:"description"`
	Content       string    `bun:"content"`
	BannerPath    string    `bun:"banner_path"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp"`
	X             float64   `bun:"x"`
	Y             float64   `bun:"y"`
	Z             float64   `bun:"z"`
}

// SimpleTag represents a tag in the database
type SimpleTag struct {
	bun.BaseModel `bun:"table:tags,alias:t"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Title         string    `bun:"title,notnull"`
	Slug          string    `bun:"slug,notnull,unique"`
	Description   string    `bun:"description"`
	Content       string    `bun:"content"`
	BannerPath    string    `bun:"banner_path"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp"`
	X             float64   `bun:"x"`
	Y             float64   `bun:"y"`
	Z             float64   `bun:"z"`
	Icon          string    `bun:"icon"` // Only for tags
}

// SimpleProject represents a project in the database
type SimpleProject struct {
	bun.BaseModel `bun:"table:projects,alias:pr"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Title         string    `bun:"title,notnull"`
	Slug          string    `bun:"slug,notnull,unique"`
	Description   string    `bun:"description"`
	Content       string    `bun:"content"`
	BannerPath    string    `bun:"banner_path"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp"`
	X             float64   `bun:"x"`
	Y             float64   `bun:"y"`
	Z             float64   `bun:"z"`
}

// SimpleCache represents a cache entry in the database
type SimpleCache struct {
	bun.BaseModel `bun:"table:caches,alias:c"`
	ID            int64   `bun:"id,pk,autoincrement"`
	Path          string  `bun:"path,notnull,unique"`
	Hash          string  `bun:"hashed,notnull"`
	X             float64 `bun:"x"`
	Y             float64 `bun:"y"`
	Z             float64 `bun:"z"`
}

// SimplePostToTag represents a many-to-many relationship between posts and tags
type SimplePostToTag struct {
	bun.BaseModel `bun:"table:post_to_tags,alias:pt"`
	PostID        int64 `bun:"post_id,pk"`
	TagID         int64 `bun:"tag_id,pk"`
}

// SimplePostToPost represents a many-to-many relationship between posts
type SimplePostToPost struct {
	bun.BaseModel   `bun:"table:post_to_posts,alias:pp"`
	SourcePostID    int64 `bun:"source_post_id,pk"`
	TargetPostID    int64 `bun:"target_post_id,pk"`
}

// SimplePostToProject represents a many-to-many relationship between posts and projects
type SimplePostToProject struct {
	bun.BaseModel `bun:"table:post_to_projects,alias:ppr"`
	PostID        int64 `bun:"post_id,pk"`
	ProjectID     int64 `bun:"project_id,pk"`
}

// SimpleProjectToTag represents a many-to-many relationship between projects and tags
type SimpleProjectToTag struct {
	bun.BaseModel `bun:"table:project_to_tags,alias:prt"`
	ProjectID     int64 `bun:"project_id,pk"`
	TagID         int64 `bun:"tag_id,pk"`
}

// SimpleProjectToProject represents a many-to-many relationship between projects
type SimpleProjectToProject struct {
	bun.BaseModel    `bun:"table:project_to_projects,alias:prpr"`
	SourceProjectID  int64 `bun:"source_project_id,pk"`
	TargetProjectID  int64 `bun:"target_project_id,pk"`
}

// SimpleTagToTag represents a many-to-many relationship between tags
type SimpleTagToTag struct {
	bun.BaseModel `bun:"table:tag_to_tags,alias:tt"`
	SourceTagID   int64 `bun:"source_tag_id,pk"`
	TargetTagID   int64 `bun:"target_tag_id,pk"`
}

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
	slog.Info("Starting database update")

	// Step 1: Remove existing database
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

	// Step 2: Connect to database
	sqldb, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return eris.Wrap(err, "failed to open database")
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
		return eris.Wrap(err, "failed to enable foreign keys")
	}

	// Step 3: Register models with Bun
	db.RegisterModel(
		(*SimplePost)(nil),
		(*SimpleTag)(nil),
		(*SimpleProject)(nil),
		(*SimpleCache)(nil),
		(*SimplePostToTag)(nil),
		(*SimplePostToPost)(nil),
		(*SimplePostToProject)(nil),
		(*SimpleProjectToTag)(nil),
		(*SimpleProjectToProject)(nil),
		(*SimpleTagToTag)(nil),
	)

	// Step 4: Create database tables
	slog.Info("creating database tables")

	// Create entity tables
	models := []interface{}{
		(*SimplePost)(nil),
		(*SimpleTag)(nil),
		(*SimpleProject)(nil),
		(*SimpleCache)(nil),
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
			model: (*SimplePostToTag)(nil),
			fks: []string{
				"(post_id) REFERENCES posts(id) ON DELETE CASCADE",
				"(tag_id) REFERENCES tags(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*SimplePostToPost)(nil),
			fks: []string{
				"(source_post_id) REFERENCES posts(id) ON DELETE CASCADE",
				"(target_post_id) REFERENCES posts(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*SimplePostToProject)(nil),
			fks: []string{
				"(post_id) REFERENCES posts(id) ON DELETE CASCADE",
				"(project_id) REFERENCES projects(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*SimpleProjectToTag)(nil),
			fks: []string{
				"(project_id) REFERENCES projects(id) ON DELETE CASCADE",
				"(tag_id) REFERENCES tags(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*SimpleProjectToProject)(nil),
			fks: []string{
				"(source_project_id) REFERENCES projects(id) ON DELETE CASCADE",
				"(target_project_id) REFERENCES projects(id) ON DELETE CASCADE",
			},
		},
		{
			model: (*SimpleTagToTag)(nil),
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

	// Step 5: Import content from markdown files
	slog.Info("Importing content from markdown files")

	// Create filesystem access
	fs := afero.NewBasePathFs(afero.NewOsFs(), assets.VaultLoc)

	// Process posts
	if err := processMarkdownDir(ctx, fs, assets.PostsLoc, func(doc *assets.Doc) error {
		post := &SimplePost{
			Title:       doc.Title,
			Slug:        doc.Slug,
			Description: doc.Description,
			Content:     doc.Content,
			BannerPath:  doc.BannerPath,
		}
		
		// Insert post
		_, err := db.NewInsert().Model(post).Exec(ctx)
		return err
	}); err != nil {
		return err
	}

	// Process projects
	if err := processMarkdownDir(ctx, fs, assets.ProjectsLoc, func(doc *assets.Doc) error {
		project := &SimpleProject{
			Title:       doc.Title,
			Slug:        doc.Slug,
			Description: doc.Description,
			Content:     doc.Content,
			BannerPath:  doc.BannerPath,
		}
		
		// Insert project
		_, err := db.NewInsert().Model(project).Exec(ctx)
		return err
	}); err != nil {
		return err
	}

	// Process tags
	if err := processMarkdownDir(ctx, fs, assets.TagsLoc, func(doc *assets.Doc) error {
		tag := &SimpleTag{
			Title:       doc.Title,
			Slug:        doc.Slug,
			Description: doc.Description,
			Content:     doc.Content,
			BannerPath:  doc.BannerPath,
			Icon:        doc.Icon,
		}
		
		// Insert tag
		_, err := db.NewInsert().Model(tag).Exec(ctx)
		return err
	}); err != nil {
		return err
	}

	slog.Info("Content import completed successfully")
	return nil
}

// processMarkdownDir processes all markdown files in a directory
func processMarkdownDir(ctx context.Context, fs afero.Fs, dirPath string, processor func(*assets.Doc) error) error {
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
		if !assets.IsAllowedDocumentType(file.Name()) {
			continue
		}
		
		filePath := dirPath + "/" + file.Name()
		slog.Info("Processing file", "path", filePath)
		
		// Read file content
		content, err := afero.ReadFile(fs, filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}
		
		// Create a doc to store the parsed data
		doc := &assets.Doc{
			Path: filePath,
			Hash: assets.Hash(content),
			// Add basic defaults
			Title: file.Name(),
			Slug:  assets.Slugify(filePath),
		}
		
		// Set default values
		if err := assets.Defaults(doc); err != nil {
			slog.Warn("Failed to set defaults", "path", filePath, "error", err)
			// Continue instead of failing
		}
		
		// Process the document
		if err := processor(doc); err != nil {
			return fmt.Errorf("failed to process %s: %w", filePath, err)
		}
		
		slog.Info("Successfully processed file", "path", filePath)
	}

	return nil
}