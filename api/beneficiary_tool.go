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

func AddBeneficiaryTool(s *server.MCPServer) {
	tool := mcp.NewTool("create_beneficiary",
		mcp.WithDescription("Create a new beneficiary for payouts"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Full name of the beneficiary"),
		),
		mcp.WithString("email",
			mcp.Required(),
			mcp.Description("Email address of the beneficiary"),
		),
		mcp.WithString("country",
			mcp.Required(),
			mcp.Description("Country code of the beneficiary (e.g., US, IN)"),
		),
		mcp.WithString("bank_account_number",
			mcp.Required(),
			mcp.Description("Bank account number of the beneficiary"),
		),
		mcp.WithString("bank_code",
			mcp.Required(),
			mcp.Description("Bank code or routing number"),
		),
		mcp.WithString("bank_name",
			mcp.Required(),
			mcp.Description("Name of the beneficiary's bank"),
		),
	)
	s.AddTool(tool, handleBeneficiaryTool)
}

func handleBeneficiaryTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	name, ok1 := arguments["name"].(string)
	email, ok2 := arguments["email"].(string)
	country, ok3 := arguments["country"].(string)
	accountNumber, ok4 := arguments["bank_account_number"].(string)
	bankCode, ok5 := arguments["bank_code"].(string)
	bankName, ok6 := arguments["bank_name"].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return nil, fmt.Errorf("invalid arguments")
	}

	// Create beneficiary request
	beneficiaryRequest := map[string]interface{}{
		"name":                name,
		"email":               email,
		"country":             country,
		"bank_account_number": accountNumber,
		"bank_code":           bankCode,
		"bank_name":           bankName,
	}

	jsonData, err := json.Marshal(beneficiaryRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	// Make API call to create beneficiary
	baseURL := "https://api-yellow.tazapay.com/v3/beneficiary"
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
	beneficiaryID := result["data"].(map[string]interface{})["id"].(string)
	status := result["data"].(map[string]interface{})["status"].(string)
	formattedResponse := fmt.Sprintf("Beneficiary created successfully!\nID: %s\nStatus: %s", beneficiaryID, status)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: formattedResponse,
			},
		},
	}, nil
}
