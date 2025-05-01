package main

import (
	"encoding/base64"
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

	// Set up Viper
	viper.SetConfigName(".tazapay-mcp-server.yml") // Name of the config file (with extension)
	viper.SetConfigType("yaml")                    // Config file type

	// Add config paths in order of priority
	viper.AddConfigPath(filepath.Join(wd, "config")) // Look in config directory first
	viper.AddConfigPath(wd)                          // Then look in current directory
	viper.AddConfigPath(os.Getenv("HOME"))           // Finally look in home directory

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Print which config file is being used
	fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())

	// Retrieve the keys
	accessKey := viper.GetString("TAZAPAY_API_KEY")
	secretKey := viper.GetString("TAZAPAY_API_SECRET")

	// Combine accessKey and secretKey with a colon
	authString := fmt.Sprintf("%s:%s", accessKey, secretKey)

	// Encode the string to Base64
	authToken := base64.StdEncoding.EncodeToString([]byte(authString))
	viper.Set("TAZAPAY_AUTH_TOKEN", authToken)

	// Create Tazapay client
	return tazapay.NewClient(accessKey, secretKey), nil
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
