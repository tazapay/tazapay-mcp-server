package constants

// Payment Link Tool constants
const (
	PaymentLinkToolName = "tazapay_generate_payment_link_tool"
	PaymentLinkToolDesc = "Generates a checkout payment link with specified invoice details and customer information"

	InvoiceCurrencyField = "invoice_currency"
	InvoiceCurrencyDesc  = "Currency in which the invoice is to be raised (e.g., USD, EUR)"

	PaymentAmountField = "payment_amount"
	PaymentAmountDesc  = "Total invoice amount to be paid in decimal format. For example, $10.12 should be 10.12."

	CustomerNameField = "customer_name"
	CustomerNameDesc  = "Full name of the customer"

	CustomerEmailField = "customer_email"
	CustomerEmailDesc  = "Email address of the customer"

	CustomerCountryField = "customer_country"
	CustomerCountryDesc  = "Country of the customer"

	TransactionDescField = "transaction_description"
	TransactionDesc      = "Short description or purpose of the transaction"
)

// FX Tool constants
const (
	FXToolName        = "tazapay_fetch_fx_tool"
	FXToolDescription = "Get FX rate from one currency to another using Tazapay FX rate"

	FXFromField       = "from"
	FXFromDescription = "Currency to convert from. 3 letter currency code (case insensitive). Example: USD, eur, inr"

	FXToField       = "to"
	FXToDescription = "Currency to convert to. 3 letter currency code (case insensitive). Example: USD, eur, inr"

	FXAmountField       = "amount"
	FXAmountDescription = "Amount to convert in decimal format. For example, $10.12 should be 10.12."
)

// Balance Fetch tool
const (
	BalanceToolName = "tazapay_fetch_balance_tool"
	BalanceToolDesc = "Get balance from Tazapay. Send currency code to fetch balance for that currency." +
		" For all the balances available in Tazapay send empty string."

	BalanceCurrencyField = "currency"
	BalanceCurrencyDesc  = "Currency to fetch balance for. 3 letter currency code (case insensitive). Example: USD, eur, inr"
)

// Create Beneficiary Tool constants
const (
	CreateBeneficiaryToolName = "tazapay_create_beneficiary_tool"
	CreateBeneficiaryToolDesc = "Create a beneficiary on Tazapay"

	BeneficiaryNameField               = "name"
	BeneficiaryEmailField              = "email"
	BeneficiaryTypeField               = "type"
	BeneficiaryTypeDesc                = "Type of beneficiary (individual or business)"
	BeneficiaryNationalIDField         = "national_identification_number"
	BeneficiaryTaxIDField              = "tax_id"
	BeneficiaryDestinationDetailsField = "destination_details"
	BeneficiaryDestinationDetailsDesc  = "Details about the beneficiary's bank account," +
		" wallet, or local payment network"
)

// Create Payout Tool constants
const (
	CreatePayoutToolName = "tazapay_create_payout_tool"
	CreatePayoutToolDesc = "Create a payout on Tazapay"

	PayoutPurposeField            = "purpose"
	PayoutAmountField             = "amount"
	PayoutCurrencyField           = "currency"
	PayoutReferenceIDField        = "reference_id"
	PayoutBeneficiaryField        = "beneficiary"
	PayoutTransactionDescField    = "transaction_description"
	PayoutStatementDescField      = "statement_descriptor"
	PayoutChargeTypeField         = "charge_type"
	PayoutTypeField               = "type"
	PayoutTypeDesc                = "Type of payout (individual or company)"
	PayoutHoldingCurrencyField    = "holding_currency"
	PayoutOnBehalfOfField         = "on_behalf_of"
	PayoutMetadataField           = "metadata"
	PayoutBeneficiaryDetailsField = "beneficiary_details"
	PayoutBeneficiaryDetailsDesc  = "Details of the beneficiary for this payout"
)

// Get Payin Tool constants
const (
	GetPayinToolName = "tazapay_get_payin_tool"
	GetPayinToolDesc = "Fetch a payin details by ID from Tazapay"
	GetPayinIDField  = "id"
	GetPayinIDDesc   = "ID of the already created payin"
)

// Get Payout Tool constants
const (
	GetPayoutToolName = "tazapay_get_payout_tool"
	GetPayoutToolDesc = "Fetch a payout details by ID from Tazapay"
	GetPayoutIDField  = "id"
	GetPayoutIDDesc   = "ID of the existing payout"
)

// Get Beneficiary Tool constants
const (
	GetBeneficiaryToolName = "tazapay_get_beneficiary_tool"
	GetBeneficiaryToolDesc = "Fetch beneficiary data by ID from Tazapay, should start with bnf_ prefix."
	GetBeneficiaryIDField  = "id"
	GetBeneficiaryIDDesc   = "ID of the existing beneficiary"
)

// Create Payin Tool constants
const (
	CreatePayinToolName             = "tazapay_create_payin_tool"
	CreatePayinToolDesc             = "Create and confirm a payin on Tazapay"
	CreatePayinInvoiceCurrencyField = "invoice_currency"
	CreatePayinInvoiceCurrencyDesc  = "Currency in which the invoice is to be raised (e.g., USD, EUR)"
	CreatePayinAmountField          = "amount"
	CreatePayinAmountDesc           = "Payment amount in decimal format. For example, $10.12 should be 10.12."
	CreatePayinCustomerDetailsField = "customer_details"
	CreatePayinCustomerDetailsDesc  = "Customer details object (name, email, country, phone)"
)

// Cancel Payin Tool constants
const (
	CancelPayinToolName = "tazapay_cancel_payin_tool"
	CancelPayinToolDesc = "Cancel a payin on Tazapay"
	CancelPayinIDField  = "id"
	CancelPayinIDDesc   = "ID of the already created payin to cancel"
)

// Confirm Payout Tool constants
const (
	ConfirmPayoutToolName    = "tazapay_confirm_payout_tool"
	ConfirmPayoutToolDesc    = "Confirm a payout by attaching a funding source on Tazapay"
	ConfirmPayoutIDField     = "id"
	ConfirmPayoutIDDesc      = "ID of the payout to confirm"
	ConfirmPayoutSourceField = "source"
	ConfirmPayoutSourceDesc  = "Funding source ID used to fund the payout (e.g., collection ID or balance ID) to confirm"
)
