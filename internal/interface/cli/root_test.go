package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

func TestRootCommandDisplaysHelp(t *testing.T) {
	var output bytes.Buffer
	command := NewRootCommand()
	command.SetOut(&output)
	command.SetArgs([]string{"--help"})

	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}

	help := output.String()
	for _, expected := range []string{"fj", "Usage:"} {
		if !strings.Contains(help, expected) {
			t.Errorf("help output does not contain %q", expected)
		}
	}
}

func TestRootCommandRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "unknown flag", args: []string{"--unknown"}},
		{name: "unexpected argument", args: []string{"unexpected"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			command := NewRootCommand()
			command.SetOut(&stdout)
			command.SetErr(&stderr)

			if code := Execute(command, test.args); code != categoryValidation.exitCode() {
				t.Fatalf("Execute() code = %d, want %d", code, categoryValidation.exitCode())
			}
			if got := stdout.String(); got != "" {
				t.Errorf("standard output = %q, want empty", got)
			}
			if got := stderr.String(); got != "Error: execute command: invalid input\n" {
				t.Errorf("standard error = %q", got)
			}
		})
	}
}

func TestMapApplicationErrorSurfacesValidationMessage(t *testing.T) {
	err := mapApplicationError(apperror.NewValidation("submit pull request review", "outcome must be comment, approve, or request-changes"), "submit pull request review")
	want := "submit pull request review: outcome must be comment, approve, or request-changes"
	if err == nil || err.Error() != want {
		t.Fatalf("unexpected error: %v", err)
	}

	fallback := mapApplicationError(apperror.ValidationError{Operation: "submit pull request review"}, "submit pull request review")
	if fallback == nil || fallback.Error() != "submit pull request review: invalid input" {
		t.Fatalf("unexpected fallback error: %v", fallback)
	}
}
