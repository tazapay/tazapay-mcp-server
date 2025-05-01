package nlp

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"tazapay-mcp-server/pkg/tazapay"
)

// Intent represents the user's intention
type Intent struct {
	Name       string
	Confidence float64
	Entities   map[string]string
}

// Agent represents an NLP agent that can understand and process natural language queries
type Agent struct {
	client  *tazapay.Client
	context map[string]string
}

// NewAgent creates a new NLP agent
func NewAgent(client *tazapay.Client) *Agent {
	return &Agent{
		client:  client,
		context: make(map[string]string),
	}
}

// Process handles a natural language query and returns a response
func (a *Agent) Process(query string) (string, error) {
	// Normalize query
	query = strings.ToLower(strings.TrimSpace(query))

	// Detect intent
	intent := a.detectIntent(query)

	// Extract entities
	intent.Entities = a.extractEntities(query)

	// Handle missing required entities
	missingEntities := a.getMissingEntities(intent)
	if len(missingEntities) > 0 {
		return fmt.Sprintf("Please provide the following information: %s", strings.Join(missingEntities, ", ")), nil
	}

	// Process the intent
	return a.handleIntent(intent)
}

// detectIntent identifies the user's intention from the query
func (a *Agent) detectIntent(query string) Intent {
	patterns := map[string]string{
		"check_balance": `(?:check|what(?:'s| is)|how much|show me).*(?:balance|money|funds)`,
		"get_fx_rate":   `(?:exchange rate|convert|fx rate|rate).*(?:from|to|USD|EUR|SGD|VND|INR)`,
		"create_payout": `(?:send|transfer|make|create|pay).*(?:\d+).*(?:USD|EUR|SGD|VND|INR)`,
		"help":          `(?:help|what can you do|commands|options)`,
	}

	maxConfidence := 0.0
	bestIntent := "unknown"

	for intent, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(query, -1)
		if len(matches) > 0 {
			// Calculate confidence based on match length and position
			matchLen := 0
			for _, match := range matches {
				matchLen += len(match)
			}
			confidence := float64(matchLen) / float64(len(query))

			if confidence > maxConfidence {
				maxConfidence = confidence
				bestIntent = intent
			}
		}
	}

	return Intent{
		Name:       bestIntent,
		Confidence: maxConfidence,
		Entities:   make(map[string]string),
	}
}

// extractEntities extracts relevant information from the query
func (a *Agent) extractEntities(query string) map[string]string {
	entities := make(map[string]string)

	// Extract currencies
	currencyPattern := regexp.MustCompile(`(?i)(USD|EUR|SGD|VND|INR)`)
	currencies := currencyPattern.FindAllString(query, -1)
	if len(currencies) > 0 {
		entities["from_currency"] = strings.ToUpper(currencies[0])
		if len(currencies) > 1 {
			entities["to_currency"] = strings.ToUpper(currencies[1])
		}
	}

	// Extract amounts
	amountPattern := regexp.MustCompile(`(\d+(?:\.\d+)?)`)
	amounts := amountPattern.FindAllString(query, -1)
	if len(amounts) > 0 {
		entities["amount"] = amounts[0]
	}

	// Extract beneficiary information
	beneficiaryPattern := regexp.MustCompile(`(?:to|for)\s+(\w+(?:\s+\w+)*)`)
	if matches := beneficiaryPattern.FindStringSubmatch(query); len(matches) > 1 {
		entities["beneficiary_name"] = matches[1]
	}

	return entities
}

// getMissingEntities checks which required entities are missing
func (a *Agent) getMissingEntities(intent Intent) []string {
	var missing []string

	switch intent.Name {
	case "get_fx_rate":
		if _, ok := intent.Entities["from_currency"]; !ok {
			missing = append(missing, "source currency")
		}
		if _, ok := intent.Entities["to_currency"]; !ok {
			missing = append(missing, "target currency")
		}
		if _, ok := intent.Entities["amount"]; !ok {
			missing = append(missing, "amount")
		}

	case "create_payout":
		if _, ok := intent.Entities["amount"]; !ok {
			missing = append(missing, "amount")
		}
		if _, ok := intent.Entities["currency"]; !ok {
			missing = append(missing, "currency")
		}
		if _, ok := intent.Entities["beneficiary_id"]; !ok {
			missing = append(missing, "beneficiary ID")
		}
	}

	return missing
}

// handleIntent processes the detected intent
func (a *Agent) handleIntent(intent Intent) (string, error) {
	switch intent.Name {
	case "check_balance":
		return a.handleCheckBalance()

	case "get_fx_rate":
		return a.handleGetFXRate(intent.Entities)

	case "create_payout":
		return a.handleCreatePayout(intent.Entities)

	case "help":
		return a.handleHelp()

	default:
		return a.handleUnknown()
	}
}

// handleCheckBalance processes balance check requests
func (a *Agent) handleCheckBalance() (string, error) {
	balance, err := a.client.GetBalance()
	if err != nil {
		return "", fmt.Errorf("error checking balance: %w", err)
	}

	var response strings.Builder
	response.WriteString("Your balances:\n")
	for _, bal := range balance.Data.Available {
		response.WriteString(fmt.Sprintf("- %s %s\n", bal.Amount, bal.Currency))
	}
	response.WriteString(fmt.Sprintf("Last updated: %s", balance.Data.UpdatedAt))

	return response.String(), nil
}

// handleGetFXRate processes exchange rate requests
func (a *Agent) handleGetFXRate(entities map[string]string) (string, error) {
	amount, _ := strconv.ParseFloat(entities["amount"], 64)
	rate, err := a.client.GetExchangeRate(
		entities["from_currency"],
		entities["to_currency"],
		amount,
	)
	if err != nil {
		return "", fmt.Errorf("error getting exchange rate: %w", err)
	}

	return fmt.Sprintf(
		"Exchange rate from %s to %s:\n"+
			"Amount: %.2f %s\n"+
			"Rate: %.6f\n"+
			"Converted: %.2f %s\n"+
			"Time: %s",
		rate.Data.InitialCurrency,
		rate.Data.FinalCurrency,
		float64(rate.Data.Amount)/100,
		rate.Data.InitialCurrency,
		rate.Data.ExchangeRate,
		float64(rate.Data.ConvertedAmount)/100,
		rate.Data.FinalCurrency,
		rate.Data.Timestamp,
	), nil
}

// handleCreatePayout processes payout creation requests
func (a *Agent) handleCreatePayout(entities map[string]string) (string, error) {
	amount, _ := strconv.ParseFloat(entities["amount"], 64)
	payout := &tazapay.Payout{
		Amount:          int64(amount * 100),
		Currency:        entities["currency"],
		Beneficiary:     entities["beneficiary_id"],
		HoldingCurrency: entities["currency"],
		Type:            "local",
		ChargeType:      "shared",
		Purpose:         "PYR002",
	}

	response, err := a.client.CreatePayout(payout)
	if err != nil {
		return "", fmt.Errorf("error creating payout: %w", err)
	}

	return fmt.Sprintf(
		"Payout created successfully!\n"+
			"ID: %s\n"+
			"Status: %s\n"+
			"Amount: %.2f %s\n"+
			"Created: %s",
		response.Data.ID,
		response.Data.Status,
		float64(response.Data.Amount)/100,
		response.Data.Currency,
		response.Data.CreatedAt,
	), nil
}

// handleHelp returns available commands and their usage
func (a *Agent) handleHelp() (string, error) {
	return "I can help you with:\n" +
		"1. Checking your balance\n" +
		"   Example: 'What's my balance?'\n\n" +
		"2. Getting exchange rates\n" +
		"   Example: 'What's the rate from USD to EUR for 100?'\n\n" +
		"3. Creating payouts\n" +
		"   Example: 'Send 100 SGD to beneficiary bnf_123'\n\n" +
		"Just ask your question in natural language!", nil
}

// handleUnknown handles unrecognized queries
func (a *Agent) handleUnknown() (string, error) {
	return "I'm not sure what you want to do. Try asking for 'help' to see what I can do!", nil
}
