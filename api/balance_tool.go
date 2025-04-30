package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
)

func AddBalanceTool(s *server.MCPServer) {
	tool := mcp.NewTool("check_balance",
		mcp.WithDescription("Check account balance for a specific currency"),
		mcp.WithString("currency",
			mcp.Required(),
			mcp.Description("Currency code (e.g., USD, EUR)"),
		),
	)
	s.AddTool(tool, handleBalanceTool)
}

func handleBalanceTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	currency, ok := arguments["currency"].(string)
	if !ok {
		return nil, fmt.Errorf("currency must be a string")
	}

	// Make API call to check balance
	baseURL := "https://api-yellow.tazapay.com/v3/balance"
	url := fmt.Sprintf("%s?currency=%s", baseURL, currency)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+viper.GetString("TAZAPAY_AUTH_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("received non-success status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Format the response in a user-friendly way
	balance := result["data"].(map[string]interface{})["balance"].(float64)
	formattedBalance := fmt.Sprintf("Your %s balance is %.2f", currency, balance)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: formattedBalance,
			},
		},
	}, nil
}
