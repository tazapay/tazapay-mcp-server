package constants

import "errors"

var (
	ErrNonSuccessStatus   = errors.New("non-success status")
	ErrInvalidType        = errors.New("invalid type for field")
	ErrNoDataInResponse   = errors.New("no data in response")
	ErrInvalidDataFormat  = errors.New("invalid data format")
	ErrMissingPaymentLink = errors.New("missing payment link in response")
	ErrMissingAuthKeys    = errors.New(
		"TAZAPAY_API_KEY or TAZAPAY_API_SECRET not set. Use -e option or provide a " +
			"`.tazapay-mcp-server.yaml` config file in your home directory",
	)
)
