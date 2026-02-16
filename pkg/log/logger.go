package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/tazapay/tazapay-mcp-server/constants"
)

// Config represents the logger configuration.
type Config struct {
	FilePath string // Custom file path; if empty, uses default
	Format   string // "text" or "json"; defaults to "text"
	Level    string // "debug", "info", "warn", "error"; defaults to "info"
}

// getDefaultLogPath returns a fallback log path near the executable.
func getDefaultLogPath() string {
	execPath, err := os.Executable()
	if err != nil {
		return filepath.Join(os.TempDir(), "logs", "tazapay-mcp-server.log")
	}
	execDir := filepath.Dir(execPath)

	return filepath.Join(execDir, "logs", "tazapay-mcp-server.log")
}

// getEnvLogPath returns the log file path from the LOG_FILE_PATH environment variable if set.
func getEnvLogPath() string {
	if envPath := Get("LOG_FILE_PATH"); envPath != "" {
		return envPath
	}

	return ""
}

// New creates a structured logger based on the given config.
// Priority for log file path:
//  1. LOG_FILE_PATH environment variable
//  2. cfg.FilePath
//  3. default path near executable
func New(cfg Config) (*slog.Logger, func(ctx context.Context), error) {
	logPath := getEnvLogPath()
	if logPath == "" {
		logPath = cfg.FilePath
	}

	if logPath == "" {
		logPath = getDefaultLogPath()
	}

	// Ensure the log directory exists.
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, constants.OpenFileMode)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Warning: Failed to open log file: %v\nFalling back to stderr\n", err)

		return fallbackLogger(cfg), func(context.Context) {}, err
	}

	// Create a multi-writer to write to both file and stdout
	multiWriter := io.MultiWriter(file, os.Stdout)

	handler := getHandler(cfg, multiWriter)
	logger := slog.New(handler)

	cleanup := func(ctx context.Context) {
		cerr := file.Close()
		if cerr != nil {
			logger.WarnContext(ctx, "Failed to close log file", "error", cerr.Error())
		}
	}

	logger.InfoContext(context.Background(), "Logs are stored in", "path", logPath)

	return logger, cleanup, nil
}

// fallbackLogger returns a stderr-based logger if file init fails.
func fallbackLogger(cfg Config) *slog.Logger {
	handler := getHandler(cfg, os.Stderr)
	return slog.New(handler)
}

// getHandler creates the appropriate slog handler.
func getHandler(cfg Config, out any) slog.Handler {
	opts := &slog.HandlerOptions{
		Level: parseLogLevel(cfg.Level),
	}

	switch strings.ToLower(cfg.Format) {
	case "json":
		if w, ok := out.(io.Writer); ok {
			return slog.NewJSONHandler(w, opts)
		}

		panic("getHandler: out is not an io.Writer for JSON format")

	default:
		if w, ok := out.(io.Writer); ok {
			return slog.NewTextHandler(w, opts)
		}

		panic("getHandler: out is not an io.Writer for text format")
	}
}

// parseLogLevel converts a string level to slog.Level.
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug

	case "warn":
		return slog.LevelWarn

	case "error":
		return slog.LevelError

	default:
		return slog.LevelInfo
	}
}
