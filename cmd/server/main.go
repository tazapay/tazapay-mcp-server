package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
	"github.com/tazapay/tazapay-mcp-server/tools"
)

func initConfig() {
	// Get the env variables if input from CLI
	viper.AutomaticEnv()

	// Try loading config file if exists (for local dev)
	home, err := os.UserHomeDir()
	if err == nil {
		// Set up Viper
		viper.AddConfigPath(home)                  // Path to look for the config file
		viper.SetConfigName(".tazapay-mcp-server") // Name of the config file (without extension)
		viper.SetConfigType("yaml")                // Config file type
		_ = viper.ReadInConfig()                   // ignore error if not found and check for
	}

	// Retrieve the keys
	accessKey := viper.GetString("TAZAPAY_API_KEY")
	secretKey := viper.GetString("TAZAPAY_API_SECRET")

	// return Error if env variables is not passed during runtime
	if accessKey == "" || secretKey == "" {
		fmt.Println("TAZAPAY_API_KEY or TAZAPAY_API_SECRET not set. Enter the following in terminal command with [-e] OPTION or add a `.tazapay-mcp-server.yaml` file in your home directory. ")
		os.Exit(1)
	}

	// Combine accessKey and secretKey with a colon
	authString := fmt.Sprintf("%s:%s", accessKey, secretKey)

	// Encode the string to Base64
	authToken := base64.StdEncoding.EncodeToString([]byte(authString))
	viper.Set("TAZAPAY_AUTH_TOKEN", authToken)
}

func main() {
	initConfig()
	// Create MCP server
	s := server.NewMCPServer(
		"tazapay",
		"0.0.1",
	)

	//Add FX tools to the server
	tools.AddFXTool(s)

	// added tool to generate payment link
	tools.AddPaymentLinkTool(s)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
