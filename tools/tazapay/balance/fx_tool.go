package balance

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	fmath "github.com/tazapay/tazapay-mcp-server/pkg/utils/math"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils/money"
	"github.com/tazapay/tazapay-mcp-server/types"
)

// FXTool defines the tool structure
type FXTool struct {
	logger *slog.Logger
}

// NewFXTool returns a new instance of the FXTool
func NewFXTool(logger *slog.Logger) *FXTool {
	logger.Info("Registering FX_Tool")

	return &FXTool{
		logger: logger,
	}
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
func (t *FXTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.InfoContext(ctx, "Handling FXTool request", slog.Any("params", req.Params.Arguments))

	args := req.Params.Arguments.(map[string]any)

	// validate and extract arguments
	params, err := validateAndExtractFXArgs(t, ctx, args)
	if err != nil {
		t.logger.Error("Argument validation failed", slog.String("error", err.Error()))
		return nil, err
	}

	// construct URL for API call
	amountInt := int(fmath.Round2Decimal(params.Amount * 100))
	url := fmt.Sprintf("%s?initial_currency=%s&final_currency=%s&amount=%d",
		constants.PaymentFxBaseURLProd, params.From, params.To, amountInt)

	t.logger.InfoContext(ctx, "Calling FX API", slog.String("url", url))

	// call FX API
	resp, err := utils.HandleGETHttpRequest(ctx, t.logger, url, constants.GetHTTPMethod)
	if err != nil {
		t.logger.Error("FX API call failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("HandleGETHttpRequest failed: %w", err)
	}

	// extract required fields from response
	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.Error("No 'data' in FX API response")
		return nil, constants.ErrNoDataInResponse
	}

	exRate, ok1 := data["exchange_rate"].(float64)
	if !ok1 {
		t.logger.Error("Invalid type for exchange_rate")
		return nil, utils.WrapFieldTypeError(ctx, t.logger, "exchange_rate")
	}

	converted, ok2 := data["converted_amount"].(float64)
	if !ok2 {
		t.logger.Error("Invalid type for converted_amount")
		return nil, utils.WrapFieldTypeError(ctx, t.logger, "converted_amount")
	}

	// Use rounding function for consistent display
	formattedExRate := fmath.Round2Decimal(exRate)
	
	// If converted amount is in cents, convert to decimal
	formattedConvertedAmount := 0.0
	// If the amount looks like cents (large number), convert it to decimal
	if converted > 100 && params.Amount < 100 {
		formattedConvertedAmount = money.Int64ToDecimal2(int64(converted))
	} else {
		formattedConvertedAmount = fmath.Round2Decimal(converted)
	}
	
	// Format with currency symbols if available
	fromCurrency := params.From
	toCurrency := params.To
	
	result := fmt.Sprintf(
		"Exchange Rate: 1 %s = %.2f %s\nConverted Amount: %.2f %s = %.2f %s",
		fromCurrency, formattedExRate, toCurrency,
		params.Amount, fromCurrency, formattedConvertedAmount, toCurrency,
	)
	t.logger.InfoContext(ctx, "FXTool result ready", slog.String("result", result))

	// return result
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}

// validateAndExtractFXArgs validates request arguments and returns structured parameters
func validateAndExtractFXArgs(t *FXTool, ctx context.Context, args map[string]any) (types.FXParams, error) {
	var p types.FXParams
	var ok bool

	if p.Amount, ok = args[constants.FXAmountField].(float64); !ok {
		return p, utils.WrapFieldTypeError(ctx, t.logger, constants.FXAmountField)
	}

	if p.From, ok = args[constants.FXFromField].(string); !ok {
		return p, utils.WrapFieldTypeError(ctx, t.logger, constants.FXFromField)
	}

	if p.To, ok = args[constants.FXToField].(string); !ok {
		return p, utils.WrapFieldTypeError(ctx, t.logger, constants.FXToField)
	}

	return p, nil
}
