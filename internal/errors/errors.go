package errors

import (
	"errors"
	"fmt"
)

// Sentinel errors - define common errors as package-level variables
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrGroupNotFound    = errors.New("group not found")
	ErrActivityNotFound = errors.New("activity not found")
	ErrInvalidTimezone  = errors.New("invalid timezone")
	ErrInvalidActivity  = errors.New("invalid activity format")
	ErrDatabase         = errors.New("database error")
	ErrTelegramAPI      = errors.New("telegram api error")
)

// Custom error types for more specific error handling

// ValidationError represents validation failures
type ValidationError struct {
	Field   string
	Value   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s=%q: %s", e.Field, e.Value, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, value, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// DatabaseError wraps database-related errors with additional context
type DatabaseError struct {
	Operation string
	Table     string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database %s failed on %s: %v", e.Operation, e.Table, e.Err)
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// NewDatabaseError creates a new database error
func NewDatabaseError(operation, table string, err error) *DatabaseError {
	return &DatabaseError{
		Operation: operation,
		Table:     table,
		Err:       err,
	}
}

// TelegramError wraps Telegram API errors
type TelegramError struct {
	Method string
	ChatID int64
	Err    error
}

func (e *TelegramError) Error() string {
	return fmt.Sprintf("telegram %s failed for chat %d: %v", e.Method, e.ChatID, e.Err)
}

func (e *TelegramError) Unwrap() error {
	return e.Err
}

// NewTelegramError creates a new Telegram error
func NewTelegramError(method string, chatID int64, err error) *TelegramError {
	return &TelegramError{
		Method: method,
		ChatID: chatID,
		Err:    err,
	}
}