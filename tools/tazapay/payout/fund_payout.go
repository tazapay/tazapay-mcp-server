package payout

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils/money"
)

// FundPayoutTool funds a payout in requires_funding state

type FundPayoutTool struct {
	logger *slog.Logger
}

func NewFundPayoutTool(logger *slog.Logger) *FundPayoutTool {
	logger.InfoContext(context.Background(), "Registering Fund_Payout_Tool")
	return &FundPayoutTool{logger: logger}
}

func (*FundPayoutTool) Definition() mcp.Tool {
	return mcp.NewTool(
		"fund_payout_tool",
		mcp.WithDescription("Fund a payout in requires_funding state by ID on Tazapay"),
		mcp.WithString("id", mcp.Required(), mcp.Description("ID of the payout in requires_funding state")),
	)
}

func (t *FundPayoutTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "Invalid arguments type for FundPayoutTool")
		return nil, fmt.Errorf("%w", constants.ErrInvalidArgumentsType)
	}

	t.logger.InfoContext(ctx, "Handling FundPayoutTool request", "args", args)

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

	url := fmt.Sprintf("%s/payout/%s/fund", constants.ProdBaseURL, id)

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, url, nil, constants.PostHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to fund payout", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in fund payout API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	status, ok := data["status"].(string)
	if !ok {
		t.logger.ErrorContext(ctx, "No status in fund payout data", "data", data)
		return nil, constants.ErrNoStatusInFundPayoutData
	}

	resultText := "Payout funded. Status: " + status

	// Convert amount from cents to decimal value if present
	if amount, exists := data["amount"].(float64); exists {
		currency, hasCurrency := data["currency"].(string)
		amountValue := money.Int64ToDecimal2(int64(amount))

		if hasCurrency {
			resultText += fmt.Sprintf("\nAmount: %s %.2f", currency, amountValue)
		} else {
			resultText += fmt.Sprintf("\nAmount: %.2f", amountValue)
		}
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: resultText},
		},
	}
	err = nil

	t.logger.InfoContext(ctx, "Successfully handled FundPayoutTool request", "result", result)

	return result, err
}
