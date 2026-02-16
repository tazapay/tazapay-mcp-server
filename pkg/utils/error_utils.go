//nolint:sloglint // slog attributes can be used
package utils

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/tazapay/tazapay-mcp-server/constants"
)

// WrapFieldTypeError creates a custom error for invalid field type and logs it.
func WrapFieldTypeError(ctx context.Context, logger *slog.Logger, field string) error {
	// Format the error message with wrapping
	err := fmt.Errorf("%w: %s", constants.ErrInvalidType, field)

	// Log the error with field info using ErrorContext
	logger.ErrorContext(ctx, "field type validation failed",
		slog.String("field", field),
		slog.String("error", err.Error()),
	)

	return err
}

// WrapInvalidAmountError creates a custom error for invalid amount format and returns it.
func WrapInvalidAmountError(currency string) error {
	// Use static error wrapping
	return fmt.Errorf("%w: %s", constants.ErrInvalidAmountFormat, currency)
}

// Missing Fields Error
func WrapMissingFieldsError(fields []string) error {
	return fmt.Errorf("%w: %s", constants.ErrMissingRequiredFields, strings.Join(fields, ", "))
}
