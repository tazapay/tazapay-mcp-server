package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"tazapay-mcp-server/internal/tazapay"

	"github.com/spf13/viper"
)

func TestPaymentCreation(t *testing.T) {
	// Get the absolute path to the config directory
	configPath, err := filepath.Abs("../config")
	if err != nil {
		t.Fatalf("Error getting config path: %v", err)
	}

	// Initialize configuration
	viper.SetConfigName(".tazapay-mcp-server")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Error reading config: %v", err)
	}

	// Get API credentials
	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")

	if apiKey == "" || apiSecret == "" {
		t.Fatal("API key or secret not found in config")
	}

	// Create Tazapay client
	client := tazapay.NewClient(apiKey, apiSecret)

	// Test payment creation with required fields
	paymentReq := &tazapay.PaymentRequest{
		Amount:          100.00,
		Currency:        "USD", // Using ISO 4217 alpha-3 currency code
		InvoiceCurrency: "USD", // Same as Currency for this test
		Description:     "Test payment",
		TransactionDesc: "Test payment for Tazapay integration",
		SuccessURL:      "https://example.com/success",
		CancelURL:       "https://example.com/cancel",
		CustomerEmail:   "test@example.com",
		CustomerName:    "Test User",
		PaymentMethods:  []string{"card", "wire_transfer"},
	}

	// Set customer details
	paymentReq.CustomerDetails.Email = paymentReq.CustomerEmail
	paymentReq.CustomerDetails.Name = paymentReq.CustomerName
	paymentReq.CustomerDetails.Phone.Number = "1234567890"
	paymentReq.CustomerDetails.Phone.CallingCode = "1" // US country code
	paymentReq.CustomerDetails.Address = "123 Test St, Test City"
	paymentReq.CustomerDetails.Country = "US"

	response, err := client.CreatePayment(paymentReq)
	if err != nil {
		t.Fatalf("Error creating payment: %v", err)
	}

	// Display payment details in a user-friendly format
	fmt.Println("\nPayment Details:")
	fmt.Println("---------------")
	fmt.Printf("Amount: $%.2f %s\n", paymentReq.Amount, paymentReq.Currency)
	fmt.Printf("Description: %s\n", paymentReq.TransactionDesc)
	fmt.Printf("Customer: %s (%s)\n", paymentReq.CustomerName, paymentReq.CustomerEmail)
	fmt.Println("\nHere is your payment link:")
	fmt.Println("------------------------")
	fmt.Printf("Checkout URL: %s\n", response.Data.URL)
	fmt.Println("\nYou can use this link to complete your payment.")
}
