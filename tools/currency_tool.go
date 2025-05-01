package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"tazapay-mcp-server/pkg/tazapay"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// CurrencyResponse represents the response from the currency metadata endpoint
type CurrencyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Currencies []struct {
			Code        string `json:"code"`
			Name        string `json:"name"`
			Type        string `json:"type"`
			IsActive    bool   `json:"is_active"`
			IsDefault   bool   `json:"is_default"`
			CountryCode string `json:"country_code"`
		} `json:"currencies"`
	} `json:"data"`
}

// AddCurrencyTool registers the currency tool with the server
func AddCurrencyTool(s *server.MCPServer, client *tazapay.Client) {
	tool := mcp.NewTool("validate_currency",
		mcp.WithDescription("Validate if a currency is a valid holding currency"),
		mcp.WithString("currency",
			mcp.Description("Currency code to validate (e.g., USD, SGD)"),
		),
	)
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleCurrencyTool(ctx, request, client)
	})
}

func handleCurrencyTool(
	ctx context.Context,
	request mcp.CallToolRequest,
	client *tazapay.Client,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	currency, ok := arguments["currency"].(string)
	if !ok {
		return nil, fmt.Errorf("currency argument is required")
	}

	// Make API call to fetch valid currencies
	baseURL := "https://api-yellow.tazapay.com/v3/metadata/payoutcurrency"
	req, err := http.NewRequest("GET", baseURL, nil)
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

	var response CurrencyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Check if the currency is valid
	isValid := false
	var validCurrencies []string
	for _, c := range response.Data.Currencies {
		validCurrencies = append(validCurrencies, c.Code)
		if c.Code == currency && c.IsActive {
			isValid = true
			break
		}
	}

	// Prepare response
	var result string
	if isValid {
		result = fmt.Sprintf("Currency %s is a valid holding currency.", currency)
	} else {
		result = fmt.Sprintf("Currency %s is not a valid holding currency.\nValid holding currencies are: %v",
			currency, validCurrencies)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}
