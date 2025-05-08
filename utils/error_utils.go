package utils

import (
	"fmt"

	"github.com/tazapay/tazapay-mcp-server/constants"
)

func WrapFieldTypeError(field string) error {
	return fmt.Errorf("%w: %s", constants.ErrInvalidType, field)
}
