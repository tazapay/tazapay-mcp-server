package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
	"github.com/tazapay/tazapay-mcp-server/constants"
)

func AddFXTool(s *server.MCPServer) {
	tool := mcp.NewTool("tazafx",
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
			mcp.Description("Amount to convert. It should be a number and should not have any decimal places."),
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
	// amount, ok3 := arguments["amount"].(int)
	amount, ok3 := 1000, true

	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("invalid arguments")
	}

	return Fxcall(from, to, amount)

}

func Fxcall(from string, to string, amount int) (*mcp.CallToolResult, error) {
	// API endpoint and parameters
	baseURL := constants.PaymentFxBaseURL_orange

	// Construct the full URL with query parameters
	url := fmt.Sprintf("%s?initial_currency=%s&final_currency=%s&amount=%d", baseURL, from, to, amount)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")

	authToken := "Basic " + viper.GetString("TAZAPAY_AUTH_TOKEN")
	req.Header.Set("Authorization", authToken)

	// Create HTTP client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("received non-success status code: %d", resp.StatusCode)
	}

	var result map[string]interface{} // or define your own struct

	// Read and decode JSON directly
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error getting fx rate : %v\n", err)
	}
	var data = result["data"].(map[string]interface{})
	// Extract the exchange rate and converted amount
	exchangeRate := data["exchange_rate"].(float64)
	convertedAmount := data["converted_amount"].(float64)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("FX rate from %v to %v is %f and converted amount is %v.", from, to, exchangeRate, convertedAmount),
			},
		},
	}, nil
}
