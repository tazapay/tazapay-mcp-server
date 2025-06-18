package payout

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils/money"
)

// GetPayoutTool fetches a payout by ID

type GetPayoutTool struct {
	logger *slog.Logger
}

func NewGetPayoutTool(logger *slog.Logger) *GetPayoutTool {
	logger.Info("Registering Get_Payout_Tool")
	return &GetPayoutTool{logger: logger}
}

func (*GetPayoutTool) Definition() mcp.Tool {
	return mcp.NewTool(
		constants.GetPayoutToolName,
		mcp.WithDescription(constants.GetPayoutToolDesc),
		mcp.WithString(constants.GetPayoutIDField, mcp.Required(), mcp.Description(constants.GetPayoutIDDesc)),
	)
}

func (t *GetPayoutTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "Invalid arguments type for GetPayoutTool")
		return nil, fmt.Errorf("%w", constants.ErrInvalidArgumentsType)
	}

	t.logger.InfoContext(ctx, "Handling GetPayoutTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.ErrorContext(ctx, "Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args["id"].(string)
	if !ok || id == "" || utils.ValidatePrefixID("pot_", id) != nil {
		t.logger.ErrorContext(ctx, constants.ErrMissingOrInvalidPayoutID.Error())
		return nil, constants.ErrMissingOrInvalidPayoutID
	}

	url := fmt.Sprintf("%s/payout/%s", constants.ProdBaseURL, id)

	resp, err := utils.HandleGETHttpRequest(ctx, t.logger, url, constants.GetHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to fetch payout", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in get payout API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	// Convert amount from cents to decimal value if present
	if amount, exists := data["amount"].(float64); exists {
		data["amount"] = money.Int64ToDecimal2(int64(amount))
		data["amount_original"] = amount
	}

	// Check for amount in sub-objects as well, such as transactions
	if transactions, ok := data["transactions"].([]any); ok {
		for i, trans := range transactions {
			if transMap, ok := trans.(map[string]any); ok {
				if amount, exists := transMap["amount"].(float64); exists {
					transMap["amount"] = money.Int64ToDecimal2(int64(amount))
					transMap["amount_original"] = amount
					transactions[i] = transMap
				}
			}
		}
		data["transactions"] = transactions
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to marshal payout data", "error", err)
		return nil, err
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: string(jsonBytes)},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled GetPayoutTool request", "result", result)

	return result, nil
}
