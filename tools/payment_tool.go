package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"tazapay-mcp-server/pkg/tazapay"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// PaymentRequest represents the request body for creating a payment
type PaymentRequest struct {
	Amount          float64 `json:"amount"`
	InvoiceCurrency string  `json:"invoice_currency"`
	Description     string  `json:"transaction_description"`
	SuccessURL      string  `json:"success_url"`
	CancelURL       string  `json:"cancel_url"`
	Customer        struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Country string `json:"country"`
		Phone   string `json:"phone,omitempty"`
		Address string `json:"address,omitempty"`
	} `json:"customer_details"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// PaymentResponse represents the response from the payment creation endpoint
type PaymentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID         string  `json:"id"`
		Amount     float64 `json:"amount"`
		Currency   string  `json:"currency"`
		Status     string  `json:"status"`
		PaymentURL string  `json:"payment_url"`
		CreatedAt  string  `json:"created_at"`
		ExpiresAt  string  `json:"expires_at"`
	} `json:"data"`
}

// AddPaymentTool registers the payment tool with the server
func AddPaymentTool(s *server.MCPServer, client *tazapay.Client) {
	tool := mcp.NewTool("create_payment",
		mcp.WithDescription("Create a payment using Tazapay Checkout API"),
		mcp.WithNumber("amount",
			mcp.Required(),
			mcp.Description("Amount to be paid"),
		),
		mcp.WithString("currency",
			mcp.Required(),
			mcp.Description("Currency code (e.g., USD, SGD)"),
		),
		mcp.WithString("description",
			mcp.Required(),
			mcp.Description("Description of the payment"),
		),
		mcp.WithString("customer_email",
			mcp.Required(),
			mcp.Description("Customer's email address"),
		),
		mcp.WithString("customer_name",
			mcp.Required(),
			mcp.Description("Customer's name"),
		),
		mcp.WithString("success_url",
			mcp.Required(),
			mcp.Description("URL to redirect after successful payment"),
		),
		mcp.WithString("cancel_url",
			mcp.Required(),
			mcp.Description("URL to redirect after cancelled payment"),
		),
		mcp.WithString("customer_phone",
			mcp.Description("Customer's phone number (optional)"),
		),
		mcp.WithString("customer_address",
			mcp.Description("Customer's address (optional)"),
		),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handlePaymentTool(ctx, request, client)
	})
}

func handlePaymentTool(
	ctx context.Context,
	request mcp.CallToolRequest,
	client *tazapay.Client,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments

	// Extract required fields
	amount, ok1 := arguments["amount"].(float64)
	currency, ok2 := arguments["currency"].(string)
	description, ok3 := arguments["description"].(string)
	customerEmail, ok4 := arguments["customer_email"].(string)
	customerName, ok5 := arguments["customer_name"].(string)
	successURL, ok6 := arguments["success_url"].(string)
	cancelURL, ok7 := arguments["cancel_url"].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return nil, fmt.Errorf("missing required fields")
	}

	// Create payment request
	paymentReq := PaymentRequest{
		Amount:          amount,
		InvoiceCurrency: currency,
		Description:     description,
		Customer: struct {
			Email   string `json:"email"`
			Name    string `json:"name"`
			Country string `json:"country"`
			Phone   string `json:"phone,omitempty"`
			Address string `json:"address,omitempty"`
		}{
			Email: customerEmail,
			Name:  customerName,
		},
		SuccessURL: successURL,
		CancelURL:  cancelURL,
	}

	// Add optional fields if provided
	if phone, ok := arguments["customer_phone"].(string); ok {
		paymentReq.Customer.Phone = phone
	}
	if address, ok := arguments["customer_address"].(string); ok {
		paymentReq.Customer.Address = address
	}

	// Make API call to create payment
	baseURL := "https://service-sandbox.tazapay.com/v3/checkout"
	reqBody, err := json.Marshal(paymentReq)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Basic "+client.GetAuthToken())

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-success status code: %d", resp.StatusCode)
	}

	var response PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Format response
	result := fmt.Sprintf("Payment created successfully!\n"+
		"ID: %s\n"+
		"Amount: %.2f %s\n"+
		"Status: %s\n"+
		"Payment URL: %s\n"+
		"Created At: %s\n"+
		"Expires At: %s",
		response.Data.ID,
		response.Data.Amount,
		response.Data.Currency,
		response.Data.Status,
		response.Data.PaymentURL,
		response.Data.CreatedAt,
		response.Data.ExpiresAt,
	)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}
