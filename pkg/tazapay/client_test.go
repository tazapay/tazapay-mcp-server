package tazapay

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func init() {
	// Set up Viper to read config file
	viper.SetConfigName(".tazapay-mcp-server")
	viper.SetConfigType("yaml")

	// Add multiple config paths
	viper.AddConfigPath(".")                                                       // Current directory
	viper.AddConfigPath("../..")                                                   // Project root
	viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), "go/tazapay-mcp-server")) // Full path

	// Try to read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Try with .yml extension
			viper.SetConfigName(".tazapay-mcp-server.yml")
			if err := viper.ReadInConfig(); err != nil {
				panic("Error reading config file: " + err.Error())
			}
		} else {
			panic("Error reading config file: " + err.Error())
		}
	}
}

func TestGetExchangeRate(t *testing.T) {
	// Initialize client with credentials from config
	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")

	if apiKey == "" || apiSecret == "" {
		t.Fatal("API credentials not found in config file")
	}

	client := NewClient(apiKey, apiSecret)

	// Test case 1: USD to EUR
	rate, err := client.GetExchangeRate("USD", "EUR", 100.0)
	if err != nil {
		t.Errorf("Error getting exchange rate: %v", err)
	}
	if rate == nil {
		t.Error("Expected non-nil exchange rate response")
	}

	// Test case 2: Invalid currency
	_, err = client.GetExchangeRate("INVALID", "EUR", 100.0)
	if err == nil {
		t.Error("Expected error for invalid currency")
	}
}

func TestCreatePayout(t *testing.T) {
	// Initialize client with credentials from config
	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")

	if apiKey == "" || apiSecret == "" {
		t.Fatal("API credentials not found in config file")
	}

	client := NewClient(apiKey, apiSecret)

	// Test case 1: Valid payout
	payout := &Payout{
		BeneficiaryDetails: BeneficiaryDetails{
			Name: "Test Beneficiary",
			Type: "individual",
			DestinationDetails: struct {
				Type string      `json:"type"`
				Bank BankDetails `json:"bank"`
			}{
				Type: "bank",
				Bank: BankDetails{
					Country:       "VN",
					Currency:      "VND",
					BankName:      "Test Bank",
					AccountNumber: "1234567890",
					BankCodes: struct {
						SwiftCode string `json:"swift_code"`
					}{
						SwiftCode: "TEST123",
					},
				},
			},
		},
		Purpose:             "PYR002",
		Amount:              1000,
		Currency:            "SGD",
		Beneficiary:         "bnf_d097upsjlbuun90ei8o0",
		HoldingCurrency:     "SGD",
		Type:                "local",
		ChargeType:          "shared",
		StatementDescriptor: "Test Payout",
	}

	response, err := client.CreatePayout(payout)
	if err != nil {
		t.Errorf("Error creating payout: %v", err)
	}
	if response == nil {
		t.Error("Expected non-nil payout response")
	}

	// Test case 2: Invalid payout (missing required fields)
	invalidPayout := &Payout{
		Amount: 1000,
	}
	_, err = client.CreatePayout(invalidPayout)
	if err == nil {
		t.Error("Expected error for invalid payout")
	}
}

func TestGetBalance(t *testing.T) {
	// Initialize client with credentials from config
	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")

	if apiKey == "" || apiSecret == "" {
		t.Fatal("API credentials not found in config file")
	}

	client := NewClient(apiKey, apiSecret)

	// Test case 1: Get balance
	balance, err := client.GetBalance()
	if err != nil {
		t.Errorf("Error getting balance: %v", err)
	}
	if balance == nil {
		t.Error("Expected non-nil balance response")
	}
}
