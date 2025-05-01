package server

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"tazapay-mcp-server/internal/tazapay"
)

// MCPServer represents the MCP server instance
type MCPServer struct {
	server *server.MCPServer
	client *tazapay.Client
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(client *tazapay.Client) *MCPServer {
	// Initialize MCP server with required parameters
	srv := server.NewMCPServer(
		"tazapay-mcp-server",
		"1.0.0",
		server.WithLogging(),
		server.WithRecovery(),
	)

	return &MCPServer{
		server: srv,
		client: client,
	}
}

// Start starts the MCP server
func (s *MCPServer) Start() error {
	// Register tools
	s.registerTools()

	// Start the server
	return server.ServeStdio(s.server)
}

// registerTools registers all available tools with the MCP server
func (s *MCPServer) registerTools() {
	// Register payment tool
	paymentTool := mcp.NewTool("create_payment")
	paymentTool.Description = "Creates a new payment using Tazapay Checkout API"
	paymentTool.InputSchema = mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"amount": map[string]interface{}{
				"type":        "number",
				"description": "The payment amount",
			},
			"currency": map[string]interface{}{
				"type":        "string",
				"description": "The payment currency (e.g., SGD)",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "The payment description",
			},
			"success_url": map[string]interface{}{
				"type":        "string",
				"description": "The URL to redirect to on successful payment",
			},
			"cancel_url": map[string]interface{}{
				"type":        "string",
				"description": "The URL to redirect to on cancelled payment",
			},
			"customer_email": map[string]interface{}{
				"type":        "string",
				"description": "The customer's email address",
			},
			"customer_name": map[string]interface{}{
				"type":        "string",
				"description": "The customer's name",
			},
			"customer_phone": map[string]interface{}{
				"type":        "string",
				"description": "The customer's phone number",
			},
			"customer_address": map[string]interface{}{
				"type":        "string",
				"description": "The customer's address",
			},
		},
		Required: []string{
			"amount",
			"currency",
			"description",
			"success_url",
			"cancel_url",
			"customer_email",
			"customer_name",
		},
	}
	s.server.AddTool(paymentTool, s.handlePaymentTool)

	// Register currency tool
	fxTool := mcp.NewTool("get_fx_rates")
	fxTool.Description = "Gets foreign exchange rates from Tazapay API"
	fxTool.InputSchema = mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"from_currency": map[string]interface{}{
				"type":        "string",
				"description": "The source currency (e.g., SGD)",
			},
			"to_currency": map[string]interface{}{
				"type":        "string",
				"description": "The target currency (e.g., USD)",
			},
			"amount": map[string]interface{}{
				"type":        "number",
				"description": "The amount to convert",
			},
		},
		Required: []string{
			"from_currency",
			"to_currency",
			"amount",
		},
	}
	s.server.AddTool(fxTool, s.handleCurrencyTool)
}

// handlePaymentTool handles payment creation requests
func (s *MCPServer) handlePaymentTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters from request
	params := request.Params.Arguments

	// Create payment request
	paymentReq := &tazapay.PaymentRequest{
		Amount:        params["amount"].(float64),
		Currency:      params["currency"].(string),
		Description:   params["description"].(string),
		SuccessURL:    params["success_url"].(string),
		CancelURL:     params["cancel_url"].(string),
		CustomerEmail: params["customer_email"].(string),
		CustomerName:  params["customer_name"].(string),
	}

	// Set customer details
	if phone, ok := params["customer_phone"].(string); ok {
		paymentReq.CustomerDetails.Phone.Number = phone
		paymentReq.CustomerDetails.Phone.CallingCode = "1" // Default to US
	}
	if address, ok := params["customer_address"].(string); ok {
		paymentReq.CustomerDetails.Address = address
	}

	// Create payment using client
	response, err := s.client.CreatePayment(paymentReq)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("failed to create payment", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Payment created successfully!\nCheckout URL: %s", response.Data.URL)), nil
}

// handleCurrencyTool handles currency conversion requests
func (s *MCPServer) handleCurrencyTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters from request
	params := request.Params.Arguments
	fromCurrency := params["from_currency"].(string)
	toCurrency := params["to_currency"].(string)
	amount := params["amount"].(float64)

	// Get exchange rate using client
	response, err := s.client.GetExchangeRate(fromCurrency, toCurrency, amount)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("failed to get FX rates", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Exchange rate from %s to %s: %.6f\nConverted amount: %.2f %s",
		response.Data.InitialCurrency,
		response.Data.FinalCurrency,
		response.Data.ExchangeRate,
		float64(response.Data.ConvertedAmount)/100,
		response.Data.FinalCurrency)), nil
}
