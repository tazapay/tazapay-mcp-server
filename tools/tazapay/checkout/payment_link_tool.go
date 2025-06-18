package checkout

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils/money"
	"github.com/tazapay/tazapay-mcp-server/types"
)

// PaymentLinkTool defines the tool structure
type PaymentLinkTool struct {
	logger *slog.Logger
}

// NewPaymentLinkTool returns a new instance of the PaymentLinkTool
func NewPaymentLinkTool(logger *slog.Logger) *PaymentLinkTool {
	logger.InfoContext(context.Background(), "Registering Payment_Link_Tool")

	return &PaymentLinkTool{
		logger: logger,
	}
}

// Definition registers this tool with the MCP platform
func (*PaymentLinkTool) Definition() mcp.Tool {
	return mcp.NewTool(
		constants.PaymentLinkToolName,
		mcp.WithDescription(constants.PaymentLinkToolDesc),
		mcp.WithString(constants.InvoiceCurrencyField, mcp.Required(), mcp.Description("Currency in which the invoice is to be raised (in uppercase, ISO-4217 standard, e.g., USD, EUR)")),
		mcp.WithNumber(constants.PaymentAmountField, mcp.Required(), mcp.Description(constants.PaymentAmountDesc)),
		mcp.WithString(constants.CustomerNameField, mcp.Required(), mcp.Description(constants.CustomerNameDesc)),
		mcp.WithString(constants.CustomerEmailField, mcp.Required(),
			mcp.Description(constants.CustomerEmailDesc),
		),
		mcp.WithString(
			constants.CustomerCountryField,
			mcp.Required(),
			mcp.Description("Country of the customer (ISO 3166 standard alpha-2 code. eg: SG, IN, US, etc.)"),
		),
		mcp.WithString(constants.TransactionDescField, mcp.Required(), mcp.Description(constants.TransactionDesc)),
	)
}

// Handle processes the tool request and returns a result
func (t *PaymentLinkTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments.(map[string]any)

	t.logger.InfoContext(ctx, "handling payment link tool request", slog.Any("args", args))

	params, err := validateAndExtractArgs(ctx, t, args)
	if err != nil {
		t.logger.ErrorContext(ctx, "argument validation failed", slog.String("error", err.Error()))
		return nil, err
	}

	payload := NewPaymentLinkRequest(&params)
	t.logger.InfoContext(ctx, "constructed payment link payload", slog.Any("payload", payload))

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, constants.PaymentLinkBaseURLProd,
		payload, constants.PostHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "payment link API call failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("HandlePOSTHttpRequest failed: %w", err)
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "no data found in payment link API response", slog.Any("response", resp))
		return nil, constants.ErrNoDataInResponse
	}

	paymentLink, ok := data["url"].(string)
	if !ok {
		t.logger.ErrorContext(ctx, "payment link missing in API response", slog.Any("data", data))
		return nil, constants.ErrMissingPaymentLink
	}

	paymentID, ok := data["id"].(string)
	if !ok {
		t.logger.ErrorContext(ctx, "checkout id missing in API response", slog.Any("data", data))
		return nil, constants.ErrNoBeneficiaryID // Use static error
	}

	// Marshal the full data to JSON for output
	fullDataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.logger.ErrorContext(ctx, "failed to marshal full data for output", slog.String("error", err.Error()))
		fullDataJSON = []byte("<failed to marshal data>")
	}

	t.logger.InfoContext(ctx, "payment link successfully generated",
		slog.String("url", paymentLink),
		slog.String("id", paymentID),
	)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf(
					"Payment Link URL: %s\nPayment Link ID: %s\nFull Data: %s",
					paymentLink,
					paymentID,
					string(fullDataJSON),
				),
			},
		},
	}, nil
}

// validateAndExtractArgs validates request arguments and returns structured parameters
func validateAndExtractArgs(ctx context.Context, t *PaymentLinkTool, args map[string]any) (types.PaymentLinkParams, error) {
	var p types.PaymentLinkParams
	var ok bool

	if p.PaymentAmount, ok = args[constants.PaymentAmountField].(float64); !ok {
		return p, utils.WrapFieldTypeError(ctx, t.logger, constants.PaymentAmountField)
	}

	if p.InvoiceCurrency, ok = args[constants.InvoiceCurrencyField].(string); !ok {
		return p, utils.WrapFieldTypeError(ctx, t.logger, constants.InvoiceCurrencyField)
	}

	if p.Description, ok = args[constants.TransactionDescField].(string); !ok {
		return p, utils.WrapFieldTypeError(ctx, t.logger, constants.TransactionDescField)
	}

	if p.CustomerName, ok = args[constants.CustomerNameField].(string); !ok {
		return p, utils.WrapFieldTypeError(ctx, t.logger, constants.CustomerNameField)
	}

	if p.CustomerEmail, ok = args[constants.CustomerEmailField].(string); !ok {
		return p, utils.WrapFieldTypeError(ctx, t.logger, constants.CustomerEmailField)
	}

	if p.CustomerCountry, ok = args[constants.CustomerCountryField].(string); !ok {
		return p, utils.WrapFieldTypeError(ctx, t.logger, constants.CustomerCountryField)
	}

	if err := utils.ValidateCurrency(p.InvoiceCurrency); err != nil {
		return p, err
	}

	if err := utils.ValidateCountry(p.CustomerCountry); err != nil {
		return p, err
	}

	return p, nil
}

// NewPaymentLinkRequest constructs the API payload from the validated parameters
func NewPaymentLinkRequest(p *types.PaymentLinkParams) types.PaymentLinkRequest {
	return types.PaymentLinkRequest{
		Amount:                 money.Decimal2ToInt64(p.PaymentAmount),
		InvoiceCurrency:        p.InvoiceCurrency,
		TransactionDescription: p.Description,
		CustomerDetails: map[string]string{
			"name":    p.CustomerName,
			"email":   p.CustomerEmail,
			"country": p.CustomerCountry,
		},
	}
}
