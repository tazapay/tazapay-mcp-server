package paymentattempt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils/money"
)

// GetPaymentAttemptTool fetches a payment attempt by ID

type GetPaymentAttemptTool struct {
	logger *slog.Logger
}

func NewGetPaymentAttemptTool(logger *slog.Logger) *GetPaymentAttemptTool {
	logger.Info("Registering Get_Payment_Attempt_Tool")
	return &GetPaymentAttemptTool{logger: logger}
}

func (t *GetPaymentAttemptTool) Definition() mcp.Tool {
	return mcp.NewTool(
		"get_payment_attempt_tool",
		mcp.WithDescription("Fetch a payment attempt by ID from Tazapay"),
		mcp.WithString("id", mcp.Required(), mcp.Description("ID of the payment attempt")),
	)
}

func (t *GetPaymentAttemptTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling GetPaymentAttemptTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args["id"].(string)
	if !ok || id == "" || utils.ValidatePrefixID("pat_", id) != nil {
		err := errors.New("missing or invalid payment attempt id, should be starting with pat_")
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	url := fmt.Sprintf("%s/payment_attempt/%s", constants.ProdBaseURL, id)

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, url, nil, constants.GetHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to fetch payment attempt", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in get payment attempt API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	// Convert amount from cents to decimal value if present
	if amount, exists := data["amount"].(float64); exists {
		data["amount"] = money.Int64ToDecimal2(int64(amount))
		data["amount_original"] = amount
	}

	// Format the data as pretty JSON
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to marshal payment attempt data", "error", err)
		return nil, err
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: string(jsonBytes)},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled GetPaymentAttemptTool request", "result", result)

	return result, nil
}
