package apperror

import (
	"errors"
	"testing"
)

func TestErrorCategoriesAreClassifiable(t *testing.T) {
	err := New(NotFound, "inspect repository", "repository not found")
	var classified Error
	if !errors.As(err, &classified) || classified.Category != NotFound {
		t.Fatalf("error = %v", err)
	}
}

func TestValidationErrorIsClassifiable(t *testing.T) {
	err := NewValidation("list repositories", "page must be at least 1")
	var classified ValidationError
	if !errors.As(err, &classified) || classified.Message != "page must be at least 1" {
		t.Fatalf("error = %v", err)
	}
}
