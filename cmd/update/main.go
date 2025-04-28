// Package main is the entry point for the application
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
	"github.com/conneroisu/conneroh.com/internal/llama"
	"github.com/conneroisu/conneroh.com/internal/logger"
	"github.com/conneroisu/conneroh.com/internal/markdown"
	"github.com/conneroisu/conneroh.com/internal/tigris"
	"github.com/rotisserie/eris"
	"github.com/spf13/afero"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	_ "modernc.org/sqlite"
)

var (
	workers    = flag.Int("workers", 20, "number of parallel workers")
	taskBuffer = flag.Int("buffer", 1000, "size of task buffer")
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
		os.Exit(1)
	}
}

// Run executes the main application logic
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

	// db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	// Initialize database tables
	err = assets.InitDB(ctx, db)
	if err != nil {
		return eris.Wrap(err, "failed to initialize database")
	}

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
	// errWg := cwg.DebugWaitGroup{}
	errWg := &sync.WaitGroup{}
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
				slog.Error("processing error", "err", err)
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

// combineErrors combines multiple errors into one
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
