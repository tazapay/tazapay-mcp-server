package constants

import "errors"

var (
	ErrNonSuccessStatus              = errors.New("non-success status")
	ErrInvalidType                   = errors.New("invalid type for field")
	ErrNoDataInResponse              = errors.New("no data received in response")
	ErrInvalidDataFormat             = errors.New("invalid data format")
	ErrMissingPaymentLink            = errors.New("missing payment link in response")
	ErrMissingAuthKeys               = errors.New("TAZAPAY_API_KEY or TAZAPAY_API_SECRET not set. Use -e option or provide a `.tazapay-mcp-server.yaml` config file in your home directory")
	ErrNoBeneficiaryID               = errors.New("no beneficiary received id in response")
	ErrInvalidAmountFormat           = errors.New("invalid amount format for currency")
	ErrMissingRequiredFields         = errors.New("missing one of the required fields")
	ErrInvalidCurrencyFormat         = errors.New("invalid currency format")
	ErrInvalidCountryFormat          = errors.New("invalid country format")
	ErrInvalidIDFormat               = errors.New("invalid id format")
	ErrMissingOrInvalidBeneficiaryID = errors.New("missing or invalid beneficiary id")
	ErrMissingOrInvalidPayoutID      = errors.New("missing or invalid payout id, should be starting with pot_")
	ErrInvalidArgumentsType          = errors.New("invalid arguments type for GetPayoutTool")
	ErrNoStatusInFundPayoutData      = errors.New("no status in fund payout data")
	ErrBeneficiaryOrDetailsRequired  = errors.New("either 'beneficiary' or 'beneficiary_details' must be provided, but not both or neither")

	// HTTP utility specific errors
	ErrFailedToCreateHTTPRequest = errors.New("failed to create HTTP request")
	ErrHTTPRequestFailed         = errors.New("HTTP request failed")
	ErrFailedToReadResponseBody  = errors.New("failed to read response body")
	ErrFailedToDecodeResponse    = errors.New("failed to decode response JSON")
)
