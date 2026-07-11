package cli

import (
	"strings"
	"testing"
)

func TestVersionCommandPrintsVersionOnly(t *testing.T) {
	command := NewRootCommandWithVersion(RepositoryDependencies{}, "1.2.3")
	var output strings.Builder
	command.SetOut(&output)
	command.SetArgs([]string{"version"})

	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "1.2.3\n"; got != want {
		t.Fatalf("version output = %q, want %q", got, want)
	}
}

func TestVersionFlagPrintsVersionOnly(t *testing.T) {
	command := NewRootCommandWithVersion(RepositoryDependencies{}, "1.2.3")
	var output strings.Builder
	command.SetOut(&output)
	command.SetArgs([]string{"--version"})

	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "1.2.3\n"; got != want {
		t.Fatalf("version output = %q, want %q", got, want)
	}
}

func TestVersionCommandRejectsArguments(t *testing.T) {
	command := NewRootCommandWithVersion(RepositoryDependencies{}, "1.2.3")
	command.SetArgs([]string{"version", "extra"})

	if err := command.Execute(); err == nil {
		t.Fatal("version command accepted an argument")
	}
}

func TestDefaultVersionIsDev(t *testing.T) {
	if got, want := defaultVersion(), "dev"; got != want {
		t.Fatalf("defaultVersion() = %q, want %q", got, want)
	}
}
