package tools

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"io"
	"net/http"
	"os"
)

func AddFXTool(s *server.MCPServer) {
	tool := mcp.NewTool("tazapay_fx",
		mcp.WithDescription("Get FX rate from one currency to another using Tazapay FX rate"),
		mcp.WithString("from",
			mcp.Required(),
			mcp.Description("Currency to convert from. It should be in 3 letter currency code. Example : USD, INR"),
		),
		mcp.WithString("to",
			mcp.Required(),
			mcp.Description("Currency to convert to. It should be in 3 letter currency code. Example : USD, INR"),
		),
		mcp.WithNumber("amount",
			mcp.Required(),
			mcp.Description("Amount to convert"),
		),
	)
	s.AddTool(tool, handleFXTool)
}

func handleFXTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	// Extract the parameters from the request
	arguments := request.Params.Arguments
	from, ok1 := arguments["from"].(string)
	to, ok2 := arguments["to"].(string)
	amount, ok3 := arguments["amount"].(float64)

	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("invalid arguments")
	}

	// Perform the currency conversion (dummy implementation)
	convertedAmount := amount * 1.2 // Replace with actual conversion logic

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Converted %f %s to %f %s.", amount, from, convertedAmount, to),
			},
		},
	}, nil
}

func fxcall() {
	// API endpoint and parameters
	baseURL := "https://service.tazapay.com/v3/fx/payout"
	initialCurrency := "USD"
	finalCurrency := "INR"
	amount := "100"

	// Construct the full URL with query parameters
	url := fmt.Sprintf("%s?initial_currency=%s&final_currency=%s&amount=%s",
		baseURL, initialCurrency, finalCurrency, amount)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")

	// Replace YOUR_AUTH_TOKEN with your actual authorization token
	authToken := "YOUR_AUTH_TOKEN"
	req.Header.Set("Authorization", authToken)

	// Create HTTP client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read and print response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))
}
