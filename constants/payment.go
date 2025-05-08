package constants

// Payment Tools constants based on env

// prod
const (
	PaymentLinkBaseURLProd = "https://service.tazapay.com/v3/checkout"
	PaymentFxBaseURLProd   = "https://service.tazapay.com/v3/fx/payout"
)

// orange
const (
	PaymentLinkBaseURLOrange = "https://api-orange.tazapay.com/v3/checkout"
	PaymentFxBaseURLOrange   = "https://api-orange.tazapay.com/v3/fx/payout"
)

// purple
const (
	PaymentLinkBaseURLPurple = "https://api-purple.tazapay.com/v3/checkout"
	PaymentFxBaseURLPurple   = "https://api-purple.tazapay.com/v3/fx/payout"
)

// brown
const (
	PaymentLinkBaseURLBrown = "https://api-brown.tazapay.com/v3/checkout"
	PaymentFxBaseURLBrown   = "https://api-brown.tazapay.com/v3/fx/payout"
)

// yellow
const (
	PaymentLinkBaseURLYellow = "https://api-yellow.tazapay.com/v3/checkout"
	PaymentFxBaseURLYellow   = "https://api-yellow.tazapay.com/v3/fx/payout"
)

// String Placeholder Constants
const (
	Miscellaneous  = "Miscellaneous "
	PostHTTPMethod = "POST"
	GetHTTPMethod  = "GET"
)
