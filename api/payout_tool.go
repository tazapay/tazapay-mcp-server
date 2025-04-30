package api

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

func AddPayoutTool(s *server.MCPServer) {
	tool := mcp.NewTool("create_payout",
		mcp.WithDescription("Create a new payout to a beneficiary"),
		mcp.WithString("beneficiary_id",
			mcp.Required(),
			mcp.Description("ID of the beneficiary to receive the payout"),
		),
		mcp.WithString("currency",
			mcp.Required(),
			mcp.Description("Currency code for the payout (e.g., USD, EUR)"),
		),
		mcp.WithNumber("amount",
			mcp.Required(),
			mcp.Description("Amount to payout"),
		),
		mcp.WithString("description",
			mcp.Required(),
			mcp.Description("Description of the payout"),
		),
	)
	s.AddTool(tool, handlePayoutTool)
}

func handlePayoutTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	beneficiaryID, ok1 := arguments["beneficiary_id"].(string)
	currency, ok2 := arguments["currency"].(string)
	amount, ok3 := arguments["amount"].(float64)
	description, ok4 := arguments["description"].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		return nil, fmt.Errorf("invalid arguments")
	}

	// Create payout request
	payoutRequest := map[string]interface{}{
		"beneficiary_id": beneficiaryID,
		"currency":       currency,
		"amount":         amount,
		"description":    description,
	}

	jsonData, err := json.Marshal(payoutRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	// Make API call to create payout
	baseURL := "https://api-yellow.tazapay.com/v3/payout"
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
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

	// Format the response
	payoutID := result["data"].(map[string]interface{})["id"].(string)
	status := result["data"].(map[string]interface{})["status"].(string)
	formattedResponse := fmt.Sprintf("Payout created successfully!\nID: %s\nStatus: %s", payoutID, status)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: formattedResponse,
			},
		},
	}, nil
}
