package tazapay

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/types"
	"github.com/tazapay/tazapay-mcp-server/utils"
)

// FXTool defines the tool structure
type FXTool struct{}

// NewFXTool returns a new instance of the FXTool
func NewFXTool() *FXTool {
	return &FXTool{}
}

// Definition registers this tool with the MCP platform
func (*FXTool) Definition() mcp.Tool {
	return mcp.NewTool(
		constants.FXToolName,
		mcp.WithDescription(constants.FXToolDescription),
		mcp.WithString(constants.FXFromField, mcp.Required(), mcp.Description(constants.FXFromDescription)),
		mcp.WithString(constants.FXToField, mcp.Required(), mcp.Description(constants.FXToDescription)),
		mcp.WithNumber(constants.FXAmountField, mcp.Required(), mcp.Description(constants.FXAmountDescription)),
	)
}

// Handle processes the tool request and returns a result
func (*FXTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments

	// validate and extract arguments
	params, err := validateAndExtractFXArgs(args)
	if err != nil {
		return nil, err
	}

	// construct URL for API call
	url := fmt.Sprintf("%s?initial_currency=%s&final_currency=%s&amount=%d",
		constants.PaymentFxBaseURLOrange, params.From, params.To, int(params.Amount))

	// call FX API
	resp, err := utils.HandleGETHttpRequest(ctx, url, constants.GetHTTPMethod)
	if err != nil {
		return nil, fmt.Errorf("HandleGETHttpRequest failed: %w", err)
	}

	// extract required fields from response
	data, ok := resp["data"].(map[string]any)
	if !ok {
		return nil, constants.ErrNoDataInResponse
	}

	exRate, ok1 := data["exchange_rate"].(float64)
	if !ok1 {
		return nil, utils.WrapFieldTypeError("exchange_rate")
	}

	converted, ok2 := data["converted_amount"].(float64)
	if !ok2 {
		return nil, utils.WrapFieldTypeError("converted_amount")
	}

	// return result
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Rate: %.2f, Converted Amount: %.2f", exRate, converted),
			},
		},
	}, nil
}

// validateAndExtractFXArgs validates request arguments and returns structured parameters
func validateAndExtractFXArgs(args map[string]any) (types.FXParams, error) {
	var p types.FXParams
	var ok bool

	if p.Amount, ok = args[constants.FXAmountField].(float64); !ok {
		return p, utils.WrapFieldTypeError(constants.FXAmountField)
	}

	if p.From, ok = args[constants.FXFromField].(string); !ok {
		return p, utils.WrapFieldTypeError(constants.FXFromField)
	}

	if p.To, ok = args[constants.FXToField].(string); !ok {
		return p, utils.WrapFieldTypeError(constants.FXToField)
	}

	return p, nil
}
