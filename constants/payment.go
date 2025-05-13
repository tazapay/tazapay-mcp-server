package constants

// Base URLs for different environments
const (
	// Production
	ProdBaseURL = "https://service.tazapay.com/v3"

	// Orange
	OrangeBaseURL = "https://api-orange.tazapay.com/v3"

	// Purple
	PurpleBaseURL = "https://api-purple.tazapay.com/v3"

	// Brown
	BrownBaseURL = "https://api-brown.tazapay.com/v3"

	// Yellow
	YellowBaseURL = "https://api-yellow.tazapay.com/v3"
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

// Orange  URLs
const (
	PaymentLinkBaseURLOrange = OrangeBaseURL + CheckoutPath
	PaymentFxBaseURLOrange   = OrangeBaseURL + FxPayoutPath
	BalanceBaseURLOrange     = OrangeBaseURL + BalancePath
)

// Purple  URLs
const (
	PaymentLinkBaseURLPurple = PurpleBaseURL + CheckoutPath
	PaymentFxBaseURLPurple   = PurpleBaseURL + FxPayoutPath
	BalanceBaseURLPurple     = PurpleBaseURL + BalancePath
)

// Brown  URLs
const (
	PaymentLinkBaseURLBrown = BrownBaseURL + CheckoutPath
	PaymentFxBaseURLBrown   = BrownBaseURL + FxPayoutPath
	BalanceBaseURLBrown     = BrownBaseURL + BalancePath
)

// Yellow  URLs
const (
	PaymentLinkBaseURLYellow = YellowBaseURL + CheckoutPath
	PaymentFxBaseURLYellow   = YellowBaseURL + FxPayoutPath
	BalanceBaseURLYellow     = YellowBaseURL + BalancePath
)

// HTTP Method Constants
const (
	PostHTTPMethod   = "POST"
	GetHTTPMethod    = "GET"
	PutHTTPMethod    = "PUT"
	DeleteHTTPMethod = "DELETE"
)
