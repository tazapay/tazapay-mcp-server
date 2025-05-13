package constants

// Base URLs for different environments
const (
	// Production
	ProdBaseURL = "https://service.tazapay.com/v3"
)

// API Path Segments
const (
	CheckoutPath = "/checkout"
	FxPayoutPath = "/fx/payout"
	BalancePath  = "/balance"
)

// Production URLs
const (
	PaymentLinkBaseURLProd = ProdBaseURL + CheckoutPath
	PaymentFxBaseURLProd   = ProdBaseURL + FxPayoutPath
	BalanceBaseURLProd     = ProdBaseURL + BalancePath
)

// HTTP Method Constants
const (
	PostHTTPMethod   = "POST"
	GetHTTPMethod    = "GET"
	PutHTTPMethod    = "PUT"
	DeleteHTTPMethod = "DELETE"
)
