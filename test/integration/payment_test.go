package integration

import (
	"tazapay-mcp-server/internal/tazapay"
	"testing"
)

func TestPaymentCreation(t *testing.T) {
	client := setupClient(t)

	paymentReq := &tazapay.PaymentRequest{
		Amount:          2000.00,
		Currency:        "USD",
		InvoiceCurrency: "USD",
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
		t.Fatalf("Error creating payment: %v", err)
	}

	if response.Data.ID == "" {
		t.Error("Expected payment ID to be non-empty")
	}
	if response.Data.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", response.Data.Status)
	}
}
