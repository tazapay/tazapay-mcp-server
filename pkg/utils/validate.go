package utils

import (
	"fmt"
	"strings"

	"github.com/tazapay/tazapay-mcp-server/constants"
)

// ValidateCurrency checks if the currency is 3 uppercase letters (ISO 4217)
func ValidateCurrency(currency string) error {
	if len(currency) != constants.Num3 || currency != strings.ToUpper(currency) {
		return fmt.Errorf("%w: must be 3 uppercase letters (e.g., 'USD')", constants.ErrInvalidCurrencyFormat)
	}

	return nil
}

// ValidateCountry checks if the country is 2 uppercase letters (ISO 3166 alpha-2)
func ValidateCountry(country string) error {
	if len(country) != constants.Num2 || country != strings.ToUpper(country) {
		return fmt.Errorf("%w: must be 2 uppercase letters (e.g., 'US')", constants.ErrInvalidCountryFormat)
	}

	return nil
}

// ValidatePrefixID checks if the id starts with the given prefix
func ValidatePrefixID(prefix, id string) error {
	if !strings.HasPrefix(id, prefix) {
		return fmt.Errorf("%w: must start with '%s'", constants.ErrInvalidIDFormat, prefix)
	}

	return nil
}
