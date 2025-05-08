package main

import (
	"encoding/base64"
	"errors"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"

	"github.com/tazapay/tazapay-mcp-server/constants"

	tools "github.com/tazapay/tazapay-mcp-server/tools/register"
)

func initConfig() error {
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
				// Config file error, ignoring since fallback is ENV
				return readErr
			}
		}
	}

	accessKey := viper.GetString("TAZAPAY_API_KEY")
	secretKey := viper.GetString("TAZAPAY_API_SECRET")

	if accessKey == "" || secretKey == "" {
		return constants.ErrMissingAuthKeys
	}

	authString := accessKey + ":" + secretKey
	authToken := base64.StdEncoding.EncodeToString([]byte(authString))
	viper.Set("TAZAPAY_AUTH_TOKEN", authToken)

	return nil
}

func main() {
	var err error
	if err = initConfig(); err != nil {
		os.Exit(1)
	}

	s := server.NewMCPServer(
		"tazapay",
		"0.0.1",
	)

	tools.RegisterTools(s)

	if err = server.ServeStdio(s); err != nil {
		os.Exit(1)
	}
}
