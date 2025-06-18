package checkout

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
)

// ExpireCheckoutTool expires a checkout session by ID

type ExpireCheckoutTool struct {
	logger *slog.Logger
}

func NewExpireCheckoutTool(logger *slog.Logger) *ExpireCheckoutTool {
	logger.Info("Registering Expire_Checkout_Tool")
	return &ExpireCheckoutTool{logger: logger}
}

func (t *ExpireCheckoutTool) Definition() mcp.Tool {
	return mcp.NewTool(
		"expire_checkout_tool",
		mcp.WithDescription("Expire a checkout session by ID on Tazapay"),
		mcp.WithString("id", mcp.Required(), mcp.Description("ID of the checkout session to expire")),
	)
}

func (t *ExpireCheckoutTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling ExpireCheckoutTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args["id"].(string)
	if !ok || id == "" {
		err := errors.New("missing or invalid checkout session id")
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	url := fmt.Sprintf("%s/checkout/%s/expire", constants.ProdBaseURL, id)

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, url, nil, constants.PostHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to expire checkout session", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in expire checkout API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	status, _ := data["status"].(string)
	resultText := "Checkout session expired. Status: " + status

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: resultText},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled ExpireCheckoutTool request", "result", result)

	return result, nil
}
