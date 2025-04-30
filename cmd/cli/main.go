package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/tazapay/tazapay-mcp-server/pkg/agent"
	"github.com/tazapay/tazapay-mcp-server/pkg/tazapay"
)

func initConfig() (string, string, error) {
	viper.SetConfigName(".tazapay-mcp-server")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return "", "", fmt.Errorf("error reading config file: %w", err)
	}

	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")

	if apiKey == "" || apiSecret == "" {
		return "", "", fmt.Errorf("api_key and api_secret must be set in config file")
	}

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
