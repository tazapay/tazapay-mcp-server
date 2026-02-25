package constants

import "sync"

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

// Production URLs (kept for backward compatibility)
const (
	PaymentLinkBaseURLProd  = ProdBaseURL + CheckoutPath
	PaymentFxBaseURLProd    = ProdBaseURL + FxPayoutPath
	BalanceBaseURLProd      = ProdBaseURL + BalancePath
	CreateBeneficiaryAPIURL = ProdBaseURL + BeneficiaryPath
	CreatePayinAPIURL       = ProdBaseURL + CreatePayin
	CreatePayoutAPIURL      = ProdBaseURL + CreatePayout
)

var (
	baseURLMu   sync.RWMutex
	customBaseURL string
)

// SetBaseURL sets a custom base URL (e.g. sandbox)
func SetBaseURL(url string) {
	baseURLMu.Lock()
	defer baseURLMu.Unlock()
	customBaseURL = url
}

// GetBaseURL returns the configured base URL, falling back to ProdBaseURL
func GetBaseURL() string {
	baseURLMu.RLock()
	defer baseURLMu.RUnlock()
	if customBaseURL != "" {
		return customBaseURL
	}
	return ProdBaseURL
}

// Dynamic URL helpers
func GetPaymentLinkBaseURL() string  { return GetBaseURL() + CheckoutPath }
func GetPaymentFxBaseURL() string    { return GetBaseURL() + FxPayoutPath }
func GetBalanceBaseURL() string      { return GetBaseURL() + BalancePath }
func GetBeneficiaryAPIURL() string   { return GetBaseURL() + BeneficiaryPath }
func GetPayinAPIURL() string         { return GetBaseURL() + CreatePayin }
func GetPayoutAPIURL() string        { return GetBaseURL() + CreatePayout }

// HTTP Method Constants
const (
	PostHTTPMethod   = "POST"
	GetHTTPMethod    = "GET"
	PutHTTPMethod    = "PUT"
	DeleteHTTPMethod = "DELETE"
)
