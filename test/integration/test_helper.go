package integration

import (
	"os"
	"path/filepath"
	"tazapay-mcp-server/internal/tazapay"
	"testing"

	"github.com/spf13/viper"
)

func setupClient(t *testing.T) *tazapay.Client {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working directory: %v", err)
	}

	viper.SetConfigName(".tazapay-mcp-server.yml")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join(wd, "..", "..", "config"))
	viper.AddConfigPath(filepath.Join(wd, "..", ".."))
	viper.AddConfigPath(wd)

	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Error reading config: %v", err)
	}

	apiKey := viper.GetString("TAZAPAY_API_KEY")
	apiSecret := viper.GetString("TAZAPAY_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		t.Fatal("API key or secret not found in config")
	}

	return tazapay.NewClient(apiKey, apiSecret)
}
