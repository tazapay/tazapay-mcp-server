package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"tazapay-mcp-server/internal/tazapay"

	"github.com/spf13/viper"
)

func TestPaymentCreation(t *testing.T) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working directory: %v", err)
	}

	// Initialize configuration
	viper.SetConfigName(".tazapay-mcp-server")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join(wd, "..", "config")) // Look in the config directory
	viper.AddConfigPath(filepath.Join(wd, ".."))           // Look in project root
	viper.AddConfigPath(os.Getenv("HOME"))                 // Look in home directory

	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Error reading config: %v", err)
	}

	// Get API credentials
	apiKey := viper.GetString("TAZAPAY_API_KEY")
	apiSecret := viper.GetString("TAZAPAY_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		t.Fatal("API key or secret not found in config")
	}

	// Create Tazapay client
	client := tazapay.NewClient(apiKey, apiSecret)

	// Test payment creation with required fields
	paymentReq := &tazapay.PaymentRequest{
		Amount:          2000.00, // Amount in currency units
		Currency:        "BRL",
		InvoiceCurrency: "BRL",
		Description:     "Test payment for Tazapay integration",
		TransactionDesc: "Test payment for Tazapay integration",
		SuccessURL:      "https://example.com/success",
		CancelURL:       "https://example.com/cancel",
		CustomerEmail:   "test@example.com",
		CustomerName:    "Test User",
		CustomerDetails: struct {
			Email string `json:"email"`
			Name  string `json:"name"`
			Phone struct {
				Number      string `json:"number"`
				CallingCode string `json:"calling_code"`
			} `json:"phone"`
			Address string `json:"address"`
			Country string `json:"country"`
		}{
			Email:   "test@example.com",
			Name:    "Test User",
			Country: "BR",
			Phone: struct {
				Number      string `json:"number"`
				CallingCode string `json:"calling_code"`
			}{
				CallingCode: "1",
				Number:      "1234567890",
			},
			Address: "123 Test St, Test City, BR",
		},
		PaymentMethods: []string{"card", "bank_transfer"},
	}

	response, err := client.CreatePayment(paymentReq)
	if err != nil {
		t.Fatalf("Error creating payment: %v", err)
	}

	// Display payment details in a user-friendly format
	fmt.Println("\nPayment Details:")
	fmt.Println("---------------")
	fmt.Printf("Amount: %.2f %s\n", paymentReq.Amount, paymentReq.Currency)
	fmt.Printf("Description: %s\n", paymentReq.Description)
	fmt.Printf("Customer: %s (%s)\n", paymentReq.CustomerName, paymentReq.CustomerEmail)
	fmt.Println("\nHere is your payment link:")
	fmt.Println("------------------------")
	fmt.Printf("Checkout URL: %s\n", response.Data.URL)
	fmt.Println("\nYou can use this link to complete your payment.")
}
