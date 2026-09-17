// Package apperrors defines sentinel, service-layer errors that are
// translated into HTTP status codes at the handler layer, keeping the
// domain/service layers free of any HTTP-specific concerns.
package apperrors

import "errors"

var (
	// ErrNotFound indicates the requested resource does not exist.
	ErrNotFound = errors.New("resource not found")
	// ErrConflict indicates a uniqueness or state conflict (e.g. duplicate SKU).
	ErrConflict = errors.New("resource conflict")
	// ErrValidation indicates the request failed business/input validation.
	ErrValidation = errors.New("validation failed")
	// ErrInsufficientStock indicates a purchase requested more units than available in inventory.
	ErrInsufficientStock = errors.New("insufficient stock")
)

// Field describes a single validation failure, used to build detailed
// validation error responses.
type Field struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationError wraps ErrValidation with per-field details.
type ValidationError struct {
	Fields []Field
}

func (e *ValidationError) Error() string {
	return ErrValidation.Error()
}

func (e *ValidationError) Unwrap() error {
	return ErrValidation
}

func NewValidationError(fields []Field) *ValidationError {
	return &ValidationError{Fields: fields}
}

// InsufficientStockError wraps ErrInsufficientStock with the specific
// product and quantities involved, so handlers/clients can display a
// precise message.
type InsufficientStockError struct {
	ProductID string
	Requested int
	Available int
}

func (e *InsufficientStockError) Error() string {
	return ErrInsufficientStock.Error()
}

func (e *InsufficientStockError) Unwrap() error {
	return ErrInsufficientStock
}

func NewInsufficientStockError(productID string, requested, available int) *InsufficientStockError {
	return &InsufficientStockError{ProductID: productID, Requested: requested, Available: available}
}
