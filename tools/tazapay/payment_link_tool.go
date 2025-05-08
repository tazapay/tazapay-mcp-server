package tazapay

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/types"
	"github.com/tazapay/tazapay-mcp-server/utils"
)

// PaymentLinkTool defines the tool structure
type PaymentLinkTool struct{}

// NewPaymentLinkTool returns a new instance of the PaymentLinkTool
func NewPaymentLinkTool() *PaymentLinkTool {
	return &PaymentLinkTool{}
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
func (*PaymentLinkTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	// validate and extract arguments
	params, err := validateAndExtractArgs(args)
	if err != nil {
		return nil, err
	}
	// construct payload
	payload := NewPaymentLinkRequest(&params)
	// call payment link API
	resp, err := utils.HandlePOSTHttpRequest(ctx, constants.PaymentLinkBaseURLOrange, payload, constants.PostHTTPMethod)
	if err != nil {
		return nil, fmt.Errorf("HandlePOSTHttpRequest failed: %w", err)
	}
	// extract payment link
	data, ok3 := resp["data"].(map[string]any)
	if !ok3 {
		return nil, constants.ErrNoDataInResponse
	}

	paymentLink, ok4 := data["url"].(string)
	if !ok4 {
		return nil, constants.ErrMissingPaymentLink
	}
	// return result
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
func validateAndExtractArgs(args map[string]any) (types.PaymentLinkParams, error) {
	var p types.PaymentLinkParams
	var ok bool

	if p.PaymentAmount, ok = args[constants.PaymentAmountField].(float64); !ok {
		return p, utils.WrapFieldTypeError(constants.PaymentAmountField)
	}

	if p.InvoiceCurrency, ok = args[constants.InvoiceCurrencyField].(string); !ok {
		return p, utils.WrapFieldTypeError(constants.InvoiceCurrencyField)
	}

	if p.Description, ok = args[constants.TransactionDescField].(string); !ok {
		return p, utils.WrapFieldTypeError(constants.TransactionDescField)
	}

	if p.CustomerName, ok = args[constants.CustomerNameField].(string); !ok {
		return p, utils.WrapFieldTypeError(constants.CustomerNameField)
	}

	if p.CustomerEmail, ok = args[constants.CustomerEmailField].(string); !ok {
		return p, utils.WrapFieldTypeError(constants.CustomerEmailField)
	}

	if p.CustomerCountry, ok = args[constants.CustomerCountryField].(string); !ok {
		return p, utils.WrapFieldTypeError(constants.CustomerCountryField)
	}

	return p, nil
}

// NewPaymentLinkRequest constructs the API payload from the validated parameters
func NewPaymentLinkRequest(p *types.PaymentLinkParams) types.PaymentLinkRequest {
	return types.PaymentLinkRequest{
		Amount:          int64(p.PaymentAmount * constants.Num100),
		InvoiceCurrency: p.InvoiceCurrency,
		TransactionDesc: p.Description,
		CustomerDetails: map[string]string{
			"name":    p.CustomerName,
			"email":   p.CustomerEmail,
			"country": p.CustomerCountry,
		},
	}
}
