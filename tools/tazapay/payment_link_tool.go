package tazapay

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	"github.com/tazapay/tazapay-mcp-server/types"
)

// PaymentLinkTool defines the tool structure
type PaymentLinkTool struct {
	logger *slog.Logger
}

// NewPaymentLinkTool returns a new instance of the PaymentLinkTool
func NewPaymentLinkTool(logger *slog.Logger) *PaymentLinkTool {
	logger.Info("Initializing PaymentLinkTool")

	return &PaymentLinkTool{
		logger: logger,
	}
}

// Definition registers this tool with the MCP platform
func (*PaymentLinkTool) Definition() mcp.Tool {
	return mcp.NewTool(
		constants.PaymentLinkToolName,
		mcp.WithDescription(constants.PaymentLinkToolDesc),
		mcp.WithString(constants.InvoiceCurrencyField, mcp.Required(), mcp.Description(constants.InvoiceCurrencyDesc)),
		mcp.WithNumber(constants.PaymentAmountField, mcp.Required(), mcp.Description(constants.PaymentAmountDesc)),
		mcp.WithString(constants.CustomerNameField, mcp.Required(), mcp.Description(constants.CustomerNameDesc)),
		mcp.WithString(constants.CustomerEmailField, mcp.Required(), mcp.Description(constants.CustomerEmailDesc)),
		mcp.WithString(constants.CustomerCountryField, mcp.Required(), mcp.Description(constants.CustomerCountryDesc)),
		mcp.WithString(constants.TransactionDescField, mcp.Required(), mcp.Description(constants.TransactionDesc)),
	)
}

// Handle processes the tool request and returns a result
func (t *PaymentLinkTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments

	t.logger.Info("handling payment link tool request", slog.Any("args", args))

	params, err := validateAndExtractArgs(t, args)
	if err != nil {
		t.logger.Error("argument validation failed", slog.String("error", err.Error()))
		return nil, err
	}

	payload := NewPaymentLinkRequest(&params)
	t.logger.Info("constructed payment link payload", slog.Any("payload", payload))

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, constants.PaymentLinkBaseURLProd,
		payload, constants.PostHTTPMethod)
	if err != nil {
		t.logger.Error("payment link API call failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("HandlePOSTHttpRequest failed: %w", err)
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.Error("no data found in payment link API response", slog.Any("response", resp))
		return nil, constants.ErrNoDataInResponse
	}

	paymentLink, ok := data["url"].(string)
	if !ok {
		t.logger.Error("payment link missing in API response", slog.Any("data", data))
		return nil, constants.ErrMissingPaymentLink
	}

	t.logger.Info("payment link successfully generated", slog.String("url", paymentLink))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "Payment Link URL: " + paymentLink,
			},
		},
	}, nil
}

// validateAndExtractArgs validates request arguments and returns structured parameters
func validateAndExtractArgs(t *PaymentLinkTool, args map[string]any) (types.PaymentLinkParams, error) {
	var p types.PaymentLinkParams
	var ok bool

	if p.PaymentAmount, ok = args[constants.PaymentAmountField].(float64); !ok {
		return p, utils.WrapFieldTypeError(t.logger, constants.PaymentAmountField)
	}

	if p.InvoiceCurrency, ok = args[constants.InvoiceCurrencyField].(string); !ok {
		return p, utils.WrapFieldTypeError(t.logger, constants.InvoiceCurrencyField)
	}

	if p.Description, ok = args[constants.TransactionDescField].(string); !ok {
		return p, utils.WrapFieldTypeError(t.logger, constants.TransactionDescField)
	}

	if p.CustomerName, ok = args[constants.CustomerNameField].(string); !ok {
		return p, utils.WrapFieldTypeError(t.logger, constants.CustomerNameField)
	}

	if p.CustomerEmail, ok = args[constants.CustomerEmailField].(string); !ok {
		return p, utils.WrapFieldTypeError(t.logger, constants.CustomerEmailField)
	}

	if p.CustomerCountry, ok = args[constants.CustomerCountryField].(string); !ok {
		return p, utils.WrapFieldTypeError(t.logger, constants.CustomerCountryField)
	}

	return p, nil
}

// NewPaymentLinkRequest constructs the API payload from the validated parameters
func NewPaymentLinkRequest(p *types.PaymentLinkParams) types.PaymentLinkRequest {
	return types.PaymentLinkRequest{
		Amount:                 int64(p.PaymentAmount * constants.Num100),
		InvoiceCurrency:        p.InvoiceCurrency,
		TransactionDescription: p.Description,
		CustomerDetails: map[string]string{
			"name":    p.CustomerName,
			"email":   p.CustomerEmail,
			"country": p.CustomerCountry,
		},
	}
}
