package customer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	"github.com/tazapay/tazapay-mcp-server/types"
)

// CreateCustomerTool creates a customer in Tazapay

type CreateCustomerTool struct {
	logger *slog.Logger
}

func NewCreateCustomerTool(logger *slog.Logger) *CreateCustomerTool {
	logger.Info("Registering Create_Customer_Tool")
	return &CreateCustomerTool{logger: logger}
}

func (t *CreateCustomerTool) Definition() mcp.Tool {
	return mcp.NewTool(
		"tazapay_create_customer_tool",
		mcp.WithDescription("Create a Customer in Tazapay."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Customer's name")),
		mcp.WithString("email", mcp.Required(), mcp.Description("Customer's email address")),
		mcp.WithString("country", mcp.Required(), mcp.Description("Customer's country. ISO 3166 standard alpha-2 code.")),
		mcp.WithString("reference_id", mcp.Description("The unique reference_id on your system representing the customer")),
		mcp.WithObject("phone", mcp.Description("Customer's phone details")),
		mcp.WithArray("billing", mcp.Description("Customer's billing details"), mcp.Items(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{"type": "string", "description": "Name"},
				"address": map[string]any{
					"type":        "object",
					"description": "Address",
					"properties": map[string]any{
						"phone": map[string]any{"type": "object", "description": "Phone"},
						"label": map[string]any{"type": "string", "description": "Denotes the type of address (Example - home, work)"},
					},
				},
			},
		})),
		mcp.WithArray("shipping", mcp.Description("Customer's shipping details"), mcp.Items(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{"type": "string", "description": "Name"},
				"address": map[string]any{
					"type":        "object",
					"description": "Address",
					"properties": map[string]any{
						"phone": map[string]any{"type": "object", "description": "Phone"},
						"label": map[string]any{"type": "string", "description": "Denotes the type of address (Example - home, work)"},
					},
				},
			},
		})),
		mcp.WithObject("metadata", mcp.Description("Set of key-value pairs to attach to the customer object")),
	)
}

func (t *CreateCustomerTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling CreateCustomerTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	url := constants.ProdBaseURL + "/customer"

	resultMap, err := utils.HandlePOSTHttpRequest(ctx, t.logger, url, args, constants.PostHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "HTTP request failed", "error", err)
		return nil, err
	}

	data, ok := resultMap["data"]
	if !ok {
		t.logger.ErrorContext(ctx, "No 'data' field in response", "response", resultMap)
		return nil, constants.ErrNoDataInResponse
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to marshal 'data' field", "error", err)
		return nil, err
	}

	var customer types.Customer
	if err := json.Unmarshal(jsonData, &customer); err != nil {
		t.logger.ErrorContext(ctx, "Failed to unmarshal 'data' field", "error", err)
		return nil, err
	}

	resultText := "Customer created with ID: " + customer.ID + ", name: " + customer.Name

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: resultText},
		},
	}, nil
}
