package payin

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
)

// CancelPayinTool represents the cancel payin tool

type CancelPayinTool struct {
	logger *slog.Logger
}

func NewCancelPayinTool(logger *slog.Logger) *CancelPayinTool {
	logger.InfoContext(context.Background(), "Registering Cancel_Payin_Tool")
	return &CancelPayinTool{logger: logger}
}

func (t *CancelPayinTool) Definition() mcp.Tool {
	
	return mcp.NewTool(
		constants.CancelPayinToolName,
		mcp.WithDescription(constants.CancelPayinToolDesc),
		mcp.WithString(constants.GetPayinIDField, mcp.Required(), mcp.Description(constants.CancelPayinIDDesc)),
	)
}

func (t *CancelPayinTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		err := errors.New("invalid arguments type: expected map[string]any")
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	t.logger.InfoContext(ctx, "Handling CancelPayinTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.ErrorContext(ctx, "Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args[constants.GetPayinIDField].(string)
	if !ok || id == "" {
		err := constants.ErrInvalidIDFormat
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	url := fmt.Sprintf("%s/payin/%s/cancel", constants.ProdBaseURL, id)

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, url, nil, constants.PostHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to cancel payin", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in cancel payin API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	statusVal, ok := data["status"]
	if !ok {
		t.logger.ErrorContext(ctx, "Missing 'status' in response data", "data", data)
		return nil, errors.New("missing 'status' in response data")
	}
	resultText := "Payin cancelled. Status: " + statusVal.(string)

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: resultText},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled CancelPayinTool request", "result", result)

	return result, nil
}
