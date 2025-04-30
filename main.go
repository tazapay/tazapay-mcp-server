package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"

	"github.com/tazapay/tazapay-mcp-server/api"
	"github.com/tazapay/tazapay-mcp-server/logs"
	"github.com/tazapay/tazapay-mcp-server/tools"
)

func initConfig() {
	viper.SetConfigName(".tazapay-mcp-server")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	// Read API keys from config
	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")

	// Generate Base64 encoded auth token
	authToken := base64.StdEncoding.EncodeToString([]byte(apiKey + ":" + apiSecret))
	os.Setenv("TZP_AUTH_TOKEN", authToken)
}

func main() {
	// Initialize configuration
	initConfig()

	// Create new MCP server
	s := server.NewMCPServer("tazapay", "0.0.1")

	// Add API tools
	api.AddBalanceTool(s)
	api.AddBeneficiaryTool(s)
	api.AddPayoutTool(s)

	// Add log analysis tool
	logs.AddLogAnalysisTool(s)

	// Add FX tool
	tools.AddFXTool(s)

	// Add add tool
	tools.AddAddTool(s)

	// Start server using stdio transport
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
