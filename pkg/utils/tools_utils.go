package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"github.com/tazapay/tazapay-mcp-server/types"
)

// GetBalances parses balance data and returns specific or all available balances.
// - If a currency is passed, it returns balance for that currency.
// - If no currency is passed, it returns all available balances.
func GetBalances(data map[string]any, currency string) (string, error) {
	// Marshal map to JSON bytes
	raw, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal balance data: %w", err)
	}

	var result types.BalanceResponse
	// Unmarshal into the BalanceResponse struct
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("failed to parse balance response: %w", err)
	}

	// Ensure data is available
	if len(result.Data.Available) == 0 {
		return "No balances found.", nil
	}

	// Normalize currency if provided
	if len(currency) > 0 && currency != "" {
		currencyCode := strings.ToUpper(currency)
		for _, balance := range result.Data.Available {
			if strings.EqualFold(balance.Currency, currencyCode) {
				amountInt, err := strconv.Atoi(balance.Amount)
				if err != nil {
					return "", fmt.Errorf("invalid amount format for %s: %w", balance.Currency, err)
				}
				amountFloat := float64(amountInt) / 100.0
				return fmt.Sprintf("%s balance: %.2f", balance.Currency, amountFloat), nil
			}
		}
		return fmt.Sprintf("No balance found for currency: %s", currencyCode), nil
	}

	// Format all balances
	output := "Available account balances:\n"
	for _, balance := range result.Data.Available {
		amountInt, err := strconv.Atoi(balance.Amount)
		if err != nil {
			return "", fmt.Errorf("invalid amount format for %s: %w", balance.Currency, err)
		}
		amountFloat := float64(amountInt) / 100.0
		output += fmt.Sprintf("- %s: %.2f\n", balance.Currency, amountFloat)
	}
	return output, nil
}



// MapToStruct converts map[string]any to any struct using JSON marshaling.
// Pass a pointer to the output struct as `out`.
func MapToStruct(input map[string]any, out any) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	if ok := json.Unmarshal(jsonData, out); ok != nil {
		return fmt.Errorf("failed to unmarshal JSON to struct: %w", ok)
	}

	return nil
}
