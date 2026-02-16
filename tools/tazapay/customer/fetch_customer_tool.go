package customer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	"github.com/tazapay/tazapay-mcp-server/types"
)

// FetchCustomerTool fetches a customer by ID

type FetchCustomerTool struct {
	logger *slog.Logger
}

func NewFetchCustomerTool(logger *slog.Logger) *FetchCustomerTool {
	logger.Info("Registering Fetch_Customer_Tool")
	return &FetchCustomerTool{logger: logger}
}

func (t *FetchCustomerTool) Definition() mcp.Tool {
	return mcp.NewTool(
		"tazapay_fetch_customer_tool",
		mcp.WithDescription("Fetch Customer Details by ID from Tazapay. ID must start with cus_."),
		mcp.WithString("id", mcp.Required(), mcp.Description("ID of the customer (must start with cus_).")),
	)
}

func (t *FetchCustomerTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling FetchCustomerTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args["id"].(string)
	if !ok || id == "" || utils.ValidatePrefixID("cus_", id) != nil {
		err := errors.New("missing or invalid customer id, should be starting with cus_")
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	url := fmt.Sprintf("%s/customer/%s", constants.ProdBaseURL, id)

	resp, err := utils.HandleGETHttpRequest(ctx, t.logger, url, "GET")
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to fetch customer", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in get customer API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	var customer types.Customer

	err = utils.MapToStruct(data, &customer)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to map data to Customer struct", "error", err)
		return nil, err
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			func() mcp.TextContent {
				jsonBytes, err := json.MarshalIndent(customer, "", "  ")
				if err != nil {
					return mcp.TextContent{Type: "text", Text: fmt.Sprintf("Customer: %+v", customer)}
				}

				return mcp.TextContent{Type: "text", Text: string(jsonBytes)}
			}(),
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled FetchCustomerTool request", "result", result)

	return result, nil
}
