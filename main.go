package main

import (
	"fmt"
	"github.com/tazapay/tazapay-mcp-server/tools"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"tazapay",
		"1.0.0",
	)
	//Add tools to the server
	//tools.AddHelloTool(s)
	tools.AddAddTool(s)
	tools.AddFXTool(s)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
