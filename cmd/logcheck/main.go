package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// last9Config holds the configuration for Last9 API
type last9Config struct {
	accessToken  string
	refreshToken string
	baseURL      string
}

func main() {
	// create mcp server instance for log checking
	s := server.NewMCPServer(
		"tazapay-logs",
		"0.0.1",
	)

	// initialize last9 configuration
	config := last9Config{
		accessToken:  os.Getenv("LAST9_ACCESS_TOKEN"),
		refreshToken: os.Getenv("LAST9_REFRESH_TOKEN"),
		baseURL:      "https://app.last9.io/api/v4/organizations/tazapay",
	}

	// register log checking tool
	s.AddTool(mcp.NewTool("check_logs",
		mcp.WithDescription("check logs from last9"),
		mcp.WithString("query",
			mcp.Description("log query to search for"),
		),
		mcp.WithNumber("time_range",
			mcp.Description("time range in minutes to check logs"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleLogCheck(ctx, request, config)
	})

	fmt.Println("starting tazapay log checking server...")

	// start server using standard input/output for communication
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("server error: %v\n", err)
		os.Exit(1)
	}
}

// handleLogCheck processes log checking requests
func handleLogCheck(
	ctx context.Context,
	request mcp.CallToolRequest,
	config last9Config,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	query, _ := arguments["query"].(string)
	timeRange, _ := arguments["time_range"].(float64)

	// create http client with auth
	client := &http.Client{}
	req, err := http.NewRequest("GET", config.baseURL+"/logs", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// set headers
	req.Header.Set("Authorization", "Bearer "+config.accessToken)
	req.Header.Set("Accept", "application/json")

	// add query parameters
	q := req.URL.Query()
	q.Add("query", query)
	q.Add("time_range", fmt.Sprintf("%.0f", timeRange))
	req.URL.RawQuery = q.Encode()

	// make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode == http.StatusUnauthorized {
		// try to refresh token
		newToken, err := refreshToken(config.refreshToken)
		if err != nil {
			return nil, fmt.Errorf("error refreshing token: %v", err)
		}
		config.accessToken = newToken
		// retry request with new token
		return handleLogCheck(ctx, request, config)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-success status code: %d", resp.StatusCode)
	}

	// parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// format response
	formattedResult := formatLogResult(result)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: formattedResult,
			},
		},
	}, nil
}

// refreshToken attempts to get a new access token using refresh token
func refreshToken(refreshToken string) (string, error) {
	// create http client
	client := &http.Client{}

	// create request to refresh token
	req, err := http.NewRequest("POST", "https://app.last9.io/api/v4/auth/refresh", nil)
	if err != nil {
		return "", fmt.Errorf("error creating refresh request: %v", err)
	}

	// set headers
	req.Header.Set("Authorization", "Bearer "+refreshToken)
	req.Header.Set("Accept", "application/json")

	// make request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making refresh request: %v", err)
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-success status code: %d", resp.StatusCode)
	}

	// parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding refresh response: %v", err)
	}

	// extract new access token
	if token, ok := result["access_token"].(string); ok {
		return token, nil
	}

	return "", fmt.Errorf("no access token in refresh response")
}

// formatLogResult formats the log data into a readable string
func formatLogResult(result map[string]interface{}) string {
	// implement log formatting logic here
	// this is a placeholder - you'll need to implement the actual formatting
	return "log results formatted here"
}
