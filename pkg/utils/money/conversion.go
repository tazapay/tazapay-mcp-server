package money

import (
	"fmt"
	"math"

	fmath "github.com/tazapay/tazapay-mcp-server/pkg/utils/math"
)

const (
	int100  = 100
	int1000 = 1000
	int10   = 10
)

// Int64ToDecimal2 - convert 2 precision int64 to float64 with 2 decimal values
func Int64ToDecimal2(x int64) float64 {
	return fmath.Round2Decimal(float64(x) / int100)
}

// Int64ToDecimal3 - convert 2 precision int64 to float64 with 3 decimal values
func Int64ToDecimal3(x int64) float64 {
	return fmath.Round3Decimal(float64(x*int10) / int1000)
}

// Int64ToDecimal0 - convert 2 precision int64 to float64 with 0 decimal values
func Int64ToDecimal0(x int64) float64 {
	return math.Round(float64(x) / int100)
}

// Decimal2ToInt64 - convert float64 with 2 decimal places to int64 (cents)
func Decimal2ToInt64(x float64) int64 {
	return int64(math.Round(x * int100))
}

// FormatCurrency - format an int64 cents value as a string with currency symbol
func FormatCurrency(amount int64, currency string) string {
	value := Int64ToDecimal2(amount)
	currencySymbol := getCurrencySymbol(currency)
	return fmt.Sprintf("%s%.2f", currencySymbol, value)
}

// getCurrencySymbol - returns the appropriate symbol for common currencies
func getCurrencySymbol(currency string) string {
	switch currency {
	case "USD":
		return "$"
	case "EUR":
		return "€"
	case "GBP":
		return "£"
	case "JPY":
		return "¥"
	case "INR":
		return "₹"
	case "BRL":
		return "R$"
	default:
		return ""
	}
}
