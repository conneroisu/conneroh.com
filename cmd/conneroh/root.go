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

	"github.com/conneroisu/conneroh.com/internal/data/css"
	"github.com/conneroisu/twerge"
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
) http.Handler {
	twerge.ClassMapStr = css.ClassMapStr
	twerge.GenClassMergeStr = css.ClassMapStr
	mux := http.NewServeMux()
	err := AddRoutes(ctx, mux)
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

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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
	start := time.Now()

	// Create a separate context for signal handling
	innerCtx, cancel := signal.NotifyContext(
		context.Background(), // Use a fresh context instead of the parent ctx
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	defer cancel()

	var (
		httpServer *http.Server
		wg         sync.WaitGroup
	)

	handler := NewServer(ctx) // Use original context for server setup

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
		slog.Info(
			"server starting",
			slog.String("address", httpServer.Addr),
			slog.String("setup-time", time.Since(start).String()),
		)
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("server error: %w", err)
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
