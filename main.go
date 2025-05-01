package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"tazapay-mcp-server/internal/tazapay"
	"tazapay-mcp-server/pkg/chat"

	"github.com/spf13/viper"
)

func initConfig() (*tazapay.Client, error) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %v", err)
	}
	fmt.Printf("Current working directory: %s\n", wd)

	// Set up Viper with explicit yml extension
	viper.SetConfigName(".tazapay-mcp-server")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(wd, "config")) // Look in config directory
	viper.AddConfigPath(wd)                          // Also look in current directory as fallback

	// Try to read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Print which config file is being used
	fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())

	// Get API credentials
	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")

	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("API key or secret not found in config file")
	}

	// Create Tazapay client
	return tazapay.NewClient(apiKey, apiSecret), nil
}

func main() {
	// Initialize configuration and create client
	client, err := initConfig()
	if err != nil {
		fmt.Printf("Error initializing: %v\n", err)
		os.Exit(1)
	}

	// Create chat interface
	chatInterface := chat.NewChatInterface(client)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		os.Exit(0)
	}()

	// Start the chat interface
	fmt.Println("Starting Tazapay Chat Interface...")
	chatInterface.Run()
}
