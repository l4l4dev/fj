package apperror

import "fmt"

type Category uint8

const (
	Internal Category = iota
	Validation
	Authentication
	NotFound
	Conflict
	Remote
)

type Error struct {
	Category  Category
	Operation string
	Message   string
}

func (err Error) Error() string {
	if err.Message != "" {
		return fmt.Sprintf("%s: %s", err.Operation, err.Message)
	}
	return fmt.Sprintf("%s: %s", err.Operation, message(err.Category))
}

type ValidationError struct {
	Operation string
	Message   string
}

func (err ValidationError) Error() string { return err.Message }
func NewValidation(operation, message string) error {
	return ValidationError{Operation: operation, Message: message}
}
func New(category Category, operation, message string) error {
	return Error{Category: category, Operation: operation, Message: message}
}

func message(category Category) string {
	switch category {
	case Validation:
		return "invalid input"
	case Authentication:
		return "authentication failed"
	case NotFound:
		return "repository not found"
	case Conflict:
		return "operation conflict"
	case Remote:
		return "remote operation failed"
	default:
		return "internal error"
	}
}
