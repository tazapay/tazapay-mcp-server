package balance

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

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
	logger.InfoContext(context.Background(), "Registering Balance_Tool")

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
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, errors.New("invalid arguments type")
	}

	currency, ok := args["currency"].(string)
	if !ok {
		return nil, errors.New("currency parameter missing or not a string")
	}

	// If empty string, fetch all balances
	if len(currency) == 0 {
		currency = ""
	} else if len(currency) == 3 {
		// Convert to uppercase if needed
		currency = strings.ToUpper(currency)
	} else {
		return nil, errors.New("currency must be 3 letters (e.g., USD, INR) or empty to fetch all balances")
	}

	t.logger.Info("handling balance tool request", slog.Any("args", args))

	resp, err := utils.HandleGETHttpRequest(ctx, t.logger, constants.BalanceBaseURLProd, constants.GetHTTPMethod)
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
