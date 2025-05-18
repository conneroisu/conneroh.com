// Package main is the entry point for the update-fixed command
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/cache"
	"github.com/conneroisu/conneroh.com/internal/cwg"
	"github.com/conneroisu/conneroh.com/internal/llama"
	"github.com/conneroisu/conneroh.com/internal/logger"
	"github.com/conneroisu/conneroh.com/internal/markdown"
	"github.com/conneroisu/conneroh.com/internal/tigris"
	"github.com/rotisserie/eris"
	"github.com/spf13/afero"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	_ "modernc.org/sqlite"
)

const (
	numWorkers    = 20
	taskBufferInt = 1000
)

var (
	workers    = flag.Int("workers", numWorkers, "number of parallel workers")
	taskBuffer = flag.Int("buffer", taskBufferInt, "size of task buffer")
)

// Models with custom, decoupled types to avoid circular references
type (
	// BaseModel is the base for all entity models
	BaseModel struct {
		bun.BaseModel 
		ID          int64     `bun:"id,pk,autoincrement"`
		Title       string    `bun:"title,notnull"`
		Slug        string    `bun:"slug,notnull,unique"`
		Description string    `bun:"description"`
		Content     string    `bun:"content"`
		BannerPath  string    `bun:"banner_path"`
		CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp"`
		X           float64   `bun:"x"`
		Y           float64   `bun:"y"`
		Z           float64   `bun:"z"`
	}

	// Post represents a blog post
	Post struct {
		BaseModel `bun:"table:posts,alias:p"`
	}

	// Tag represents a tag
	Tag struct {
		BaseModel `bun:"table:tags,alias:t"`
		Icon      string `bun:"icon"`
	}

	// Project represents a project
	Project struct {
		BaseModel `bun:"table:projects,alias:pr"`
	}

	// Cache represents a cached file
	Cache struct {
		bun.BaseModel `bun:"table:caches,alias:c"`
		ID            int64   `bun:"id,pk,autoincrement"`
		Path          string  `bun:"path,notnull,unique"`
		Hash          string  `bun:"hashed,notnull"`
		X             float64 `bun:"x"`
		Y             float64 `bun:"y"`
		Z             float64 `bun:"z"`
	}

	// PostToTag represents a many-to-many relationship between posts and tags
	PostToTag struct {
		bun.BaseModel `bun:"table:post_to_tags,alias:pt"`
		PostID        int64 `bun:"post_id,pk"`
		TagID         int64 `bun:"tag_id,pk"`
	}

	// PostToPost represents a many-to-many relationship between posts
	PostToPost struct {
		bun.BaseModel `bun:"table:post_to_posts,alias:pp"`
		SourcePostID  int64 `bun:"source_post_id,pk"`
		TargetPostID  int64 `bun:"target_post_id,pk"`
	}

	// PostToProject represents a many-to-many relationship between posts and projects
	PostToProject struct {
		bun.BaseModel `bun:"table:post_to_projects,alias:ppr"`
		PostID        int64 `bun:"post_id,pk"`
		ProjectID     int64 `bun:"project_id,pk"`
	}

	// ProjectToTag represents a many-to-many relationship between projects and tags
	ProjectToTag struct {
		bun.BaseModel `bun:"table:project_to_tags,alias:prt"`
		ProjectID     int64 `bun:"project_id,pk"`
		TagID         int64 `bun:"tag_id,pk"`
	}

	// ProjectToProject represents a many-to-many relationship between projects
	ProjectToProject struct {
		bun.BaseModel    `bun:"table:project_to_projects,alias:prpr"`
		SourceProjectID  int64 `bun:"source_project_id,pk"`
		TargetProjectID  int64 `bun:"target_project_id,pk"`
	}

	// TagToTag represents a many-to-many relationship between tags
	TagToTag struct {
		bun.BaseModel `bun:"table:tag_to_tags,alias:tt"`
		SourceTagID   int64 `bun:"source_tag_id,pk"`
		TargetTagID   int64 `bun:"target_tag_id,pk"`
	}
)

func main() {
	flag.Parse()
	slog.SetDefault(logger.DefaultProdLogger)

	// Create context that will be canceled on interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP)
	defer stop()

	if err := Run(ctx, os.Getenv, *workers, *taskBuffer); err != nil {
		formattedStr := eris.ToString(err, true)
		fmt.Println(formattedStr)
		panic(err)
	}
}

// Run executes the main application logic.
func Run(ctx context.Context, getenv func(string) string, numWorkers, bufferSize int) error {
	// Initialize error collection
	var errs []error
	errMutex := &sync.Mutex{}

	// Create database connection
	sqldb, err := sql.Open("sqlite", assets.DBName())
	if err != nil {
		return eris.Wrap(err, "failed to open database")
	}
	defer sqldb.Close()

	// Initialize BUN DB
	db := bun.NewDB(sqldb, sqlitedialect.New())

	if os.Getenv("DEBUG") == "true" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	// Enable foreign keys in SQLite
	_, err = db.ExecContext(ctx, "PRAGMA foreign_keys = ON;")
	if err != nil {
		return eris.Wrap(err, "failed to enable foreign keys")
	}
	
	// We need to create minimal database schema to ensure it works
	slog.Info("Creating database schema if it doesn't exist")
	
	// Register entity models first, then relationship models using our custom types
	// This avoids circular reference issues
	db.RegisterModel(
		(*Post)(nil),
		(*Tag)(nil),
		(*Project)(nil),
		(*Cache)(nil),
	)
	db.RegisterModel(
		(*PostToTag)(nil),
		(*PostToPost)(nil),
		(*PostToProject)(nil),
		(*ProjectToTag)(nil),
		(*ProjectToProject)(nil),
		(*TagToTag)(nil),
	)
	
	// Create entity tables if they don't exist yet
	for _, model := range []interface{}{
		(*Post)(nil),
		(*Tag)(nil),
		(*Project)(nil),
		(*Cache)(nil),
	} {
		_, err := db.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			return eris.Wrapf(err, "failed to create table for %T", model)
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
			return eris.Wrapf(err, "failed to create relationship table for %T", rel.model)
		}
	}
	
	slog.Info("Database schema created successfully")

	// Initialize filesystem
	fs := afero.NewBasePathFs(afero.NewOsFs(), assets.VaultLoc)

	// Initialize dependencies
	ol, err := llama.NewOllamaClient(getenv)
	if err != nil {
		return eris.Wrap(err, "failed to create Ollama client")
	}

	ti, err := tigris.New(getenv)
	if err != nil {
		return eris.Wrap(err, "failed to create Tigris client")
	}

	md := markdown.NewMD(fs)

	// Create processor with specified number of workers
	processor := cache.NewProcessor(fs, db, ti, ol, md, bufferSize)

	// Create a context with cancellation
	processingCtx, cancelProcessing := context.WithCancel(ctx)
	defer cancelProcessing()

	// Start the processor workers
	processor.Start(processingCtx, numWorkers)

	// Handle errors from the processor
	errWg := &cwg.DebugWaitGroup{}
	errWg.Add(1)
	go func() {
		defer errWg.Done()
		for {
			select {
			case err, ok := <-processor.Errors():
				if !ok {
					return // Channel closed
				}
				errMutex.Lock()
				errs = append(errs, err)
				fmt.Print("processing error", "err", err)
				errMutex.Unlock()
			case <-processingCtx.Done():
				return
			}
		}
	}()

	// Scan the filesystem and submit tasks
	slog.Info("scanning filesystem", "path", assets.VaultLoc)
	if err := processor.ScanFS(); err != nil {
		errMutex.Lock()
		errs = append(errs, eris.Wrap(err, "failed to scan filesystem"))
		errMutex.Unlock()
		cancelProcessing() // Cancel processing on scan error

		return combineErrors(errs)
	}

	// Set up a forced exit after an absolute timeout (3 minutes)
	// This is an emergency fallback in case all other mechanisms fail
	forceExitCh := make(chan struct{})
	go func() {
		// Wait for 3 minutes then force exit
		time.Sleep(3 * time.Minute)
		slog.Error("EMERGENCY TIMEOUT - forcing exit after 3 minutes")
		close(forceExitCh)
	}()

	// Wait for completion with a timeout, or handle Ctrl+C
	completionTimeout := 2 * time.Minute

	// Create separate channels for completion waiting
	completionDone := make(chan bool)

	// Wait for completion in a separate goroutine
	go func() {
		completed := processor.WaitForCompletion(completionTimeout)
		completionDone <- completed
	}()

	// Wait for either completion, context cancellation, or force exit
	select {
	case completed := <-completionDone:
		if completed {
			slog.Info("tasks completed successfully")
		} else {
			slog.Warn("completion timeout reached - forcing shutdown")
			cancelProcessing()
		}
	case <-ctx.Done():
		slog.Info("received termination signal")
		cancelProcessing()
	case <-forceExitCh:
		slog.Error("emergency exit triggered - shutdown forced")
		cancelProcessing()
	}

	// Wait for error handling to finish
	errWg.Wait()

	// Clean up processor
	processor.Close()

	// Check if any errors occurred
	errMutex.Lock()
	defer errMutex.Unlock()

	if len(errs) > 0 {
		return combineErrors(errs)
	}

	slog.Info("processing completed successfully")

	return nil
}

// combineErrors combines multiple errors into one.
func combineErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var combinedErr error
	for _, err := range errs {
		if combinedErr == nil {
			combinedErr = err
		} else {
			combinedErr = eris.Wrap(combinedErr, err.Error())
		}
	}

	return eris.Wrap(combinedErr, "processing completed with errors")
}