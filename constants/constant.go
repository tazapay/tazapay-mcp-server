package constants

const (
	Num100       = 100
	Error        = "error"
	OpenFileMode = 0o666
	Num64        = 64
	Num2         = 2
	Num3         = 3

	// Added for http_utils.go magic string/number linting
	StrTAZAPAYAuthToken           = "TAZAPAY_AUTH_TOKEN"
	StrFailedToCreateHTTPRequest  = "Failed to create HTTP request"
	StrErrorCreatingRequest       = "error creating request: %w"
	StrHTTPRequestFailed          = "HTTP request failed"
	StrErrorMakingRequest         = "error making request: %w"
	StrFailedToReadResponseBody   = "Failed to read response body"
	StrErrorReadingResponseBody   = "error reading response body: %w"
	StrNonSuccessHTTPResponse     = "Non-success HTTP response"
	StrStatusCode                 = "status_code"
	StrBody                       = "body"
	StrWrappedErrorWithBody       = "%w: %v, body: %s"
	StrFailedToDecodeResponseJSON = "Failed to decode response JSON"
	StrErrorDecodingResponse      = "error decoding response: %w"

	// Common string keys for payloads, logs, and validation
	KeyError              = "error"
	KeyData               = "data"
	KeyBeneficiaryDetails = "beneficiary_details"
	KeyCurrency           = "currency"
	KeyPurpose            = "purpose"
	KeyTransactionDesc    = "transaction_description"
	KeyAddress            = "address"
	KeyCountry            = "country"
	KeyBank               = "bank"
	KeyWallet             = "wallet"

	// string constants for repeated schema keys
	KeyObject      = "object"
	KeyProperties  = "properties"
	KeyType        = "type"
	KeyString      = "string"
	KeyBoolean     = "boolean"
	KeyDescription = "description"

	// string constants for transport types
	TransportTypeStdio          = "stdio"
	TransportTypeStreamableHTTP = "streamablehttp"
)
