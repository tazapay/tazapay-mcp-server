package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/tazapay/tazapay-mcp-server/pkg/agent"
	"github.com/tazapay/tazapay-mcp-server/pkg/tazapay"
)

func initConfig() (string, string, error) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("error getting working directory: %w", err)
	}

	// Try both .yml and .yaml extensions
	configPaths := []string{
		filepath.Join(wd, ".tazapay-mcp-server.yml"),
		filepath.Join(wd, ".tazapay-mcp-server.yaml"),
	}

	var configFile string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configFile = path
			break
		}
	}

	if configFile == "" {
		return "", "", fmt.Errorf("config file not found. Tried: %v", configPaths)
	}

	fmt.Printf("Using config file: %s\n", configFile)

	// Read the config file directly
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return "", "", fmt.Errorf("error reading config file: %w", err)
	}

	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")

	if apiKey == "" || apiSecret == "" {
		return "", "", fmt.Errorf("api_key and api_secret must be set in config file")
	}

	fmt.Println("Successfully loaded configuration")
	return apiKey, apiSecret, nil
}

func main() {
	// Initialize configuration
	apiKey, apiSecret, err := initConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create Tazapay client
	client := tazapay.NewClient(apiKey, apiSecret)

	// Create AI agent
	agent := agent.NewAgent(client)

	// Create reader for user input
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Tazapay AI Assistant!")
	fmt.Println("I can help you with:")
	fmt.Println("- Checking account balance")
	fmt.Println("- Creating beneficiaries")
	fmt.Println("- Making payouts")
	fmt.Println("- Checking exchange rates")
	fmt.Println("\nType 'exit' to quit.")

	// Main interaction loop
	for {
		fmt.Print("\nHow can I help you? > ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		// Process the message using the AI agent
		response, err := agent.HandleMessage(context.Background(), input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Println(response)
	}
}
