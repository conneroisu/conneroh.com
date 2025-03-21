package conneroh

import (
	"context"
	"errors"
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
		slog.Info(
			"request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
		)
		mux.ServeHTTP(w, r)
	})
	var handler http.Handler = slogLogHandler
	return handler
}

// NewDb creates a new database connection.
func NewDb(_ func(string) string) (*data.Database[master.Queries], error) {
	return data.NewDb(master.New, "master.db")
}

// Run is the entry point for the application.
func Run(
	ctx context.Context,
	getenv func(string) string,
) error {
	db, err := NewDb(getenv)
	if err != nil {
		return err
	}

	// Create a separate context for signal handling
	innerCtx, stop := signal.NotifyContext(
		context.Background(), // Use a fresh context instead of the parent ctx
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	var (
		httpServer *http.Server
		wg         sync.WaitGroup
	)

	handler := NewServer(ctx, db) // Use original context for server setup

	// Configure server with timeouts
	httpServer = &http.Server{
		Addr: net.JoinHostPort(
			defaultHost,
			defaultPort,
		),
		Handler:           handler,
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
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("server error: %w", err)
			slog.Error("server error", slog.String("error", err.Error()))
		}
	}()

	// Wait for either server error or shutdown signal
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-innerCtx.Done():
		// Signal received, initiate graceful shutdown
		slog.Info("shutdown signal received, shutting down server...")

		// Create shutdown context with timeout
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(), // Use a fresh context for shutdown
			shutdownTimeout,
		)
		defer cancel()

		// Attempt graceful shutdown
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			slog.Error("error during server shutdown",
				slog.String("error", err.Error()),
				slog.Duration("timeout", shutdownTimeout),
			)
		}

		// Wait for all goroutines to finish
		slog.Info("waiting for server shutdown to complete")
		wg.Wait()
		slog.Info("server shutdown completed")
		return nil
	case <-ctx.Done():
		// Parent context cancelled
		slog.Info("parent context cancelled, shutting down...")

		// Create shutdown context with timeout
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(), // Use a fresh context for shutdown
			shutdownTimeout,
		)
		defer cancel()

		// Attempt graceful shutdown
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			slog.Error("error during server shutdown",
				slog.String("error", err.Error()),
				slog.Duration("timeout", shutdownTimeout),
			)
		}

		// Wait for all goroutines to finish
		wg.Wait()
		return nil
	}
}
