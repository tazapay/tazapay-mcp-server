package money

import (
	"testing"

	fmath "github.com/tazapay/tazapay-mcp-server/pkg/utils/math"
)

func TestInt64ToDecimal2(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected float64
	}{
		{"whole number", 100, 1.0},
		{"with cents", 123, 1.23},
		{"zero", 0, 0.0},
		{"large number", 1000000, 10000.0},
		{"negative", -123, -1.23},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Int64ToDecimal2(test.input)
			if result != test.expected {
				t.Errorf("Int64ToDecimal2(%d) = %f; want %f", test.input, result, test.expected)
			}
		})
	}
}

func TestDecimal2ToInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected int64
	}{
		{"whole number", 1.0, 100},
		{"with decimal", 1.23, 123},
		{"zero", 0.0, 0},
		{"large number", 10000.0, 1000000},
		{"negative", -1.23, -123},
		{"rounding up", 1.235, 124},
		{"rounding down", 1.234, 123},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Decimal2ToInt64(test.input)
			if result != test.expected {
				t.Errorf("Decimal2ToInt64(%f) = %d; want %d", test.input, result, test.expected)
			}
		})
	}
}

func TestFormatCurrency(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		currency string
		expected string
	}{
		{"USD", 1234, "USD", "$12.34"},
		{"EUR", 1234, "EUR", "€12.34"},
		{"GBP", 1234, "GBP", "£12.34"},
		{"INR", 1234, "INR", "₹12.34"},
		{"unknown", 1234, "XYZ", "12.34"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := FormatCurrency(test.amount, test.currency)
			if result != test.expected {
				t.Errorf("FormatCurrency(%d, %s) = %s; want %s",
					test.amount, test.currency, result, test.expected)
			}
		})
	}
}

func TestRoundingConsistency(t *testing.T) {
	// This test ensures that the conversion from float to int64 and back is consistent
	testValues := []float64{1.23, 45.67, 89.01, 99.99, 0.01}

	for _, value := range testValues {
		intValue := Decimal2ToInt64(value)
		floatValue := Int64ToDecimal2(intValue)

		// Should be the same after rounding
		expected := fmath.Round2Decimal(value)
		if floatValue != expected {
			t.Errorf("Rounding inconsistency: %f -> %d -> %f; expected %f",
				value, intValue, floatValue, expected)
		}
	}
}
