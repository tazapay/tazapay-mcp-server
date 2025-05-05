package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
	"github.com/tazapay/tazapay-mcp-server/constants"
)

// AddPaymentLinkTool adds the paymentlink tool to the server
func AddPaymentLinkTool(s *server.MCPServer) {
	tool := mcp.NewTool("tazapaymentlinktool",
		mcp.WithDescription("Generate a payment link using Tazapay checkout API"),
		mcp.WithString("invoice_currency",
			mcp.Required(),
			mcp.Description("Currency for the payment. Should be in 3 letter currency code. Example: USD, INR"),
		),
		mcp.WithNumber("payment_amount",
			mcp.Required(),
			mcp.Description("Amount for the payment. Should be a float number."),
		),
		mcp.WithString("customer_name",
			mcp.Required(),
			mcp.Description("Name of the customer"),
		),
		mcp.WithString("customer_email",
			mcp.Required(),
			mcp.Description("Email address of the customer"),
		),
		mcp.WithString("customer_country",
			mcp.Required(),
			mcp.Description("Country of the customer in ISO 2-letter code. Example: US, IN"),
		),
		mcp.WithString("transaction_description",
			mcp.Required(),
			mcp.Description("Description of the transaction"),
		),
	)
	s.AddTool(tool, handlePaymentLinkTool)
}

// handlePaymentLinkTool handles the paymentlink tool call and performs checks validation
func handlePaymentLinkTool(ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments

	// Check if all required fields are present
	invoiceCurrency, param1 := arguments["invoice_currency"].(string)
	if invoiceCurrency == "" {
		return nil, fmt.Errorf("missing required field: invoice_currency")
	}

	amount, param2 := arguments["payment_amount"].(float64)
	if amount == 0 {
		return nil, fmt.Errorf("missing required field: amount")
	}

	customerName, param3 := arguments["customer_name"].(string)
	if customerName == "" {
		return nil, fmt.Errorf("missing required field: customer_name")
	}

	customerEmail, param4 := arguments["customer_email"].(string)
	if customerEmail == "" {
		return nil, fmt.Errorf("missing required field: customer_email")
	}

	customerCountry, param5 := arguments["customer_country"].(string)
	if customerCountry == "" {
		return nil, fmt.Errorf("missing required field: customer_country")
	}

	transactionDescription, param6 := arguments["transaction_description"].(string)
	if transactionDescription == "" {
		transactionDescription = constants.Miscellaneous
	}

	// Check if all required fields are present
	if !param1 || !param2 || !param3 || !param4 || !param5 || !param6 {
		return nil, fmt.Errorf("invalid arguments")
	}

	// Call the CreatePaymentLink function
	return CreatePaymentLink(invoiceCurrency, amount*100, customerName, customerEmail, customerCountry, transactionDescription)
}

// CreatePaymentLink creates a payment link
func CreatePaymentLink(invoiceCurrency string, amount float64, customerName string, customerEmail string, customerCountry string, transactionDescription string,
) (*mcp.CallToolResult, error) {
	baseURL := constants.PaymentLinkBaseURL_prod
	requestBody := map[string]any{
		"amount":                  amount,
		"invoice_currency":        invoiceCurrency,
		"transaction_description": transactionDescription,
		"customer_details": map[string]string{
			"name":    customerName,
			"email":   customerEmail,
			"country": customerCountry,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request body: %v", err)
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	authToken := "Basic " + viper.GetString("TAZAPAY_AUTH_TOKEN")
	req.Header.Set("Authorization", authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-success status code: %v, body: %s", resp.Status, string(bodyBytes))
	}

	var result map[string]any

	// Read and decode JSON directly
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error getting fx rate : %v\n", err)
	}

	data, _ := result["data"].(map[string]any)
	paymentLink, ok4 := data["url"].(string)
	if !ok4 {
		return nil, fmt.Errorf("error getting payment link")
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Payment Link URL: %v", paymentLink),
			},
		},
	}, nil
}
