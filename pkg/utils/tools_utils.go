package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tazapay/tazapay-mcp-server/pkg/utils/money"
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
	if unmarshalErr := json.Unmarshal(raw, &result); unmarshalErr != nil {
		return "", fmt.Errorf("failed to parse balance response: %w", unmarshalErr)
	}

	// Ensure data is available
	if len(result.Data.Available) == 0 {
		return "No balances found.", nil
	}
	// Normalize currency if provided
	if currency != "" {
		currencyCode := strings.ToUpper(currency)
		for _, balance := range result.Data.Available {
			if strings.EqualFold(balance.Currency, currencyCode) {
				amountFloat := money.Int64ToDecimal2(balance.Amount)
				return fmt.Sprintf("%s balance: %.2f", balance.Currency, amountFloat), nil
			}
		}

		return "No balance found for currency: " + currencyCode, nil
	}
	// Format all balances
	output := "Available account balances:\n"

	for _, balance := range result.Data.Available {
		amountFloat := money.Int64ToDecimal2(balance.Amount)
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

// MoveBankCodesToNested moves bank code fields from the top level of the bank map into a nested bank_codes map.
func MoveBankCodesToNested(dest map[string]any) {
	bank, ok := dest["bank"].(map[string]any)
	if !ok {
		return
	}
	bankCodeFields := []string{
		"swift_code", "bic_code", "ifsc_code", "aba_code",
		"sort_code", "branch_code", "bsb_code", "bank_code", "cnaps",
	}
	bankCodes := make(map[string]any)

	for _, field := range bankCodeFields {
		if val, exists := bank[field]; exists {
			bankCodes[field] = val

			delete(bank, field)
		}
	}

	if len(bankCodes) > 0 {
		if existing, ok2 := bank["bank_codes"].(map[string]any); ok2 {
			for k, v := range bankCodes {
				existing[k] = v
			}
			bank["bank_codes"] = existing
		} else {
			bank["bank_codes"] = bankCodes
		}
		dest["bank"] = bank
	}
}
