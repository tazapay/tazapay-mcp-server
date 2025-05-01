package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config represents the server configuration
type Config struct {
	APIKey    string `yaml:"api_key"`
	APISecret string `yaml:"api_secret"`
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %v", err)
	}

	// Try to find the config file in the config directory first
	configPath := filepath.Join(wd, "config", ".tazapay-mcp-server.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Try the current directory
		configPath = filepath.Join(wd, ".tazapay-mcp-server.yml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found in %s or config directory", wd)
		}
	}

	fmt.Printf("Using config file: %s\n", configPath)

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Parse the config
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	// Validate the config
	if cfg.APIKey == "" || cfg.APISecret == "" {
		return nil, fmt.Errorf("API key and secret must be provided in config file")
	}

	return &cfg, nil
}
