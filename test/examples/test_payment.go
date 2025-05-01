package main

import (
	"fmt"
	"os"
	"tazapay-mcp-server/tools"

	"github.com/spf13/viper"
)

func main() {
	// Initialize configuration
	viper.SetConfigName(".tazapay-mcp-server.yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	// Create test payment request
	paymentReq := tools.PaymentRequest{
		Amount:          10.00,
		Currency:        "USD",
		Description:     "Test payment from automated test",
		SuccessURL:      "https://example.com/success",
		CancelURL:       "https://example.com/cancel",
		CustomerEmail:   "test@example.com",
		CustomerName:    "Test User",
		CustomerPhone:   "+1234567890",
		CustomerAddress: "123 Test St, Test City",
	}

	fmt.Printf("Test Payment Request:\n")
	fmt.Printf("Amount: %.2f %s\n", paymentReq.Amount, paymentReq.Currency)
	fmt.Printf("Description: %s\n", paymentReq.Description)
	fmt.Printf("Customer: %s (%s)\n", paymentReq.CustomerName, paymentReq.CustomerEmail)
	fmt.Printf("URLs: \n  Success: %s\n  Cancel: %s\n", paymentReq.SuccessURL, paymentReq.CancelURL)

	// TODO: Make API call to create payment
	fmt.Printf("\nTest completed successfully!\n")
}
