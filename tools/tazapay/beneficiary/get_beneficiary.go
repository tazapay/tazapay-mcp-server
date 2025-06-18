package beneficiary

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

// GetBeneficiaryTool fetches a beneficiary by ID

type GetBeneficiaryTool struct {
	logger *slog.Logger
}

func NewGetBeneficiaryTool(logger *slog.Logger) *GetBeneficiaryTool {
	logger.Info("Registering Get_Beneficiary_Tool")
	return &GetBeneficiaryTool{logger: logger}
}

func (t *GetBeneficiaryTool) Definition() mcp.Tool {
	return mcp.NewTool(
		constants.GetBeneficiaryToolName,
		mcp.WithDescription(constants.GetBeneficiaryToolDesc),
		mcp.WithString(constants.GetBeneficiaryIDField, mcp.Required(), mcp.Description(constants.GetBeneficiaryIDDesc)),
	)
}

func (t *GetBeneficiaryTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling GetBeneficiaryTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.ErrorContext(ctx, "Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args["id"].(string)
	if !ok || id == "" {
		err := errors.New("missing or invalid beneficiary id")
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	t.logger.Debug("Fetching beneficiary", "id", id)

	// Validate beneficiary id prefix using ValidatePrefixId
	if err := utils.ValidatePrefixID("bnf_", id); err != nil {
		t.logger.ErrorContext(ctx, err.Error())
		return nil, err
	}

	url := fmt.Sprintf("%s/beneficiary/%s", constants.ProdBaseURL, id)

	t.logger.Debug("URL", "url", url)

	resp, err := utils.HandleGETHttpRequest(ctx, t.logger, url, constants.GetHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to fetch beneficiary", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in get beneficiary API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	// Use the BeneficiaryDetails struct from types and MapToStruct utility
	var beneficiary types.Beneficiary

	err = utils.MapToStruct(data, &beneficiary)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to map data to BeneficiaryDetails struct", "error", err)
		return nil, err
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			// Marshal beneficiary struct to JSON for human-readable output
			func() mcp.TextContent {
				jsonBytes, err := json.MarshalIndent(beneficiary, "", "  ")
				if err != nil {
					return mcp.TextContent{Type: "text", Text: fmt.Sprintf("Beneficiary: %+v", beneficiary)}
				}

				return mcp.TextContent{Type: "text", Text: string(jsonBytes)}
			}(),
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled GetBeneficiaryTool request", "result", result)

	return result, nil
}
