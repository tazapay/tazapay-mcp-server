package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
)

func AddPaymentLinkTool(s *server.MCPServer) {
	tool := mcp.NewTool("tazapaymentlinktool",
		mcp.WithDescription("Generate a payment link using Tazapay checkout API"),
		mcp.WithString("invoice_currency",
			mcp.Required(),
			mcp.Description("Currency for the payment. Should be in 3 letter currency code. Example: USD, INR"),
		),
		mcp.WithNumber("payment amount",
			mcp.Required(),
			mcp.Description("Amount for the payment. Should be a number."),
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

func handlePaymentLinkTool(ctx context.Context,request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	invoiceCurrency, param1 := arguments["invoice_currency"].(string)
	if invoiceCurrency == "" {
		return nil, fmt.Errorf("missing required field: invoice_currency")
	}
	amount, param2 := arguments["amount"].(float64)
	if amount == 0  {
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
		return nil, fmt.Errorf("missing required field: transaction_description")
	}

	if !param1 || !param2 || !param3 || !param4 || !param5 || !param6 {
		return nil, fmt.Errorf("invalid arguments")
	}

	return CreatePaymentLink(invoiceCurrency, amount, customerName, customerEmail, customerCountry, transactionDescription)
}

func CreatePaymentLink(invoiceCurrency string,amount float64,customerName string,customerEmail string,customerCountry string,transactionDescription string,
) (*mcp.CallToolResult, error) {
	baseURL := "https://service.tazapay.com/v3/checkout"
	requestBody := map[string]interface{}{
		"invoice_currency": invoiceCurrency,
		"amount":           amount,
		"customer_details": map[string]string{
			"name":    customerName,
			"email":   customerEmail,
			"country": customerCountry,
		},
		"transaction_description": transactionDescription,
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
		return nil, fmt.Errorf("received non-success status code: %d", resp.StatusCode)
	}

	// need to check response body to check returned payload and then parse it 

	return &mcp.CallToolResult{}, nil
}
