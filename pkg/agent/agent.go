package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/tazapay/tazapay-mcp-server/pkg/tazapay"
)

// Agent represents an AI agent that handles natural language interactions
type Agent struct {
	client *tazapay.Client
}

// NewAgent creates a new AI agent
func NewAgent(client *tazapay.Client) *Agent {
	return &Agent{
		client: client,
	}
}

// HandleMessage processes a natural language message and returns a response
func (a *Agent) HandleMessage(ctx context.Context, message string) (string, error) {
	// Convert message to lowercase for easier matching
	msg := strings.ToLower(message)

	// Check for balance-related queries
	if strings.Contains(msg, "balance") || strings.Contains(msg, "how much") {
		currency := extractCurrency(msg)
		if currency == "" {
			return "Please specify a currency (e.g., USD, EUR) to check the balance.", nil
		}

		balance, err := a.client.GetBalance(currency)
		if err != nil {
			return "", fmt.Errorf("error checking balance: %w", err)
		}

		return fmt.Sprintf("Your %s balance is %.2f", currency, balance), nil
	}

	// Check for beneficiary-related queries
	if strings.Contains(msg, "beneficiary") || strings.Contains(msg, "add recipient") {
		// Extract beneficiary details from message
		name, email, country := extractBeneficiaryDetails(msg)
		if name == "" || email == "" || country == "" {
			return "Please provide beneficiary details including name, email, and country.", nil
		}

		beneficiaryID, err := a.client.CreateBeneficiary(name, email, country)
		if err != nil {
			return "", fmt.Errorf("error creating beneficiary: %w", err)
		}

		return fmt.Sprintf("Beneficiary created successfully! ID: %s", beneficiaryID), nil
	}

	// Check for payout-related queries
	if strings.Contains(msg, "payout") || strings.Contains(msg, "send money") {
		// Extract payout details from message
		beneficiaryID, amount, currency := extractPayoutDetails(msg)
		if beneficiaryID == "" || amount == 0 || currency == "" {
			return "Please provide payout details including beneficiary ID, amount, and currency.", nil
		}

		payoutID, err := a.client.CreatePayout(beneficiaryID, amount, currency)
		if err != nil {
			return "", fmt.Errorf("error creating payout: %w", err)
		}

		return fmt.Sprintf("Payout created successfully! ID: %s", payoutID), nil
	}

	// Check for FX rate queries
	if strings.Contains(msg, "exchange rate") || strings.Contains(msg, "convert") {
		from, to := extractCurrencies(msg)
		if from == "" || to == "" {
			return "Please specify both source and target currencies (e.g., USD to EUR).", nil
		}

		rate, err := a.client.GetFXRates(from, to)
		if err != nil {
			return "", fmt.Errorf("error getting FX rates: %w", err)
		}

		return fmt.Sprintf("Exchange rate from %s to %s: %.4f", from, to, rate), nil
	}

	// Default response for unrecognized queries
	return "I can help you with:\n" +
		"- Checking account balance\n" +
		"- Creating beneficiaries\n" +
		"- Making payouts\n" +
		"- Checking exchange rates\n\n" +
		"Please ask your question in natural language.", nil
}

// Helper functions to extract information from natural language
func extractCurrency(msg string) string {
	// Look for common currency codes
	currencies := []string{"usd", "eur", "gbp", "jpy", "aud", "cad", "chf", "cny", "hkd", "sgd"}
	for _, currency := range currencies {
		if strings.Contains(msg, currency) {
			return strings.ToUpper(currency)
		}
	}
	return ""
}

func extractBeneficiaryDetails(msg string) (name, email, country string) {
	// This is a simple implementation. In a real application, you would use
	// more sophisticated NLP techniques to extract these details.
	// For now, we'll return empty strings to indicate missing information.
	return "", "", ""
}

func extractPayoutDetails(msg string) (beneficiaryID string, amount float64, currency string) {
	// This is a simple implementation. In a real application, you would use
	// more sophisticated NLP techniques to extract these details.
	// For now, we'll return empty values to indicate missing information.
	return "", 0, ""
}

func extractCurrencies(msg string) (from, to string) {
	// Look for currency pairs in the message
	// This is a simple implementation. In a real application, you would use
	// more sophisticated NLP techniques to extract these details.
	// For now, we'll return empty strings to indicate missing information.
	return "", ""
}
