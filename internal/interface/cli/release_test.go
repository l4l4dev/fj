package cli

import (
	"context"
	"strings"
	"testing"

	applicationrelease "github.com/l4l4dev/fj/internal/application/release"
)

type releaseListerStub struct {
	result []applicationrelease.Release
}

func (s releaseListerStub) List(context.Context, applicationrelease.ListRequest) ([]applicationrelease.Release, error) {
	return s.result, nil
}

func TestReleaseListOutput(t *testing.T) {
	command := newReleaseListCommand(releaseListerStub{result: []applicationrelease.Release{
		{ID: 1, TagName: "v1.0.0", Title: "First release", Draft: false, Prerelease: false},
		{ID: 2, TagName: "v1.1.0-rc1", Title: "Release candidate", Draft: false, Prerelease: true},
		{ID: 3, TagName: "v0.1.0", Title: "Work in progress", Draft: true, Prerelease: false},
	}})
	command.SetArgs([]string{"alice/project"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	want := "Releases:\nv1.0.0 First release [published]\nv1.1.0-rc1 Release candidate [prerelease]\nv0.1.0 Work in progress [draft]\n"
	if output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleaseListEmptyOutput(t *testing.T) {
	command := newReleaseListCommand(releaseListerStub{})
	command.SetArgs([]string{"alice/project"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if output.String() != "Releases:\nNo releases found.\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleaseListRejectsInvalidInput(t *testing.T) {
	tests := [][]string{
		{},
		{"invalid"},
		{"alice/project", "--page", "0"},
		{"alice/project", "--limit", "0"},
		{"alice/project", "--limit", "101"},
	}
	for _, args := range tests {
		command := newReleaseListCommand(releaseListerStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}
