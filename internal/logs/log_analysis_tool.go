package logs

import (
	"context"
	"fmt"

	"tazapay-mcp-server/internal/tazapay"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// HandleLogAnalysisTool returns a handler function for the log analysis tool
func HandleLogAnalysisTool(client *tazapay.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get count parameter with default value of 10
		count := 10
		if countParam, ok := request.Params.Arguments["count"]; ok {
			if countFloat, ok := countParam.(float64); ok {
				count = int(countFloat)
			}
		}

		// Get logs from Tazapay API
		logs, err := client.GetLogs(count)
		if err != nil {
			return nil, fmt.Errorf("error getting logs: %w", err)
		}

		// Format logs for response
		formattedLogs := make([]map[string]interface{}, len(logs))
		for i, log := range logs {
			if logMap, ok := log.(map[string]interface{}); ok {
				formattedLogs[i] = logMap
			}
		}

		// Return formatted response
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Log Analysis Results:\nTotal logs: %d\nLogs: %v", len(formattedLogs), formattedLogs),
				},
			},
		}, nil
	}
}

func analyzeLogs(logs []interface{}) string {
	// Count event types and statuses
	eventTypes := make(map[string]int)
	statuses := make(map[string]int)

	for _, log := range logs {
		if logMap, ok := log.(map[string]interface{}); ok {
			if eventType, ok := logMap["event_type"].(string); ok {
				eventTypes[eventType]++
			}
			if status, ok := logMap["status"].(string); ok {
				statuses[status]++
			}
		}
	}

	// Generate insights
	insights := "Log Analysis Summary:\n\n"
	insights += "Event Types:\n"
	for eventType, count := range eventTypes {
		insights += fmt.Sprintf("- %s: %d\n", eventType, count)
	}

	insights += "\nStatus Distribution:\n"
	for status, count := range statuses {
		insights += fmt.Sprintf("- %s: %d\n", status, count)
	}

	return insights
}
