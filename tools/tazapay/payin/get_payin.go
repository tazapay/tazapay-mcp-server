package payin

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

// GetPayinTool fetches a payin by ID

type GetPayinTool struct {
	logger *slog.Logger
}

func NewGetPayinTool(logger *slog.Logger) *GetPayinTool {
	logger.Info("Registering Get_Payin_Tool")
	return &GetPayinTool{logger: logger}
}

func (t *GetPayinTool) Definition() mcp.Tool {
	return mcp.NewTool(
		constants.GetPayinToolName,
		mcp.WithDescription(constants.GetPayinToolDesc),
		mcp.WithString(constants.GetPayinIDField, mcp.Required(), mcp.Description(constants.GetPayinIDDesc)),
	)
}

func (t *GetPayinTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling GetPayinTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args["id"].(string)
	if !ok || id == "" || utils.ValidatePrefixID("pay_", id) != nil {
		err := errors.New("missing or invalid payin id, should be starting with pay_")
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	url := fmt.Sprintf("%s/payin/%s", constants.ProdBaseURL, id)

	resp, err := utils.HandleGETHttpRequest(ctx, t.logger, url, constants.GetHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to fetch payin", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in get payin API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	// Convert amount from cents to decimal value if present
	if amount, exists := data["amount"].(float64); exists {
		data["amount"] = money.Int64ToDecimal2(int64(amount))
		data["amount_original"] = amount
	}

	// Marshal the data to pretty JSON
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to marshal payin data to JSON", "error", err)
		return nil, err
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: string(jsonBytes)},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled GetPayinTool request", "result", result)

	return result, nil
}
