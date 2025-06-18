package constants

// Base URLs for different environments
const (
	// Production
	ProdBaseURL = "https://service.tazapay.com/v3"
)

// API Path Segments
const (
	CheckoutPath    = "/checkout"
	FxPayoutPath    = "/fx/payout"
	BalancePath     = "/balance"
	BeneficiaryPath = "/beneficiary"
	CreatePayin     = "/payin"
	CreatePayout    = "/payout"
)

// Production URLs
const (
	PaymentLinkBaseURLProd  = ProdBaseURL + CheckoutPath
	PaymentFxBaseURLProd    = ProdBaseURL + FxPayoutPath
	BalanceBaseURLProd      = ProdBaseURL + BalancePath
	CreateBeneficiaryAPIURL = ProdBaseURL + BeneficiaryPath
	CreatePayinAPIURL       = ProdBaseURL + CreatePayin
	CreatePayoutAPIURL      = ProdBaseURL + CreatePayout
)

// HTTP Method Constants
const (
	PostHTTPMethod   = "POST"
	GetHTTPMethod    = "GET"
	PutHTTPMethod    = "PUT"
	DeleteHTTPMethod = "DELETE"
)
