package cli

import (
	"bytes"
	"strings"
	"testing"
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
