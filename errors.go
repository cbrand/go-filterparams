package filterparams

import (
	"fmt"
)

// UnsupportedOperationError indicates that an operation was passed
// which is unsupported.
type UnsupportedOperationError struct {
	Operation string
}

// Error returns the formatted error message.
func (u *UnsupportedOperationError) Error() string {
	return fmt.Sprintf("The operation \"%s\" is unsupported", u.Operation)
}

// NewUnsupportedOperation generates the error with the passed operation as
// indication.
func NewUnsupportedOperation(operation string) *UnsupportedOperationError {
	return &UnsupportedOperationError{
		Operation: operation,
	}
}
