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

type releaseInspectorStub struct {
	request applicationrelease.InspectRequest
	detail  applicationrelease.ReleaseDetail
}

func (s *releaseInspectorStub) Inspect(_ context.Context, request applicationrelease.InspectRequest) (applicationrelease.ReleaseDetail, error) {
	s.request = request
	return s.detail, nil
}

func TestReleaseInspectOutputWithAssets(t *testing.T) {
	inspector := &releaseInspectorStub{detail: applicationrelease.ReleaseDetail{
		TagName: "v1.0.0",
		Title:   "First release",
		Notes:   "Release notes",
		Assets:  []applicationrelease.Asset{{ID: 1, Name: "fj_darwin_arm64.tar.gz", Size: 1024}},
	}}
	command := newReleaseInspectCommand(inspector)
	command.SetArgs([]string{"alice/project", "v1.0.0"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if inspector.request.Owner != "alice" || inspector.request.Name != "project" || inspector.request.Tag != "v1.0.0" {
		t.Fatalf("unexpected request: %+v", inspector.request)
	}
	want := "Release: v1.0.0\nTitle: First release\nState: published\nNotes: Release notes\nAssets:\n- #1 fj_darwin_arm64.tar.gz (1024 bytes)\n"
	if output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleaseInspectOutputWithoutAssets(t *testing.T) {
	inspector := &releaseInspectorStub{detail: applicationrelease.ReleaseDetail{TagName: "v0.1.0", Title: "Draft", Draft: true}}
	command := newReleaseInspectCommand(inspector)
	command.SetArgs([]string{"alice/project", "v0.1.0"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	want := "Release: v0.1.0\nTitle: Draft\nState: draft\nNotes: -\nAssets: none\n"
	if output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleaseInspectRejectsInvalidInput(t *testing.T) {
	tests := [][]string{
		{},
		{"alice/project"},
		{"invalid", "v1.0.0"},
		{"alice/project", "  "},
	}
	for _, args := range tests {
		command := newReleaseInspectCommand(&releaseInspectorStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}

type releaseCreatorStub struct {
	request applicationrelease.CreateRequest
}

func (s *releaseCreatorStub) Create(_ context.Context, request applicationrelease.CreateRequest) (applicationrelease.ReleaseDetail, error) {
	s.request = request
	return applicationrelease.ReleaseDetail{TagName: request.Tag, Title: request.Title, Draft: true, Prerelease: request.Prerelease}, nil
}

func TestReleaseCreateOutputAndRequest(t *testing.T) {
	creator := &releaseCreatorStub{}
	command := newReleaseCreateCommand(creator)
	command.SetArgs([]string{"alice/project", "--tag", "v1.0.0", "--title", "First release", "--notes", "Notes", "--prerelease"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if creator.request.Owner != "alice" || creator.request.Name != "project" || creator.request.Tag != "v1.0.0" || creator.request.Title != "First release" || creator.request.Notes != "Notes" || !creator.request.Prerelease {
		t.Fatalf("unexpected request: %+v", creator.request)
	}
	want := "Release created as draft: v1.0.0\nTitle: First release\nPrerelease: true\n"
	if output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleaseCreateRejectsInvalidInput(t *testing.T) {
	tests := [][]string{
		{"invalid", "--tag", "v1.0.0", "--title", "First release"},
		{"alice/project", "--title", "First release"},
		{"alice/project", "--tag", "v 1", "--title", "First release"},
		{"alice/project", "--tag", "v1.0.0"},
		{"alice/project", "--tag", "v1.0.0", "--title", "  "},
	}
	for _, args := range tests {
		command := newReleaseCreateCommand(&releaseCreatorStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}

type releaseUpdaterStub struct {
	request applicationrelease.UpdateRequest
	detail  applicationrelease.ReleaseDetail
}

func (stub *releaseUpdaterStub) Update(_ context.Context, request applicationrelease.UpdateRequest) (applicationrelease.ReleaseDetail, error) {
	stub.request = request
	if stub.detail.TagName == "" {
		stub.detail.TagName = request.Tag
	}
	return stub.detail, nil
}

func TestReleaseUpdateSendsOnlyChangedTitleField(t *testing.T) {
	updater := &releaseUpdaterStub{}
	command := newReleaseUpdateCommand(updater)
	command.SetArgs([]string{"alice/project", "v1.0.0", "--title", "New title"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if updater.request.Title == nil || *updater.request.Title != "New title" || updater.request.Notes != nil || updater.request.Prerelease != nil {
		t.Fatalf("unexpected request: %+v", updater.request)
	}
	if output.String() != "Release updated: v1.0.0\nChanged fields: title\nState: published\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleaseUpdateSendsOnlyChangedPrereleaseField(t *testing.T) {
	updater := &releaseUpdaterStub{detail: applicationrelease.ReleaseDetail{TagName: "v1.0.0", Prerelease: false}}
	command := newReleaseUpdateCommand(updater)
	command.SetArgs([]string{"alice/project", "v1.0.0", "--prerelease=false"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if updater.request.Title != nil || updater.request.Notes != nil || updater.request.Prerelease == nil || *updater.request.Prerelease != false {
		t.Fatalf("unexpected request: %+v", updater.request)
	}
	if output.String() != "Release updated: v1.0.0\nChanged fields: prerelease\nState: published\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleaseUpdateRejectsInvalidInput(t *testing.T) {
	tests := [][]string{
		{"alice/project", "v1.0.0"},
		{"alice/project", "  ", "--title", "New title"},
		{"invalid", "v1.0.0", "--title", "New title"},
		{"alice/project"},
	}
	for _, args := range tests {
		command := newReleaseUpdateCommand(&releaseUpdaterStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}

type releasePublisherStub struct {
	called  bool
	request applicationrelease.PublishRequest
	detail  applicationrelease.ReleaseDetail
}

func (s *releasePublisherStub) Publish(_ context.Context, request applicationrelease.PublishRequest) (applicationrelease.ReleaseDetail, error) {
	s.called = true
	s.request = request
	return s.detail, nil
}

func TestReleasePublishOutputAndRequest(t *testing.T) {
	inspector := &releaseInspectorStub{detail: applicationrelease.ReleaseDetail{ID: 7, TagName: "v1.0.0", Title: "First release", Draft: true}}
	publisher := &releasePublisherStub{detail: applicationrelease.ReleaseDetail{ID: 7, TagName: "v1.0.0", Title: "First release", Draft: false}}
	command := newReleasePublishCommand(inspector, publisher)
	command.SetArgs([]string{"alice/project", "v1.0.0"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if !publisher.called || publisher.request.ID != 7 || publisher.request.Owner != "alice" || publisher.request.Name != "project" {
		t.Fatalf("unexpected request: %+v", publisher.request)
	}
	if output.String() != "Release published: v1.0.0\nTitle: First release\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleasePublishRejectsAlreadyPublished(t *testing.T) {
	inspector := &releaseInspectorStub{detail: applicationrelease.ReleaseDetail{ID: 7, TagName: "v1.0.0", Draft: false}}
	publisher := &releasePublisherStub{}
	command := newReleasePublishCommand(inspector, publisher)
	command.SetArgs([]string{"alice/project", "v1.0.0"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err == nil {
		t.Fatal("expected an error when the release is already published")
	}
	if publisher.called {
		t.Fatal("publisher must not be called when the release is already published")
	}
	if strings.Contains(output.String(), "Release published:") {
		t.Fatalf("expected no success output on failure, got %q", output.String())
	}
}

func TestReleasePublishRejectsInvalidInput(t *testing.T) {
	tests := [][]string{
		{"alice/project"},
		{"invalid", "v1.0.0"},
		{"alice/project", "  "},
		{},
	}
	for _, args := range tests {
		command := newReleasePublishCommand(&releaseInspectorStub{}, &releasePublisherStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}
