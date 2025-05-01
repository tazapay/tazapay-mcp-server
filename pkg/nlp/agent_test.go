package nlp

import (
	"testing"

	"tazapay-mcp-server/pkg/tazapay"
)

func TestDetectIntent(t *testing.T) {
	agent := NewAgent(tazapay.NewClient("test", "test"))

	tests := []struct {
		query          string
		expectedIntent string
	}{
		{"what is my balance", "check_balance"},
		{"show me my funds", "check_balance"},
		{"check balance", "check_balance"},
		{"what's the exchange rate from usd to eur", "get_fx_rate"},
		{"convert 100 usd to eur", "get_fx_rate"},
		{"send 100 sgd to john", "create_payout"},
		{"make a payment of 50 usd", "create_payout"},
		{"help me", "help"},
		{"what can you do", "help"},
		{"random text", "unknown"},
	}

	for _, test := range tests {
		intent := agent.detectIntent(test.query)
		if intent.Name != test.expectedIntent {
			t.Errorf("Query '%s': expected intent '%s', got '%s'",
				test.query, test.expectedIntent, intent.Name)
		}
	}
}

func TestExtractEntities(t *testing.T) {
	agent := NewAgent(tazapay.NewClient("test", "test"))

	tests := []struct {
		query    string
		expected map[string]string
	}{
		{
			"convert 100 USD to EUR",
			map[string]string{
				"amount":        "100",
				"from_currency": "USD",
				"to_currency":   "EUR",
			},
		},
		{
			"send 50.5 SGD to John",
			map[string]string{
				"amount":           "50.5",
				"from_currency":    "SGD",
				"beneficiary_name": "John",
			},
		},
		{
			"what's my balance",
			map[string]string{},
		},
	}

	for _, test := range tests {
		entities := agent.extractEntities(test.query)
		for key, expected := range test.expected {
			if got := entities[key]; got != expected {
				t.Errorf("Query '%s', entity '%s': expected '%s', got '%s'",
					test.query, key, expected, got)
			}
		}
	}
}

func TestGetMissingEntities(t *testing.T) {
	agent := NewAgent(tazapay.NewClient("test", "test"))

	tests := []struct {
		intent          Intent
		expectedMissing []string
	}{
		{
			Intent{
				Name: "get_fx_rate",
				Entities: map[string]string{
					"amount": "100",
				},
			},
			[]string{"source currency", "target currency"},
		},
		{
			Intent{
				Name: "create_payout",
				Entities: map[string]string{
					"amount":   "100",
					"currency": "USD",
				},
			},
			[]string{"beneficiary ID"},
		},
		{
			Intent{
				Name:     "check_balance",
				Entities: map[string]string{},
			},
			nil,
		},
	}

	for _, test := range tests {
		missing := agent.getMissingEntities(test.intent)
		if len(missing) != len(test.expectedMissing) {
			t.Errorf("Intent '%s': expected %d missing entities, got %d",
				test.intent.Name, len(test.expectedMissing), len(missing))
			continue
		}
		for i, expected := range test.expectedMissing {
			if missing[i] != expected {
				t.Errorf("Intent '%s': expected missing entity '%s', got '%s'",
					test.intent.Name, expected, missing[i])
			}
		}
	}
}
