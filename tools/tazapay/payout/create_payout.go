package payout

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils/money"
	"github.com/tazapay/tazapay-mcp-server/types"
)

// CreatePayoutTool represents the create payout tool
type CreatePayoutTool struct {
	logger *slog.Logger
}

// NewCreatePayoutTool returns a new instance of the CreatePayoutTool
func NewCreatePayoutTool(logger *slog.Logger) *CreatePayoutTool {
	logger.InfoContext(context.Background(), "Registering Create_Payout_Tool")
	return &CreatePayoutTool{logger: logger}
}

func (t *CreatePayoutTool) Definition() mcp.Tool {
	t.logger.InfoContext(context.Background(), "Registering CreatePayoutTool with MCP")
	return newCreatePayoutToolSchema()
}

func newCreatePayoutToolSchema() mcp.Tool {
	return mcp.NewTool(
		"create_payout_tool",
		mcp.WithDescription("Create a payout on Tazapay"),
		mcp.WithNumber(
			"amount",
			mcp.Required(),
			mcp.Description("Amount in cents. For example, $10.12 should be 1012."),
		),
		mcp.WithString(
			"currency",
			mcp.Required(),
			mcp.Description("ISO 4217 standard. This is the payout currency."),
		),
		mcp.WithString(
			constants.KeyPurpose,
			mcp.Required(),
			mcp.Description(
				"Reason for payout. Should be in form of PYR0XX where XX is a number from 01 TO 28. "+
					"also default is PYR001.",
			),
		),
		mcp.WithString(
			constants.KeyTransactionDesc,
			mcp.Required(),
			mcp.Description("Additional Details for the payout."),
		),
		mcp.WithString(
			"reference_id",
			mcp.Description("Reference ID of the payout on your system."),
		),
		mcp.WithString(
			"statement_descriptor",
			mcp.Description("Statement Descriptor for the payout."),
		),
		mcp.WithString(
			"charge_type",
			mcp.Description("For wire transfers only."),
			mcp.Enum("shared", "ours"),
		),
		mcp.WithString(
			constants.KeyType,
			mcp.Enum("local", "swift", "wallet"),
		),
		mcp.WithString(
			"holding_currency",
			mcp.Description(
				"ISO 4217 standard, in uppercase. This is one of your balance "+
					"currencies whose balance will fund the payout.",
			),
		),
		mcp.WithString("on_behalf_of",
			mcp.Description("ID of the entity the payout is created on behalf of."),
		),
		mcp.WithString("metadata",
			mcp.Description("Set of key-value pairs to attach to the payout object."),
		),
		mcp.WithString("beneficiary",
			mcp.Description("ID of an existing payout beneficiary."),
		),
		beneficiaryDetailsSchema(),
		mcp.WithArray("documents", mcp.Items(documentSchema()), mcp.Description("Attach documents to the payout")),
		mcp.WithObject("logistics_tracking_details", mcp.Properties(logisticsTrackingDetailsSchema())),
	)
}

func beneficiaryDetailsSchema() mcp.ToolOption {
	return mcp.WithObject(constants.KeyBeneficiaryDetails,
		mcp.Properties(map[string]any{
			"name": map[string]any{
				constants.KeyType:        constants.KeyString,
				constants.KeyDescription: "Beneficiary Name",
			},
			"email": map[string]any{
				constants.KeyType:        constants.KeyString,
				constants.KeyDescription: "Beneficiary email address",
			},
			constants.KeyType: map[string]any{
				constants.KeyType:        constants.KeyString,
				"enum":                   []string{"business", "individual"},
				constants.KeyDescription: "business or individual",
			},
			"tax_id": map[string]any{
				constants.KeyType:        constants.KeyString,
				constants.KeyDescription: "Tax ID of the beneficiary.",
			},
			"national_identification_number": map[string]any{
				constants.KeyType:        constants.KeyString,
				constants.KeyDescription: "National ID of the individual.",
			},
			"phone": map[string]any{
				constants.KeyType: constants.KeyObject,
				constants.KeyProperties: map[string]any{
					"calling_code": map[string]any{
						constants.KeyType:        constants.KeyString,
						constants.KeyDescription: "Calling country code (e.g., '1' for US)",
					},
					"number": map[string]any{
						constants.KeyType:        constants.KeyString,
						constants.KeyDescription: "Phone Number",
					},
				},
			},
			"address": map[string]any{
				constants.KeyType: constants.KeyObject,
				constants.KeyProperties: map[string]any{
					"line1": map[string]any{constants.KeyType: constants.KeyString},
					"line2": map[string]any{constants.KeyType: constants.KeyString},
					"city":  map[string]any{constants.KeyType: constants.KeyString},
					"state": map[string]any{constants.KeyType: constants.KeyString},
					"country": map[string]any{
						constants.KeyType:        constants.KeyString,
						constants.KeyDescription: "Country (ISO 3166-1 alpha_2 country code)",
					},
					"postal_code": map[string]any{constants.KeyType: constants.KeyString},
				},
			},
			"destination_details": destinationDetailsSchema(),
			"document":            documentSchema(),
		}),
	)
}

func destinationDetailsSchema() map[string]any {
	return map[string]any{
		constants.KeyType: constants.KeyObject,
		constants.KeyProperties: map[string]any{
			"type": map[string]any{
				constants.KeyType:        constants.KeyString,
				"enum":                   []string{"bank", "wallet", "local_payment_network"},
				constants.KeyDescription: "bank, wallet or local_payment_network",
			},
			constants.KeyBank:       bankSchema(),
			constants.KeyWallet:     walletSchema(),
			"local_payment_network": localPaymentNetworkSchema(),
		},
	}
}

func bankSchema() map[string]any {
	return map[string]any{
		constants.KeyType: constants.KeyObject,
		constants.KeyProperties: map[string]any{
			"account_number": map[string]any{
				constants.KeyType: constants.KeyString,
				constants.KeyDescription: "Bank Account Number. Either account_number" +
					" or IBAN is mandatory",
			},
			"iban": map[string]any{
				constants.KeyType:        constants.KeyString,
				constants.KeyDescription: "IBAN. Either account_number or iban is mandatory",
			},
			"bank_name": map[string]any{
				constants.KeyType:        constants.KeyString,
				constants.KeyDescription: "Name of the bank",
			},
			"branch_name": map[string]any{
				constants.KeyType:        constants.KeyString,
				constants.KeyDescription: "Name of the branch",
			},
			constants.KeyCountry: map[string]any{
				constants.KeyType:        constants.KeyString,
				constants.KeyDescription: "Two-letter country code (ISO 3166-1 alpha-2)",
			},
			constants.KeyCurrency: map[string]any{
				constants.KeyType:        constants.KeyString,
				constants.KeyDescription: "Three-letter ISO currency code, in uppercase",
			},
			"purpose_code": map[string]any{
				constants.KeyType:        constants.KeyString,
				constants.KeyDescription: "Purpose Code for INR bank accounts",
			},
			"account_type": map[string]any{
				constants.KeyType:        constants.KeyString,
				constants.KeyDescription: "Type of bank account",
			},
			"bank_codes": bankCodesSchema(),
			"firc_required": map[string]any{
				constants.KeyType: constants.KeyBoolean,
				constants.KeyDescription: "Pass true if you want FIRC for a payout to India. " +
					"By default, it is set to false",
			},
		},
	}
}

func bankCodesSchema() map[string]any {
	description := "Various banking codes required for international transfers. " +
		"All code types such as swift_code, bic_code, ifsc_code, aba_code, sort_code, " +
		"branch_code, bsb_code, bank_code, cnaps are grouped under 'bank_codes'. " +
		"For example: {\"swift_code\":\"HSBCHCMCHKH\",\"ifsc_code\":\"HDFC0001234\"}"

	return map[string]any{
		constants.KeyType:        constants.KeyObject,
		constants.KeyDescription: description,
		constants.KeyProperties: map[string]any{
			"swift_code":  map[string]any{constants.KeyType: constants.KeyString},
			"bic_code":    map[string]any{constants.KeyType: constants.KeyString},
			"ifsc_code":   map[string]any{constants.KeyType: constants.KeyString},
			"aba_code":    map[string]any{constants.KeyType: constants.KeyString},
			"sort_code":   map[string]any{constants.KeyType: constants.KeyString},
			"branch_code": map[string]any{constants.KeyType: constants.KeyString},
			"bsb_code":    map[string]any{constants.KeyType: constants.KeyString},
			"bank_code":   map[string]any{constants.KeyType: constants.KeyString},
			"cnaps":       map[string]any{constants.KeyType: constants.KeyString},
		},
	}
}

func walletSchema() map[string]any {
	return map[string]any{
		constants.KeyType: constants.KeyObject,
		constants.KeyProperties: map[string]any{
			"deposit_address": map[string]any{constants.KeyType: constants.KeyString},
			constants.KeyType: map[string]any{constants.KeyType: constants.KeyString},
			constants.KeyCurrency: map[string]any{
				constants.KeyType: constants.KeyString,
				constants.KeyDescription: "Currency in which the payout " +
					"is to be made (in uppercase, ISO-4217 standard, " +
					"e.g., USD, EUR)",
			},
		},
	}
}

func localPaymentNetworkSchema() map[string]any {
	return map[string]any{
		constants.KeyType: constants.KeyObject,
		constants.KeyProperties: map[string]any{
			constants.KeyType:  map[string]any{constants.KeyType: constants.KeyString},
			"deposit_key_type": map[string]any{constants.KeyType: constants.KeyString},
			"deposit_key":      map[string]any{constants.KeyType: constants.KeyString},
		},
	}
}

func documentSchema() map[string]any {
	return map[string]any{
		constants.KeyType: constants.KeyObject,
		constants.KeyProperties: map[string]any{
			constants.KeyType: map[string]any{
				constants.KeyType:        constants.KeyString,
				constants.KeyDescription: "Type of document. Possible values - invoice, other",
			},
			"url": map[string]any{
				constants.KeyType:        constants.KeyString,
				constants.KeyDescription: "Dynamically downloadable URL",
			},
		},
	}
}

func logisticsTrackingDetailsSchema() map[string]any {
	return map[string]any{
		"logistics_provider_name": map[string]any{constants.KeyType: constants.KeyString},
		"logistics_provider_code": map[string]any{constants.KeyType: constants.KeyString},
		"tracking_number":         map[string]any{constants.KeyType: constants.KeyString},
	}
}

// Handle processes tool requests
func (t *CreatePayoutTool) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, argsOk := req.Params.Arguments.(map[string]any)
	if !argsOk {
		return nil, constants.ErrInvalidArgumentsType
	}

	t.logger.InfoContext(ctx, "Handling CreatePayoutTool request", "args", args)

	defer func() {
		if r := recover(); r != nil {
			t.logger.ErrorContext(ctx, "Panic recovered in Handle", "panic", r)
		}
	}()

	// Validate all arguments
	if err := t.validatePayoutArgs(ctx, args); err != nil {
		t.logger.ErrorContext(ctx, "Validation failed", constants.KeyError, err)
		return nil, err
	}

	// Convert amount from float64 to int64 cents before mapping to struct
	if amount, ok := args["amount"].(float64); ok {
		args["amount"] = money.Decimal2ToInt64(amount)
	}

	var payload types.PayoutRequest
	if err := utils.MapToStruct(args, &payload); err != nil {
		t.logger.ErrorContext(ctx, "Failed to map arguments to struct", constants.KeyError, err)
		return nil, err
	}

	return t.processPayout(ctx, &payload)
}

// processPayout handles the actual payout creation logic
func (t *CreatePayoutTool) processPayout(ctx context.Context,
	payload *types.PayoutRequest,
) (*mcp.CallToolResult, error) {
	hasBeneficiary := payload.Beneficiary != ""
	validBeneficiary := utils.ValidatePrefixID("bnf_", payload.Beneficiary) == nil

	if hasBeneficiary && validBeneficiary {
		return t.processPayoutWithBeneficiary(ctx, payload)
	}

	return t.processPayoutWithDetails(ctx, payload)
}

// processPayoutWithBeneficiary handles payout creation using existing beneficiary ID
func (t *CreatePayoutTool) processPayoutWithBeneficiary(ctx context.Context,
	payload *types.PayoutRequest,
) (*mcp.CallToolResult, error) {
	temp := *payload // create a copy

	type payoutRequestMap map[string]any
	var payloadMap payoutRequestMap

	b, err := json.Marshal(temp)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to marshal temp payload", constants.KeyError, err)
		return nil, err
	}

	unmarshalErr := json.Unmarshal(b, &payloadMap)
	if unmarshalErr != nil {
		t.logger.ErrorContext(ctx, "Failed to unmarshal temp payload", constants.KeyError, unmarshalErr)
		return nil, unmarshalErr
	}

	delete(payloadMap, constants.KeyBeneficiaryDetails)

	return t.createPayoutRequest(ctx, payloadMap)
}

// processPayoutWithDetails handles payout creation using beneficiary details
func (t *CreatePayoutTool) processPayoutWithDetails(ctx context.Context,
	payload *types.PayoutRequest,
) (*mcp.CallToolResult, error) {
	payload.Beneficiary = ""
	return t.createPayoutRequest(ctx, payload)
}

// createPayoutRequest makes the API call and handles the response
func (t *CreatePayoutTool) createPayoutRequest(ctx context.Context,
	payload any,
) (*mcp.CallToolResult, error) {
	resp, err := utils.HandlePOSTHttpRequest(ctx, t.logger, constants.CreatePayoutAPIURL,
		payload, constants.PostHTTPMethod)
	if err != nil {
		t.logger.ErrorContext(ctx, "Failed to create payout", constants.KeyError, err)
		return nil, err
	}

	data, ok := resp[constants.KeyData].(map[string]any)
	if !ok {
		t.logger.ErrorContext(ctx, "No data in create payout API response", constants.KeyData, resp)
		return nil, constants.ErrNoDataInResponse
	}

	payoutID, ok := data["id"].(string)
	if !ok || payoutID == "" {
		t.logger.ErrorContext(ctx, "No payout ID in response", constants.KeyData, data)
		return nil, constants.ErrNoBeneficiaryID
	}

	resultText := "Payout created with ID: " + payoutID

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: resultText},
		},
	}
	t.logger.InfoContext(
		ctx,
		"Successfully handled CreatePayoutTool request",
		"result",
		result,
	)

	return result, nil
}

// validatePayoutArgs validates the arguments for creating a payout
func (t *CreatePayoutTool) validatePayoutArgs(ctx context.Context,
	args map[string]any,
) error {
	//  either beneficiary or beneficiary_details, but not both or neither
	_, hasBeneficiary := args["beneficiary"]
	_, hasBeneficiaryDetails := args[constants.KeyBeneficiaryDetails]

	if (hasBeneficiary && hasBeneficiaryDetails) || (!hasBeneficiary &&
		!hasBeneficiaryDetails) {
		return constants.ErrBeneficiaryOrDetailsRequired
	}

	// In Handle, before mapping to struct, move any top-level bank codes into bank_codes
	if bdRaw, hasBD := args[constants.KeyBeneficiaryDetails]; hasBD {
		if err := validateBeneficiaryDetails(ctx, bdRaw, t.logger); err != nil {
			return err
		}
	}

	// Basic validation for required fields (amount, currency, purpose, transaction_description)
	missingFields := make([]string, 0)
	if args["amount"] == nil {
		missingFields = append(missingFields, "amount")
	}

	if args[constants.KeyCurrency] == nil || args[constants.KeyCurrency] == "" {
		missingFields = append(missingFields, constants.KeyCurrency)
	}

	if args[constants.KeyPurpose] == nil || args[constants.KeyPurpose] == "" {
		missingFields = append(missingFields, constants.KeyPurpose)
	}

	if args[constants.KeyTransactionDesc] == nil || args[constants.KeyTransactionDesc] == "" {
		missingFields = append(missingFields, constants.KeyTransactionDesc)
	}

	if len(missingFields) > 0 {
		err := utils.WrapMissingFieldsError(missingFields)
		t.logger.ErrorContext(ctx, err.Error())

		return err
	}

	return nil
}

// validateCurrencyField validates a currency field in the given map
func validateCurrencyField(ctx context.Context, data map[string]any, logger *slog.Logger) error {
	currency, ok := data[constants.KeyCurrency].(string)
	if !ok || currency == "" {
		return nil
	}

	if err := utils.ValidateCurrency(currency); err != nil {
		logger.ErrorContext(ctx, err.Error())
		return err
	}

	return nil
}

// validateCountryField validates a country field in the given map
func validateCountryField(ctx context.Context, data map[string]any, logger *slog.Logger) error {
	country, ok := data[constants.KeyCountry].(string)
	if !ok || country == "" {
		return nil
	}

	if err := utils.ValidateCountry(country); err != nil {
		logger.ErrorContext(ctx, err.Error())
		return err
	}

	return nil
}

// processBankCodes extracts and consolidates bank codes into a bank_codes field
func processBankCodes(bank map[string]any) {
	bankCodes := make(map[string]any)
	fields := []string{
		"swift_code", "bic_code", "ifsc_code", "aba_code",
		"sort_code", "branch_code", "bsb_code", "bank_code", "cnaps",
	}

	for _, field := range fields {
		if val, exists := bank[field]; exists {
			bankCodes[field] = val

			delete(bank, field)
		}
	}

	if len(bankCodes) > 0 {
		bank["bank_codes"] = bankCodes
	}
}

// validateBankDetails validates bank details within destination_details
func validateBankDetails(ctx context.Context, destinationDetails map[string]any,
	logger *slog.Logger,
) error {
	bankRaw, hasBank := destinationDetails[constants.KeyBank]
	if !hasBank {
		return nil
	}

	bank, ok := bankRaw.(map[string]any)
	if !ok {
		return nil
	}

	processBankCodes(bank)

	if err := validateCurrencyField(ctx, bank, logger); err != nil {
		return err
	}

	return validateCountryField(ctx, bank, logger)
}

// validateWalletDetails validates wallet details within destination_details
func validateWalletDetails(ctx context.Context, destinationDetails map[string]any,
	logger *slog.Logger,
) error {
	walletRaw, hasWallet := destinationDetails[constants.KeyWallet]
	if !hasWallet {
		return nil
	}

	wallet, ok := walletRaw.(map[string]any)
	if !ok {
		return nil
	}

	return validateCurrencyField(ctx, wallet, logger)
}

// validateAddressDetails validates address details within beneficiary details
func validateAddressDetails(ctx context.Context, bd map[string]any,
	logger *slog.Logger,
) error {
	addressRaw, hasAddress := bd[constants.KeyAddress]
	if !hasAddress {
		return nil
	}

	address, ok := addressRaw.(map[string]any)
	if !ok {
		return nil
	}

	return validateCountryField(ctx, address, logger)
}

// validateDestinationDetails validates the destination_details section
func validateDestinationDetails(ctx context.Context, bd map[string]any,
	logger *slog.Logger,
) error {
	ddRaw, hasDD := bd["destination_details"]
	if !hasDD {
		return nil
	}

	destinationDetails, ok := ddRaw.(map[string]any)
	if !ok {
		return nil
	}

	if err := validateBankDetails(ctx, destinationDetails, logger); err != nil {
		return err
	}

	return validateWalletDetails(ctx, destinationDetails, logger)
}

func validateBeneficiaryDetails(ctx context.Context, bdRaw any,
	logger *slog.Logger,
) error {
	bd, ok := bdRaw.(map[string]any)
	if !ok {
		return nil // or return error if you want to enforce type
	}

	if err := validateDestinationDetails(ctx, bd, logger); err != nil {
		return err
	}

	return validateAddressDetails(ctx, bd, logger)
}
