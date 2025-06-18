package payin

import (
	"context"
	"errors"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils/money"
)

// CreatePayinTool represents the create payin tool
// Similar structure to CreateBeneficiaryTool and PaymentLinkTool

type CreatePayinTool struct {
	logger *slog.Logger
}

func NewCreatePayinTool(logger *slog.Logger) *CreatePayinTool {
	logger.Info("Registering Create_Payin_Tool")
	return &CreatePayinTool{logger: logger}
}

func (t *CreatePayinTool) Definition() mcp.Tool {
	return mcp.NewTool(
		"create_payin_tool",
		mcp.WithDescription("Create and confirm a payin on Tazapay"),
		mcp.WithString("invoice_currency", mcp.Required(), mcp.Description("Currency in which the invoice is to be raised (in uppercase, ISO-4217 standard, e.g., USD, EUR)")),
		mcp.WithNumber("amount", mcp.Required()),
		mcp.WithObject("customer_details", mcp.Required(),
			mcp.Properties(map[string]any{
				"name":    map[string]any{"type": "string"},
				"email":   map[string]any{"type": "string"},
				"country": map[string]any{"type": "string", "description": "ISO 3166 standard alpha-2 code. eg: SG, IN, US, etc."},
				"phone": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"calling_code": map[string]any{"type": "string"},
						"number":       map[string]any{"type": "string"},
					},
				},
			}),
		),
		mcp.WithString("customer"),
		mcp.WithString("success_url", mcp.Required()),
		mcp.WithString("cancel_url", mcp.Required()),
		mcp.WithString("webhook_url"),
		mcp.WithString("transaction_description", mcp.Required()),
		mcp.WithObject("shipping_details",
			mcp.Properties(map[string]any{
				"name": map[string]any{"type": "string"},
				"address": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"line1":       map[string]any{"type": "string"},
						"line2":       map[string]any{"type": "string"},
						"city":        map[string]any{"type": "string"},
						"state":       map[string]any{"type": "string"},
						"country":     map[string]any{"type": "string", "description": "ISO 3166 standard alpha-2 code. eg: SG, IN, US, etc."},
						"postal_code": map[string]any{"type": "string"},
					},
				},
				"phone": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"calling_code": map[string]any{"type": "string"},
						"number":       map[string]any{"type": "string"},
					},
				},
				"label": map[string]any{"type": "string"},
			}),
		),
		mcp.WithObject("billing_details",
			mcp.Properties(map[string]any{
				"name": map[string]any{"type": "string"},
				"address": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"line1":       map[string]any{"type": "string"},
						"line2":       map[string]any{"type": "string"},
						"city":        map[string]any{"type": "string"},
						"state":       map[string]any{"type": "string"},
						"country":     map[string]any{"type": "string", "description": "ISO 3166 standard alpha-2 code. eg: SG, IN, US, etc."},
						"postal_code": map[string]any{"type": "string"},
					},
				},
				"phone": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"calling_code": map[string]any{"type": "string"},
						"number":       map[string]any{"type": "string"},
					},
				},
				"label": map[string]any{"type": "string"},
			}),
		),
		mcp.WithObject("transaction_documents"),
		mcp.WithObject("metadata"),
		mcp.WithString("reference_id"),
		mcp.WithBoolean("confirm"),
		mcp.WithString("statement_descriptor"),
		mcp.WithObject("payment_method_details"),
		mcp.WithString("session_id"),
	)
}

func (t *CreatePayinTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling CreatePayinTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	// Validate currency
	if currency, ok := args["invoice_currency"].(string); ok && currency != "" {
		if err := utils.ValidateCurrency(currency); err != nil {
			t.logger.ErrorContext(ctx, err.Error())
			return nil, err
		}
	}
	// Validate country in customer_details if present
	if customerDetails, ok := args["customer_details"].(map[string]any); ok {
		if country, ok := customerDetails["country"].(string); ok && country != "" {
			if err := utils.ValidateCountry(country); err != nil {
				t.logger.ErrorContext(ctx, err.Error())
				return nil, err
			}
		}
	}

	// Build payload with only required and non-empty fields
	payload := make(map[string]any)

	// Required fields
	if v, ok := args["amount"]; ok {
		payload["amount"] = money.Decimal2ToInt64(v.(float64))
	}

	if v, ok := args["invoice_currency"]; ok {
		payload["invoice_currency"] = v
	}

	if v, ok := args["transaction_description"]; ok {
		payload["transaction_description"] = v
	}

	if v, ok := args["confirm"]; ok {
		payload["confirm"] = v
	}

	// Optionally add other fields only if they are non-empty
	optionalFields := []string{
		"customer_details", "customer", "success_url", "cancel_url", "webhook_url", "shipping_details", "billing_details", "transaction_documents", "metadata", "reference_id", "statement_descriptor", "payment_method_details", "session_id",
	}
	for _, field := range optionalFields {
		if v, ok := args[field]; ok {
			// Only add if not nil, not empty map, and not empty string
			switch val := v.(type) {
			case string:
				if val != "" {
					payload[field] = val
				}

			case map[string]any:
				if len(val) > 0 {
					payload[field] = val
				}

			case []any:
				if len(val) > 0 {
					payload[field] = val
				}

			case nil:
				// skip
			default:
				payload[field] = val
			}
		}
	}

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, constants.CreatePayinAPIURL, payload, constants.PostHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to create payin", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in create payin API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	payinID, ok := data["id"].(string)
	if !ok || payinID == "" {
		t.logger.ErrorContext(ctx, "No payin ID in response", "data", data)
		return nil, errors.New("no payin ID in response")
	}

	resultText := "Payin created with ID: " + payinID

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: resultText},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled CreatePayinTool request", "result", result)

	return result, nil
}
