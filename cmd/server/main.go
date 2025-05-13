package main

import (
	"context"
	"encoding/base64"
	"errors"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"

	"github.com/tazapay/tazapay-mcp-server/constants"

	logs "github.com/tazapay/tazapay-mcp-server/pkg/logs"
	tools "github.com/tazapay/tazapay-mcp-server/tools/register"
)

func initConfig(logger *slog.Logger) error {
	viper.AutomaticEnv()

	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(home)
		viper.SetConfigName(".tazapay-mcp-server")
		viper.SetConfigType("yaml")

		readErr := viper.ReadInConfig()
		if readErr != nil {
			var notFoundErr viper.ConfigFileNotFoundError
			if !errors.As(readErr, &notFoundErr) {
				logger.Error("Config read error", "error", readErr)
				return readErr
			}
		}
	}

	accessKey := viper.GetString("TAZAPAY_API_KEY")
	secretKey := viper.GetString("TAZAPAY_API_SECRET")

	if accessKey == "" || secretKey == "" {
		logger.Error("Missing API credentials")
		return constants.ErrMissingAuthKeys
	}

	authString := accessKey + ":" + secretKey
	authToken := base64.StdEncoding.EncodeToString([]byte(authString))
	viper.Set("TAZAPAY_AUTH_TOKEN", authToken)

	logger.Info("Configuration initialized")

	return nil
}

func main() {
	// Create a logger configuration
	logConfig := logs.Config{
		Level:    "info",                           // Example log level
		Format:   "json",                           // Example log format
		FilePath: viper.GetString("LOG_FILE_PATH"), // Optional file path for logs, if needed
	}

	// Create the logger
	logger, cleanup, err := logs.New(logConfig) // Empty path = default path near binary
	defer cleanup(context.Background())

	// Exit if logger failed
	if err != nil {
		os.Exit(1)
	}

	if err := initConfig(logger); err != nil {
		logger.Error("failed to initialize config", "error", err)
		os.Exit(1)
	}

	s := server.NewMCPServer("tazapay", "0.0.1")

	tools.RegisterTools(s, logger)

	logger.Info("Started Tazapay MCP Server.")

	if err := server.ServeStdio(s); err != nil {
		logger.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}
