package log_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	log "github.com/tazapay/tazapay-mcp-server/pkg/logs"
)

func TestNewLoggerWithDefaultConfig(t *testing.T) {
	cfg := log.Config{}
	logger, closeFn, err := log.New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	defer closeFn(t.Context())

	logger.InfoContext(t.Context(), "default logger test")
}

func TestNewLoggerWithCustomConfig(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "custom.log")

	cfg := log.Config{
		FilePath: logPath,
		Format:   "json",
		Level:    "debug",
	}

	logger, closeFn, err := log.New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	defer closeFn(t.Context())

	ctx := t.Context()
	logger.DebugContext(ctx, "debug message")
	logger.InfoContext(ctx, "info message")

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "debug message") || !strings.Contains(content, "info message") {
		t.Errorf("log file does not contain expected messages:\n%s", content)
	}
}

func TestNewLoggerWithInvalidFilePath(t *testing.T) {
	// Redirect os.Stderr to capture fallback log output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe() //nolint: errcheck // No need to check error
	os.Stderr = w

	cfg := log.Config{
		FilePath: "/invalid/path/to/logfile.log",
		Format:   "text",
		Level:    "info",
	}

	logger, closeFn, err := log.New(cfg)
	if err != nil {
		t.Fatalf("expected fallback logger without error, got: %v", err)
	}
	defer closeFn(t.Context())

	logger.InfoContext(t.Context(), "fallback log test")

	// Close writer and restore stderr
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		t.Fatalf("failed to read from stderr pipe: %v", err)
	}

	if !strings.Contains(buf.String(), "fallback") {
		t.Errorf("expected fallback warning in stderr, got: %s", buf.String())
	}
}
