package beneficiary

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
)

// UpdateBeneficiaryTool updates an existing beneficiary by ID

type UpdateBeneficiaryTool struct {
	logger *slog.Logger
}

func NewUpdateBeneficiaryTool(logger *slog.Logger) *UpdateBeneficiaryTool {
	logger.InfoContext(context.Background(), "Registering Update_Beneficiary_Tool")
	return &UpdateBeneficiaryTool{logger: logger}
}

func (t *UpdateBeneficiaryTool) Definition() mcp.Tool {

	return mcp.NewTool(
		"update_beneficiary_tool",
		mcp.WithDescription("Update an existing beneficiary by ID on Tazapay"),
		mcp.WithString("id", mcp.Required(), mcp.Description("ID of the existing beneficiary")),
		mcp.WithString("name"),
		mcp.WithString("type"),
		mcp.WithObject("address",
			mcp.Properties(map[string]any{
				"line1":       map[string]any{"type": "string"},
				"line2":       map[string]any{"type": "string"},
				"city":        map[string]any{"type": "string"},
				"state":       map[string]any{"type": "string"},
				"country":     map[string]any{"type": "string"},
				"postal_code": map[string]any{"type": "string"},
			}),
		),
		mcp.WithObject("phone",
			mcp.Properties(map[string]any{
				"calling_code": map[string]any{"type": "string"},
				"number":       map[string]any{"type": "string"},
			}),
		),
		mcp.WithObject("destination_details",
			mcp.Properties(map[string]any{
				"type": map[string]any{"type": "string"},
				"bank": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"account_number": map[string]any{"type": "string"},
						"iban":           map[string]any{"type": "string"},
						"bank_name":      map[string]any{"type": "string"},
						"branch_name":    map[string]any{"type": "string"},
						"country":        map[string]any{"type": "string"},
						"currency":       map[string]any{"type": "string"},
						"purpose_code":   map[string]any{"type": "string"},
						"bank_codes": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"swift_code":  map[string]any{"type": "string"},
								"bic_code":    map[string]any{"type": "string"},
								"ifsc_code":   map[string]any{"type": "string"},
								"aba_code":    map[string]any{"type": "string"},
								"sort_code":   map[string]any{"type": "string"},
								"branch_code": map[string]any{"type": "string"},
								"bsb_code":    map[string]any{"type": "string"},
								"bank_code":   map[string]any{"type": "string"},
								"cnaps":       map[string]any{"type": "string"},
							},
						},
						"firc_required": map[string]any{"type": "boolean"},
						"account_type":  map[string]any{"type": "string"},
					},
				},
				"wallet": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"deposit_address": map[string]any{"type": "string"},
						"type":            map[string]any{"type": "string"},
						"currency":        map[string]any{"type": "string"},
					},
				},
				"local_payment_network": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"type":             map[string]any{"type": "string"},
						"deposit_key_type": map[string]any{"type": "string"},
						"deposit_key":      map[string]any{"type": "string"},
					},
				},
			}),
		),
		mcp.WithObject("metadata"),
		mcp.WithString("tax_id"),
		mcp.WithObject("documents",
			mcp.Properties(map[string]any{
				"type": map[string]any{"type": "string"},
				"url":  map[string]any{"type": "string"},
			}),
		),
		mcp.WithString("national_identification_number"),
	)
}

func (t *UpdateBeneficiaryTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling UpdateBeneficiaryTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.ErrorContext(ctx, "Panic recovered in Handle", "panic", r)
		}
	}()

	id, ok := args["id"].(string)
	if !ok || id == "" {
		err := constants.ErrMissingOrInvalidBeneficiaryID
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}
	// Validate beneficiary id prefix using ValidatePrefixId
	if err := utils.ValidatePrefixID("bnf_", id); err != nil {
		t.logger.ErrorContext(ctx, err.Error())
		return nil, err
	}

	// Remove id from args for payload
	delete(args, "id")
	payload := args

	url := fmt.Sprintf("%s/beneficiary/%s", constants.ProdBaseURL, id)

	resp, err := utils.HandlePUTHttpRequest(ctx, t.logger, url, payload, constants.PutHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to update beneficiary", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in update beneficiary API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: fmt.Sprintf("Beneficiary updated: %+v", data)},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled UpdateBeneficiaryTool request", "result", result)

	return result, nil
}
