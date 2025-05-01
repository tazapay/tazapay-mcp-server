package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tazapay-mcp-server/internal/config"
	"tazapay-mcp-server/internal/tazapay"
	"tazapay-mcp-server/pkg/chat"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize Tazapay client
	client := tazapay.NewClient(cfg.APIKey, cfg.APISecret)

	// Create chat interface
	chatInterface := chat.NewChatInterface(client)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down server...")
		os.Exit(0)
	}()

	// Start the chat interface
	fmt.Println("Starting Tazapay Chat Interface...")
	chatInterface.Run()
}
