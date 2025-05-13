package tazapay

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
)

// BalanceTool represents the balance tool
type BalanceTool struct {
	logger *slog.Logger
}

// NewBalanceTool creates a new balance tool
func NewBalanceTool(logger *slog.Logger) *BalanceTool {
	return &BalanceTool{
		logger: logger,
	}
}

// Definition returns the tool definition
func (t *BalanceTool) Definition() mcp.Tool {
	return mcp.NewTool(
		constants.BalanceToolName,
		mcp.WithDescription(constants.BalanceToolDesc),
		mcp.WithString(constants.BalanceCurrencyField, mcp.Description(constants.BalanceCurrencyDesc)),
	)
}

// Handle processes tool requests
func (t *BalanceTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	currency, _ := args["currency"].(string)

	path := constants.BalancePath
	resp, err := utils.HandleGETHttpRequest(ctx, t.logger, constants.OrangeBaseURL+path, constants.GetHTTPMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	text, err := utils.GetBalances(resp, currency)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: text,
			},
		},
	}, nil
}
