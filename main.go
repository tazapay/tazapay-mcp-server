package main

import (
	"encoding/base64"
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
	"github.com/tazapay/tazapay-mcp-server/tools"
	"os"
)

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return
	}
	// Set up Viper
	viper.AddConfigPath(home)                  // Path to look for the config file
	viper.SetConfigName(".tazapay-mcp-server") // Name of the config file (without extension)
	viper.SetConfigType("yaml")                // Config file type
	// Path to look for the config file in the current directory

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}

	// Retrieve the keys
	accessKey := viper.GetString("TAZAPAY_API_KEY")
	secretKey := viper.GetString("TAZAPAY_API_SECRET")

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
	//Add tools to the server
	//tools.AddHelloTool(s)
	tools.AddAddTool(s)
	tools.AddFXTool(s)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}

	// Uncomment the following lines to test the Fxcall function
	//res, err := tools.Fxcall("USD", "INR", 1000)
	//if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	return
	//}
	//fmt.Printf("Result: %v\n", res)

}
