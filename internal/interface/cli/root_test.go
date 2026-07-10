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
		want string
	}{
		{name: "unknown flag", args: []string{"--unknown"}, want: "unknown flag: --unknown"},
		{name: "unexpected argument", args: []string{"unexpected"}, want: "unknown command \"unexpected\" for \"fj\""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			command := NewRootCommand()
			command.SetOut(&output)
			command.SetErr(&output)
			command.SetArgs(test.args)

			err := command.Execute()
			if err == nil || err.Error() != test.want {
				t.Fatalf("Execute() error = %v, want %q", err, test.want)
			}
			if got := output.String(); got != "Error: "+test.want+"\n" {
				t.Errorf("error output = %q, want %q", got, "Error: "+test.want+"\n")
			}
		})
	}
}
