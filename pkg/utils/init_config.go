package utils

import (
	"context"
	"encoding/base64"
	"errors"
	"os"

	"log/slog"

	"github.com/spf13/viper"

	"github.com/tazapay/tazapay-mcp-server/constants"
)

// InitConfig initializes configuration using viper and sets up auth token
func InitConfig(logger *slog.Logger) error {
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
				logger.ErrorContext(context.Background(), "Config read error", "error", readErr)
				return readErr
			}
		}
	}

	accessKey := viper.GetString("TAZAPAY_API_KEY")
	secretKey := viper.GetString("TAZAPAY_API_SECRET")

	if accessKey == "" || secretKey == "" {
		logger.ErrorContext(context.Background(), "Missing API credentials")
		return constants.ErrMissingAuthKeys
	}

	authString := accessKey + ":" + secretKey
	authToken := base64.StdEncoding.EncodeToString([]byte(authString))
	viper.Set("TAZAPAY_AUTH_TOKEN", authToken)

	logger.InfoContext(context.Background(), "Configuration initialized")

	return nil
}
