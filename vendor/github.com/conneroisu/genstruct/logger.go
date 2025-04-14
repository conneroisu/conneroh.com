package genstruct

import (
	"flag"
	"io"
	"log/slog"
	"os"
)

var (
	// Default global logger
	defaultLogger *slog.Logger

	// Command line flags
	verbosity  string
	logFormat  string
	logOutput  string
	logHandler slog.Handler
)

// Verbosity levels
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

// InitLogger initializes the logger with command line flags
func InitLogger() *slog.Logger {
	// Only parse flags once
	if !flag.Parsed() {
		// Define flags
		flag.StringVar(&verbosity, "v", LevelInfo, "Log verbosity (debug, info, warn, error)")
		flag.StringVar(&logFormat, "log-format", "text", "Log format (text, json)")
		flag.StringVar(&logOutput, "log-output", "stderr", "Log output (stderr, stdout, or a file path)")
		flag.Parse()
	}

	// Determine log level
	var level slog.Level
	switch verbosity {
	case LevelDebug:
		level = slog.LevelDebug
	case LevelInfo:
		level = slog.LevelInfo
	case LevelWarn:
		level = slog.LevelWarn
	case LevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Determine output writer
	var output io.Writer
	switch logOutput {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		// Try to open file
		file, err := os.OpenFile(logOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// Fall back to stderr
			output = os.Stderr
		} else {
			output = file
		}
	}

	// Create handler based on format
	opts := &slog.HandlerOptions{Level: level}
	if logFormat == "json" {
		logHandler = slog.NewJSONHandler(output, opts)
	} else {
		logHandler = slog.NewTextHandler(output, opts)
	}

	// Create and store logger
	defaultLogger = slog.New(logHandler)
	return defaultLogger
}

// GetLogger returns the default logger or initializes a new one if it doesn't exist
func GetLogger() *slog.Logger {
	if defaultLogger == nil {
		return InitLogger()
	}
	return defaultLogger
}

// WithLevel returns a logger with the specified level
func WithLevel(level slog.Level) *slog.Logger {
	if logHandler == nil {
		InitLogger()
	}
	
	var handler slog.Handler
	if logFormat == "json" {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	}
	
	return slog.New(handler)
}