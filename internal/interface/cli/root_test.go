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
