package cli

import "testing"

func TestInjectedVersionIsAvailableToCompositionRoot(t *testing.T) {
	command := NewRootCommandWithVersion(RepositoryDependencies{}, "1.2.3")
	if got := versionFromContext(command.Context()); got != "1.2.3" {
		t.Fatalf("composition version = %q, want %q", got, "1.2.3")
	}
}
