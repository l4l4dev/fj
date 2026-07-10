package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestExecuteReturnsDistinctSafeProcessOutcomes(t *testing.T) {
	const secret = "secret-token"
	tests := []struct {
		name       string
		category   errorCategory
		wantCode   int
		wantOutput string
	}{
		{name: "internal", category: categoryInternal, wantCode: 1, wantOutput: "Error: list repositories: internal error\n"},
		{name: "validation", category: categoryValidation, wantCode: 2, wantOutput: "Error: list repositories: invalid input\n"},
		{name: "authentication", category: categoryAuthentication, wantCode: 3, wantOutput: "Error: list repositories: authentication failed\n"},
		{name: "remote", category: categoryRemote, wantCode: 4, wantOutput: "Error: list repositories: remote operation failed\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stderr bytes.Buffer
			command := &cobra.Command{
				Use:           "test",
				SilenceErrors: true,
				SilenceUsage:  true,
				RunE: func(*cobra.Command, []string) error {
					return newCommandError(test.category, "list repositories", errors.New(secret))
				},
			}
			command.SetErr(&stderr)

			if code := Execute(command, nil); code != test.wantCode {
				t.Fatalf("Execute() code = %d, want %d", code, test.wantCode)
			}
			if got := stderr.String(); got != test.wantOutput {
				t.Errorf("standard error = %q, want %q", got, test.wantOutput)
			}
			if strings.Contains(stderr.String(), secret) {
				t.Errorf("standard error exposes sensitive cause: %q", stderr.String())
			}
		})
	}
}

func TestExecuteHidesUnclassifiedInternalCause(t *testing.T) {
	const secret = "secret internal detail"
	var stderr bytes.Buffer
	command := &cobra.Command{
		Use:           "test",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(*cobra.Command, []string) error {
			return errors.New(secret)
		},
	}
	command.SetErr(&stderr)

	if code := Execute(command, nil); code != categoryInternal.exitCode() {
		t.Fatalf("Execute() code = %d, want %d", code, categoryInternal.exitCode())
	}
	if got := stderr.String(); got != "Error: execute command: internal error\n" {
		t.Errorf("standard error = %q", got)
	}
	if strings.Contains(stderr.String(), secret) {
		t.Errorf("standard error exposes internal cause: %q", stderr.String())
	}
}
