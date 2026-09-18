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

// ConfirmPayoutTool confirms a payout by attaching a funding source

type ConfirmPayoutTool struct {
	logger *slog.Logger
}

func NewConfirmPayoutTool(logger *slog.Logger) *ConfirmPayoutTool {
	logger.InfoContext(context.Background(), "Registering Confirm_Payout_Tool")
	return &ConfirmPayoutTool{logger: logger}
}

func (*ConfirmPayoutTool) Definition() mcp.Tool {
	return mcp.NewTool(
		constants.ConfirmPayoutToolName,
		mcp.WithDescription(constants.ConfirmPayoutToolDesc),
		mcp.WithString(constants.ConfirmPayoutIDField, mcp.Required(), mcp.Description(constants.ConfirmPayoutIDDesc)),
		mcp.WithString(constants.ConfirmPayoutSourceField, mcp.Required(), mcp.Description(constants.ConfirmPayoutSourceDesc)),
	)
}

func (t *ConfirmPayoutTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "Invalid arguments type for ConfirmPayoutTool")
		return nil, fmt.Errorf("%w", constants.ErrInvalidArgumentsType)
	}

	t.logger.InfoContext(ctx, "Handling ConfirmPayoutTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.ErrorContext(ctx, "Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args[constants.ConfirmPayoutIDField].(string)
	if !ok || id == "" || utils.ValidatePrefixID("pot_", id) != nil {
		t.logger.ErrorContext(ctx, constants.ErrMissingOrInvalidPayoutID.Error())
		return nil, constants.ErrMissingOrInvalidPayoutID
	}

	source, ok := args[constants.ConfirmPayoutSourceField].(string)
	if !ok || source == "" {
		t.logger.ErrorContext(ctx, constants.ErrMissingOrInvalidSource.Error())
		return nil, constants.ErrMissingOrInvalidSource
	}

	url := fmt.Sprintf("%s/payout/%s/confirm", constants.ProdBaseURL, id)

	payload := map[string]string{
		constants.ConfirmPayoutSourceField: source,
	}

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, url, payload, constants.PostHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to confirm payout", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in confirm payout API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	// convert amount from cents to decimal if present
	if amount, exists := data["amount"].(float64); exists {
		data["amount"] = money.Int64ToDecimal2(int64(amount))
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to marshal confirm payout data", "error", err)
		return nil, err
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: string(jsonBytes)},
		},
	}

	t.logger.InfoContext(ctx, "Successfully handled ConfirmPayoutTool request", "result", result)

	return result, nil
}
