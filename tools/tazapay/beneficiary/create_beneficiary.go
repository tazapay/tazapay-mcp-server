package beneficiary

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	"github.com/tazapay/tazapay-mcp-server/types"
)

// CreateBeneficiaryTool represents the create beneficiary tool
type CreateBeneficiaryTool struct {
	logger *slog.Logger
}

// NewCreateBeneficiaryTool returns a new instance of the CreateBeneficiaryTool
func NewCreateBeneficiaryTool(logger *slog.Logger) *CreateBeneficiaryTool {
	logger.Info("Registering Create_Beneficiary_Tool")
	return &CreateBeneficiaryTool{logger: logger}
}

// Definition : registers this tool with the MCP
func (t *CreateBeneficiaryTool) Definition() mcp.Tool {

	return mcp.NewTool(
		constants.CreateBeneficiaryToolName,
		mcp.WithDescription("Create a new beneficiary for payouts with comprehensive destination details including bank accounts, wallets, or local payment networks"),

		// Basic beneficiary information
		mcp.WithString("name", mcp.Required(), mcp.Description("Full legal name of the beneficiary as it appears on their bank account or official documents")),
		mcp.WithString("email", mcp.Description("Email address of the beneficiary for notifications and communication")),
		mcp.WithString("type", mcp.Required(), mcp.Enum("individual", "business"), mcp.Description("Type of beneficiary - 'individual' for persons or 'business' for companies")),

		// Identification fields
		mcp.WithString("national_identification_number", mcp.Description("National ID number, passport number, or other government-issued identification")),
		mcp.WithString("tax_id", mcp.Description("Tax identification number (TIN, SSN, VAT number, etc.) depending on jurisdiction")),

		// Destination details - the core payment information
		mcp.WithObject("destination_details", mcp.Required(), mcp.Description("Payment destination configuration - specify how the beneficiary will receive funds"),
			mcp.Properties(map[string]any{
				"type": map[string]any{
					"type":        "string",
					"enum":        []string{"bank", "wallet", "local_payment_network"},
					"description": "Payment method type: 'bank' for traditional banking, 'wallet' for digital wallets, 'local_payment_network' for regional payment systems",
				},
				"bank": map[string]any{
					"type":        "object",
					"description": "Bank account details for traditional wire transfers and ACH payments",
					"properties": map[string]any{
						"account_number": map[string]any{
							"type":        "string",
							"description": "Bank account number (required if IBAN not provided)",
						},
						"iban": map[string]any{
							"type":        "string",
							"description": "International Bank Account Number (required if account_number not provided)",
						},
						"bank_name": map[string]any{
							"type":        "string",
							"description": "Full official name of the bank",
						},
						"branch_name": map[string]any{
							"type":        "string",
							"description": "Name of the specific branch (if applicable)",
						},
						"country": map[string]any{
							"type":        "string",
							"description": "Bank's country using ISO 3166-1 alpha-2 code (e.g., US, GB, SG, IN, AU)",
						},
						"currency": map[string]any{
							"type":        "string",
							"description": "Currency for receiving funds in ISO 4217 format (e.g., USD, EUR, GBP, SGD, INR)",
						},
						"purpose_code": map[string]any{
							"type":        "string",
							"description": "Purpose code required for certain jurisdictions (especially India INR transactions)",
						},
						"firc_required": map[string]any{
							"type":        "boolean",
							"description": "Set to true if Foreign Inward Remittance Certificate is required (primarily for India)",
						},
						"account_type": map[string]any{
							"type":        "string",
							"description": "Type of bank account (e.g., checking, savings, business)",
						},
						"bank_codes": map[string]any{
							"type":        "object",
							"description": "Various banking codes required for international transfers. All code types such as swift_code, bic_code, ifsc_code, aba_code, sort_code, branch_code, bsb_code, bank_code, cnaps are grouped under 'bank_codes'. For example: {\"swift_code\":\"HSBCHCMCHKH\",\"ifsc_code\":\"HDFC0001234\"}",
							"properties": map[string]any{
								"swift_code": map[string]any{
									"type":        "string",
									"description": "SWIFT/BIC code for international wire transfers (8 or 11 characters)",
								},
								"bic_code": map[string]any{
									"type":        "string",
									"description": "Bank Identifier Code (alternative to SWIFT code)",
								},
								"ifsc_code": map[string]any{
									"type":        "string",
									"description": "Indian Financial System Code (for India banks only)",
								},
								"aba_code": map[string]any{
									"type":        "string",
									"description": "ABA routing number (for US banks only, 9 digits)",
								},
								"sort_code": map[string]any{
									"type":        "string",
									"description": "Sort code (for UK banks only, 6 digits)",
								},
								"branch_code": map[string]any{
									"type":        "string",
									"description": "Branch identifier code",
								},
								"bsb_code": map[string]any{
									"type":        "string",
									"description": "Bank State Branch code (for Australian banks only)",
								},
								"bank_code": map[string]any{
									"type":        "string",
									"description": "National bank code (format varies by country)",
								},
								"cnaps": map[string]any{
									"type":        "string",
									"description": "China National Advanced Payment System code (for China banks only)",
								},
							},
						},
					},
				},
				"wallet": map[string]any{
					"type":        "object",
					"description": "Digital wallet configuration for cryptocurrency or digital payment services",
					"properties": map[string]any{
						"deposit_address": map[string]any{
							"type":        "string",
							"description": "Wallet address, account ID, or deposit identifier for the digital wallet",
						},
						"type": map[string]any{
							"type":        "string",
							"description": "Type of wallet or digital payment service (e.g., bitcoin, ethereum, paypal)",
						},
						"currency": map[string]any{
							"type":        "string",
							"description": "Digital currency or token type in ISO 4217 or standard format (e.g., BTC, ETH, USD)",
						},
					},
				},
				"local_payment_network": map[string]any{
					"type":        "object",
					"description": "Regional payment network configuration (e.g., UPI, PIX, FPS)",
					"properties": map[string]any{
						"type": map[string]any{
							"type":        "string",
							"description": "Name of the local payment network (e.g., UPI, PIX, FPS)",
						},
						"deposit_key_type": map[string]any{
							"type":        "string",
							"description": "Type of identifier used (e.g., mobile, email, account_id)",
						},
						"deposit_key": map[string]any{
							"type":        "string",
							"description": "The actual identifier value (phone number, email, or account identifier)",
						},
					},
				},
			}),
		),

		// Contact information
		mcp.WithObject("phone", mcp.Description("Beneficiary's phone number for verification and communication"),
			mcp.Properties(map[string]any{
				"number": map[string]any{
					"type":        "string",
					"description": "Phone number without country code (e.g., 1234567890)",
				},
				"calling_code": map[string]any{
					"type":        "string",
					"description": "International calling code without + sign (e.g., 1 for US, 44 for UK, 91 for India)",
				},
			}),
		),

		// Physical address
		mcp.WithObject("address", mcp.Description("Physical address of the beneficiary for compliance and verification"),
			mcp.Properties(map[string]any{
				"line1": map[string]any{
					"type":        "string",
					"description": "Primary address line (street number, street name)",
				},
				"line2": map[string]any{
					"type":        "string",
					"description": "Secondary address line (apartment, suite, unit number)",
				},
				"city": map[string]any{
					"type":        "string",
					"description": "City or town name",
				},
				"state": map[string]any{
					"type":        "string",
					"description": "State, province, or region",
				},
				"postal_code": map[string]any{
					"type":        "string",
					"description": "ZIP code, postal code, or equivalent",
				},
				"country": map[string]any{
					"type":        "string",
					"description": "Country using ISO 3166-1 alpha-2 code (e.g., US, GB, SG, IN, AU)",
				},
			}),
		),

		// Supporting documents
		mcp.WithObject("document", mcp.Description("Supporting document for identity verification or compliance"),
			mcp.Properties(map[string]any{
				"type": map[string]any{
					"type":        "string",
					"description": "Type of document (e.g., passport, drivers_license, utility_bill, bank_statement)",
				},
				"url": map[string]any{
					"type":        "string",
					"description": "Publicly accessible URL to the document file (PDF, JPG, PNG)",
				},
			}),
		),
	)
}

// Handle processes tool requests
func (t *CreateBeneficiaryTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, _ := req.Params.Arguments.(map[string]any)
	t.logger.InfoContext(ctx, "Handling CreateBeneficiaryTool request", "args", args)

	// Preprocess: Move bank code fields into bank_codes if present at top level of bank
	if dest, ok := args[constants.BeneficiaryDestinationDetailsField].(map[string]any); ok {
		utils.MoveBankCodesToNested(dest)
		args[constants.BeneficiaryDestinationDetailsField] = dest
	}

	defer func() {
		if r := recover(); r != nil {
			t.logger.Error("Panic recovered in Handle", "panic", r)
		}
	}()

	var payload types.CreateBeneficiaryRequest
	if err := utils.MapToStruct(args, &payload); err != nil {
		t.logger.ErrorContext(ctx, "Failed to map arguments to struct", "error", err)
		return nil, err
	}

	t.logger.InfoContext(ctx, "Mapped arguments to struct", "payload", payload)
	// Basic validation for required fields
	if payload.Name == "" || payload.Type == "" || payload.DestinationDetails.Type == "" {
		err := utils.WrapMissingFieldsError([]string{"name", "type", "account_id", "destination_details.type"})
		t.logger.ErrorContext(ctx, err.Error())

		return nil, err
	}

	// Validate currency and country in destination_details if present
	if dest, ok := args[constants.BeneficiaryDestinationDetailsField].(map[string]any); ok {
		if bank, ok := dest["bank"].(map[string]any); ok {
			if currency, ok := bank["currency"].(string); ok && currency != "" {
				if err := utils.ValidateCurrency(currency); err != nil {
					t.logger.ErrorContext(ctx, err.Error())
					return nil, err
				}
			}

			if country, ok := bank["country"].(string); ok && country != "" {
				if err := utils.ValidateCountry(country); err != nil {
					t.logger.ErrorContext(ctx, err.Error())
					return nil, err
				}
			}
		}

		if wallet, ok := dest["wallet"].(map[string]any); ok {
			if currency, ok := wallet["currency"].(string); ok && currency != "" {
				if err := utils.ValidateCurrency(currency); err != nil {
					t.logger.ErrorContext(ctx, err.Error())
					return nil, err
				}
			}
		}
	}
	// Validate country in address if present
	if address, ok := args["address"].(map[string]any); ok {
		if country, ok := address["country"].(string); ok && country != "" {
			if err := utils.ValidateCountry(country); err != nil {
				t.logger.ErrorContext(ctx, err.Error())
				return nil, err
			}
		}
	}

	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, constants.CreateBeneficiaryAPIURL, payload, constants.PostHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to create beneficiary", "error", err)
		return nil, err
	}

	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in create beneficiary API response", "resp", resp)
		return nil, constants.ErrNoDataInResponse
	}

	// The beneficiary ID is in data["id"]
	beneficiaryID, ok := data["id"].(string)
	if !ok || beneficiaryID == "" {
		t.logger.ErrorContext(ctx, "No beneficiary ID in response", "data", data)
		return nil, constants.ErrNoBeneficiaryID
	}

	// Optionally, you can include the destination as well
	destinationID, _ := data["destination"].(string)

	resultText := "Beneficiary created with ID: " + beneficiaryID
	if destinationID != "" {
		resultText += ", destinationID: " + destinationID
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: resultText},
		},
	}
	t.logger.InfoContext(ctx, "Successfully handled CreateBeneficiaryTool request", "result", result)

	return result, nil
}
