package conneroh

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/master"
)

const (
	defaultHost       = "0.0.0.0"
	defaultPort       = "8080"
	shutdownTimeout   = 10 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	readHeaderTimeout = 5 * time.Second
)

// NewServer creates a new web-ui server
func NewServer(
	ctx context.Context,
	db *data.Database[master.Queries],
) http.Handler {
	mux := http.NewServeMux()
	err := AddRoutes(
		ctx,
		mux,
		db,
	)
	if err != nil {
		slog.Error("error adding routes", slog.String("error", err.Error()))
		log.Fatal(err)
	}
	slogLogHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", slog.String("method", r.Method), slog.String("url", r.URL.String()))
		mux.ServeHTTP(w, r)
	})
	var handler http.Handler = slogLogHandler
	return handler
}

// Run is the entry point for the application.
func Run(
	ctx context.Context,
	_ func(string) string,
) error {

	innerCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	var (
		httpServer *http.Server
		wg         sync.WaitGroup
	)

	db, err := data.NewDb(master.New, &data.Config{
		Schema:   master.Schema,
		URI:      "file://dummy",
		Seed:     master.Seed,
		FileName: "test.db",
	})
	if err != nil {
		return err
	}

	// Configure server with timeouts
	httpServer = &http.Server{
		Addr:              net.JoinHostPort(defaultHost, defaultPort),
		Handler:           NewServer(innerCtx, db),
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	serverErrors := make(chan error, 1)
	// Start server
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("server starting", slog.String("address", httpServer.Addr))

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- fmt.Errorf("server error: %w", err)
			slog.Error("server error", slog.String("error", err.Error()))
		}
	}()

	// Shutdown handler
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Wait for context cancellation (signal or parent context)
		<-innerCtx.Done()

		slog.Info("shutting down server...")

		// Create shutdown context with timeout
		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		// Attempt graceful shutdown
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("error during server shutdown",
				slog.String("error", err.Error()),
				slog.Duration("timeout", shutdownTimeout),
			)
		}
	}()

	// Wait for server error or successful shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-innerCtx.Done():
		// Wait for graceful shutdown to complete
		wg.Wait()
		slog.Info("server shutdown completed")
		return nil
	case <-ctx.Done():
		slog.Info("server shutdown ctx cancelled")
		return nil
	}
}
