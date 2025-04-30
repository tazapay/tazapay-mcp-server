package agent

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tazapay/tazapay-mcp-server/pkg/tazapay"
)

// Context represents the conversation context
type Context struct {
	LastIntent    string
	LastEntities  map[string]string
	MissingFields map[string]bool
	Conversation  []Message
}

// Message represents a message in the conversation
type Message struct {
	Role    string
	Content string
	Time    time.Time
}

// Agent represents an AI agent that handles natural language interactions
type Agent struct {
	client  *tazapay.Client
	context *Context
}

// NewAgent creates a new AI agent
func NewAgent(client *tazapay.Client) *Agent {
	return &Agent{
		client: client,
		context: &Context{
			LastEntities:  make(map[string]string),
			MissingFields: make(map[string]bool),
			Conversation:  make([]Message, 0),
		},
	}
}

// HandleMessage processes a natural language message and returns a response
func (a *Agent) HandleMessage(ctx context.Context, message string) (string, error) {
	// Add message to conversation history
	a.context.Conversation = append(a.context.Conversation, Message{
		Role:    "user",
		Content: message,
		Time:    time.Now(),
	})

	// Analyze intent and extract entities
	intent, entities := a.analyzeMessage(message)
	a.context.LastIntent = intent
	a.context.LastEntities = entities

	// Handle the intent
	response, err := a.handleIntent(intent, entities)
	if err != nil {
		// Log error to MCP server
		a.logError(err)
		return "", err
	}

	// Add response to conversation history
	a.context.Conversation = append(a.context.Conversation, Message{
		Role:    "assistant",
		Content: response,
		Time:    time.Now(),
	})

	return response, nil
}

// analyzeMessage performs sophisticated NLP analysis
func (a *Agent) analyzeMessage(message string) (string, map[string]string) {
	msg := strings.ToLower(message)
	entities := make(map[string]string)

	// Intent detection using regex patterns
	intent := a.detectIntent(msg)

	// Entity extraction
	entities = a.extractEntities(msg)

	// Context-based entity resolution
	a.resolveEntitiesFromContext(entities)

	return intent, entities
}

// detectIntent identifies the user's intent using regex patterns
func (a *Agent) detectIntent(message string) string {
	patterns := map[string]*regexp.Regexp{
		"check_balance":      regexp.MustCompile(`(?:check|what is|how much|show me).*(?:balance|money|funds)`),
		"create_beneficiary": regexp.MustCompile(`(?:add|create|new).*(?:beneficiary|recipient|receiver)`),
		"create_payout":      regexp.MustCompile(`(?:send|transfer|make|create).*(?:payout|payment|money)`),
		"check_fx":           regexp.MustCompile(`(?:exchange rate|convert|fx rate|rate).*(?:from|to)`),
	}

	for intent, pattern := range patterns {
		if pattern.MatchString(message) {
			return intent
		}
	}

	return "unknown"
}

// extractEntities performs sophisticated entity extraction
func (a *Agent) extractEntities(message string) map[string]string {
	entities := make(map[string]string)

	// Currency extraction with support for both codes and names
	currencyPattern := regexp.MustCompile(`(?i)(?:USD|EUR|GBP|JPY|AUD|CAD|CHF|CNY|HKD|SGD|dollar|euro|pound|yen)`)
	currencies := currencyPattern.FindAllString(message, -1)
	if len(currencies) > 0 {
		entities["currency"] = normalizeCurrency(currencies[0])
	}

	// Amount extraction with support for different formats
	amountPattern := regexp.MustCompile(`(?:\$|€|£)?\s*(\d+(?:,\d{3})*(?:\.\d{2})?)`)
	amounts := amountPattern.FindAllStringSubmatch(message, -1)
	if len(amounts) > 0 {
		entities["amount"] = strings.ReplaceAll(amounts[0][1], ",", "")
	}

	// Email extraction
	emailPattern := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	emails := emailPattern.FindAllString(message, -1)
	if len(emails) > 0 {
		entities["email"] = emails[0]
	}

	// Country extraction
	countryPattern := regexp.MustCompile(`(?i)(?:in|from|to|country)\s+([A-Za-z]+(?: [A-Za-z]+)*)`)
	countries := countryPattern.FindAllStringSubmatch(message, -1)
	if len(countries) > 0 {
		entities["country"] = strings.TrimSpace(countries[0][1])
	}

	return entities
}

// resolveEntitiesFromContext fills in missing entities from previous context
func (a *Agent) resolveEntitiesFromContext(entities map[string]string) {
	for key, value := range a.context.LastEntities {
		if _, exists := entities[key]; !exists {
			entities[key] = value
		}
	}
}

// handleIntent processes the detected intent with extracted entities
func (a *Agent) handleIntent(intent string, entities map[string]string) (string, error) {
	switch intent {
	case "check_balance":
		return a.handleCheckBalance(entities)
	case "create_beneficiary":
		return a.handleCreateBeneficiary(entities)
	case "create_payout":
		return a.handleCreatePayout(entities)
	case "check_fx":
		return a.handleCheckFX(entities)
	default:
		return a.handleUnknownIntent()
	}
}

// handleCheckBalance processes balance check requests
func (a *Agent) handleCheckBalance(entities map[string]string) (string, error) {
	balance, err := a.client.GetBalance()
	if err != nil {
		return "", fmt.Errorf("error checking balance: %w", err)
	}

	return fmt.Sprintf("Your balance is %.2f %s", balance.Data.Balance, balance.Data.Currency), nil
}

// handleCreateBeneficiary processes beneficiary creation requests
func (a *Agent) handleCreateBeneficiary(entities map[string]string) (string, error) {
	missing := make([]string, 0)
	for _, field := range []string{"name", "email", "country"} {
		if _, ok := entities[field]; !ok {
			missing = append(missing, field)
			a.context.MissingFields[field] = true
		}
	}

	if len(missing) > 0 {
		return fmt.Sprintf("Please provide the following information: %s", strings.Join(missing, ", ")), nil
	}

	beneficiary := &tazapay.Beneficiary{
		Name:    entities["name"],
		Email:   entities["email"],
		Country: entities["country"],
	}
	response, err := a.client.CreateBeneficiary(beneficiary)
	if err != nil {
		return "", fmt.Errorf("error creating beneficiary: %w", err)
	}

	return fmt.Sprintf("Beneficiary created successfully! ID: %s", response), nil
}

// handleCreatePayout processes payout creation requests
func (a *Agent) handleCreatePayout(entities map[string]string) (string, error) {
	missing := make([]string, 0)
	for _, field := range []string{"beneficiary_id", "amount", "currency"} {
		if _, ok := entities[field]; !ok {
			missing = append(missing, field)
			a.context.MissingFields[field] = true
		}
	}

	if len(missing) > 0 {
		return fmt.Sprintf("Please provide the following information: %s", strings.Join(missing, ", ")), nil
	}

	amount, _ := strconv.ParseFloat(entities["amount"], 64)
	payout := &tazapay.Payout{
		BeneficiaryID: entities["beneficiary_id"],
		Amount:        amount,
		Currency:      entities["currency"],
	}
	response, err := a.client.CreatePayout(payout)
	if err != nil {
		return "", fmt.Errorf("error creating payout: %w", err)
	}

	return fmt.Sprintf("Payout created successfully! ID: %s", response), nil
}

// handleCheckFX processes FX rate requests
func (a *Agent) handleCheckFX(entities map[string]string) (string, error) {
	missing := make([]string, 0)
	for _, field := range []string{"from_currency", "to_currency"} {
		if _, ok := entities[field]; !ok {
			missing = append(missing, field)
			a.context.MissingFields[field] = true
		}
	}

	if len(missing) > 0 {
		return "Please specify both source and target currencies (e.g., USD to EUR)", nil
	}

	rate, err := a.client.GetExchangeRate(
		entities["from_currency"],
		entities["to_currency"],
	)
	if err != nil {
		return "", fmt.Errorf("error getting exchange rate: %w", err)
	}

	return fmt.Sprintf("Exchange rate from %s to %s: %.4f",
		entities["from_currency"],
		entities["to_currency"],
		rate,
	), nil
}

// handleUnknownIntent handles unrecognized intents
func (a *Agent) handleUnknownIntent() (string, error) {
	return "I can help you with:\n" +
		"- Checking account balance (e.g., 'What's my USD balance?')\n" +
		"- Creating beneficiaries (e.g., 'Add a new beneficiary named John')\n" +
		"- Making payouts (e.g., 'Send 100 USD to beneficiary 123')\n" +
		"- Checking exchange rates (e.g., 'What's the rate from USD to EUR?')\n\n" +
		"Please ask your question in natural language.", nil
}

// logError logs errors to the MCP server
func (a *Agent) logError(err error) {
	// TODO: Implement error logging to MCP server
	fmt.Printf("Error logged: %v\n", err)
}

// normalizeCurrency normalizes currency names to codes
func normalizeCurrency(currency string) string {
	currencyMap := map[string]string{
		"dollar": "USD",
		"euro":   "EUR",
		"pound":  "GBP",
		"yen":    "JPY",
	}

	if code, ok := currencyMap[strings.ToLower(currency)]; ok {
		return code
	}

	return strings.ToUpper(currency)
}
