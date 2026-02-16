package log_test

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	log "github.com/tazapay/tazapay-mcp-server/pkg/log"
)

func getTestContext(t *testing.T) context.Context {
	t.Helper()
	// Go 1.20+ provides t.Context(), otherwise fallback
	tCtxMethod := reflect.ValueOf(t).MethodByName("Context")
	if tCtxMethod.IsValid() {
		results := tCtxMethod.Call(nil)
		if len(results) == 1 {
			if ctx, ok := results[0].Interface().(context.Context); ok {
				return ctx
			}
		}
	}
	return t.Context()
}

func TestNewLoggerWithDefaultConfig(t *testing.T) {
	cfg := log.Config{}
	logger, closeFn, err := log.New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	defer closeFn(getTestContext(t))

	logger.InfoContext(getTestContext(t), "default logger test")
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
	defer closeFn(getTestContext(t))

	ctx := getTestContext(t)
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

func TestLoggerUsesEnvLogFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "envlog.log")

	err := log.Set("LOG_FILE_PATH", logPath)
	if err != nil {
		t.Fatalf("failed to set LOG_FILE_PATH: %v", err)
	}
	defer os.Unsetenv("LOG_FILE_PATH")

	cfg := log.Config{
		FilePath: "should_not_be_used.log",
		Format:   "text",
		Level:    "info",
	}

	logger, closeFn, err := log.New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	defer closeFn(getTestContext(t))

	ctx := getTestContext(t)
	logger.InfoContext(ctx, "env log path message")

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "env log path message") {
		t.Errorf("log file does not contain expected message from env path:\n%s", content)
	}
}
