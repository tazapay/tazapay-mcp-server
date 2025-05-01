package integration

import (
	"os"
	"path/filepath"
	"tazapay-mcp-server/internal/tazapay"
	"testing"

	"github.com/spf13/viper"
)

func TestAllTools(t *testing.T) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working directory: %v", err)
	}

	// Initialize configuration
	viper.SetConfigName(".tazapay-mcp-server.yml")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join(wd, "..", "..", "config"))
	viper.AddConfigPath(filepath.Join(wd, "..", ".."))
	viper.AddConfigPath(wd)

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

	// Test 1: Check Balance
	t.Run("Check Balance", func(t *testing.T) {
		testCheckBalance(t, client)
	})

	// Test 2: Create Payment
	t.Run("Create Payment", func(t *testing.T) {
		testCreatePayment(t, client)
	})

	// Test 3: Create Beneficiary
	t.Run("Create Beneficiary", func(t *testing.T) {
		beneficiaryID := testCreateBeneficiary(t, client)
		if beneficiaryID != "" {
			// Test 4: Create Payout (requires beneficiary)
			t.Run("Create Payout", func(t *testing.T) {
				testCreatePayout(t, client, beneficiaryID)
			})
		}
	})

	// Test 5: Get FX Rates
	t.Run("Get FX Rates", func(t *testing.T) {
		testGetFXRates(t, client)
	})
}

func testCheckBalance(t *testing.T, client *tazapay.Client) {
	balance, err := client.GetBalance()
	if err != nil {
		t.Errorf("❌ Balance check failed: %v", err)
		return
	}

	t.Log("✅ Balance check successful")
	t.Log("Available balances:")
	for _, bal := range balance.Data.Available {
		t.Logf("- %s %s", bal.Amount, bal.Currency)
	}
	t.Logf("Last updated: %s", balance.Data.UpdatedAt)
}

func testCreatePayment(t *testing.T, client *tazapay.Client) {
	paymentReq := &tazapay.PaymentRequest{
		Amount:          100.00,
		Currency:        "USD",
		InvoiceCurrency: "USD",
		Description:     "Test payment creation",
		TransactionDesc: "Test payment creation",
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
			Country: "US",
			Phone: struct {
				Number      string `json:"number"`
				CallingCode string `json:"calling_code"`
			}{
				CallingCode: "1",
				Number:      "1234567890",
			},
			Address: "123 Test St, New York, US",
		},
		PaymentMethods: []string{"card"},
	}

	response, err := client.CreatePayment(paymentReq)
	if err != nil {
		t.Errorf("❌ Payment creation failed: %v", err)
		return
	}

	t.Log("✅ Payment creation successful")
	t.Logf("Payment ID: %s", response.Data.ID)
	t.Logf("Status: %s", response.Data.Status)
	t.Logf("Checkout URL: %s", response.Data.URL)
}

func testCreateBeneficiary(t *testing.T, client *tazapay.Client) string {
	beneficiary := &tazapay.Beneficiary{
		AccountID: "acc_d00eij4qqkhc9e5ats4g",
		Name:      "Test Beneficiary",
		Type:      "individual",
		Email:     "test.beneficiary@example.com",
		DestinationDetails: tazapay.DestinationDetails{
			Type: "bank",
			Bank: tazapay.BankDetails{
				Country:       "US",
				Currency:      "USD",
				BankName:      "Test Bank",
				AccountNumber: "1234567890",
				BankCodes: struct {
					SwiftCode string `json:"swift_code"`
				}{
					SwiftCode: "TESTUS33",
				},
			},
		},
		Phone: tazapay.Phone{
			CallingCode: "1",
			Number:      "1234567890",
		},
	}

	response, err := client.CreateBeneficiary(beneficiary)
	if err != nil {
		t.Errorf("❌ Beneficiary creation failed: %v", err)
		return ""
	}

	t.Log("✅ Beneficiary creation successful")
	t.Logf("Beneficiary ID: %s", response.Data.ID)
	return response.Data.ID
}

func testCreatePayout(t *testing.T, client *tazapay.Client, beneficiaryID string) {
	payout := &tazapay.Payout{
		BeneficiaryDetails: tazapay.BeneficiaryDetails{
			Name: "Test Beneficiary",
			Type: "individual",
			DestinationDetails: struct {
				Type string              `json:"type"`
				Bank tazapay.BankDetails `json:"bank"`
			}{
				Type: "bank",
				Bank: tazapay.BankDetails{
					BankName: "Test Bank",
					Currency: "USD",
					Country:  "US",
				},
			},
		},
		Purpose:             "PYR002",
		Amount:              5000,
		Currency:            "USD",
		Beneficiary:         beneficiaryID,
		HoldingCurrency:     "USD",
		Type:                "local",
		ChargeType:          "shared",
		StatementDescriptor: "Test payout",
	}

	response, err := client.CreatePayout(payout)
	if err != nil {
		t.Errorf("❌ Payout creation failed: %v", err)
		return
	}

	t.Log("✅ Payout creation successful")
	t.Logf("Payout ID: %s", response.Data.ID)
	t.Logf("Status: %s", response.Data.Status)
	t.Logf("Amount: %d %s", response.Data.Amount, response.Data.Currency)
}

func testGetFXRates(t *testing.T, client *tazapay.Client) {
	fromCurrency := "USD"
	toCurrency := "EUR"
	amount := 1000.00

	rate, err := client.GetExchangeRate(fromCurrency, toCurrency, amount)
	if err != nil {
		t.Errorf("❌ FX rate check failed: %v", err)
		return
	}

	t.Log("✅ FX rate check successful")
	t.Logf("From: %s", rate.Data.InitialCurrency)
	t.Logf("To: %s", rate.Data.FinalCurrency)
	t.Logf("Amount: %.2f %s", float64(rate.Data.Amount)/100, rate.Data.InitialCurrency)
	t.Logf("Exchange Rate: %.6f", rate.Data.ExchangeRate)
	t.Logf("Converted Amount: %.2f %s", float64(rate.Data.ConvertedAmount)/100, rate.Data.FinalCurrency)
	t.Logf("Timestamp: %s", rate.Data.Timestamp)
}
