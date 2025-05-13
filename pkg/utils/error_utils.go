package utils

import (
	"fmt"
	"log/slog"

	"github.com/tazapay/tazapay-mcp-server/constants"
)

// WrapFieldTypeError creates a custom error for invalid field type and logs it.
func WrapFieldTypeError(logger *slog.Logger, field string) error {
	// Format the error message with wrapping
	err := fmt.Errorf("%w: %s", constants.ErrInvalidType, field)

	// Log the error with field info
	logger.Error("field type validation failed", slog.String("field", field), slog.String("error", err.Error()))

	return err
}

// WrapInvalidAmountError creates a custom error for invalid amount format and returns it.
func WrapInvalidAmountError(currency string) error {
	// Using fmt.Errorf for better formatting
	return fmt.Errorf("invalid amount format for currency: %s", currency)
}
