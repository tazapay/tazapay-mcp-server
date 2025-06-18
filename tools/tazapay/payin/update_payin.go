package payin

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
)

// UpdatePayinTool represents the update payin tool

type UpdatePayinTool struct {
	logger *slog.Logger
}

func NewUpdatePayinTool(logger *slog.Logger) *UpdatePayinTool {
	logger.InfoContext(context.Background(), "Registering Update_Payin_Tool")
	return &UpdatePayinTool{logger: logger}
}

func (t *UpdatePayinTool) Definition() mcp.Tool {
	
	return mcp.NewTool(
		"update_payin_tool",
		mcp.WithDescription("Update a payin on Tazapay without confirming it"),
		mcp.WithString("id", mcp.Required(), mcp.Description("ID of the already created payin")),
		mcp.WithObject("customer_details",
			mcp.Properties(map[string]any{
				"customer": map[string]any{"type": "string"},
			}),
		),
		mcp.WithString("success_url", mcp.Required()),
		mcp.WithString("cancel_url", mcp.Required()),
		mcp.WithObject("shipping_details"),
		mcp.WithObject("billing_details"),
		mcp.WithObject("transaction_documents"),
		mcp.WithObject("metadata"),
		mcp.WithString("reference_id"),
		mcp.WithString("statement_descriptor"),
		mcp.WithObject("payment_method_details", mcp.Required()),
	)
}

func (t *UpdatePayinTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		err := constants.ErrInvalidType
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	t.logger.InfoContext(ctx, "Handling UpdatePayinTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.ErrorContext(ctx, "Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args["id"].(string)
	if !ok || id == "" {
		err := constants.ErrInvalidIDFormat
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	// Remove id from args for payload
	delete(args, "id") // no error to check for delete in Go, safe to ignore
	payload := args

	url := fmt.Sprintf("%s/payin/%s", constants.ProdBaseURL, id)

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, url, payload, constants.PutHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to update payin", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in update payin API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	status, ok := data["status"].(string)
	if !ok {
		t.logger.ErrorContext(ctx, "No status in update payin API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	resultText := "Payin updated. Status: " + status

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: resultText},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled UpdatePayinTool request", "result", result)

	return result, nil
}
