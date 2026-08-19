package cli

import (
	"context"
	"os"
	"path/filepath"
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

type releaseAssetUploaderStub struct {
	called  bool
	request applicationrelease.UploadAssetRequest
}

func (s *releaseAssetUploaderStub) UploadAsset(_ context.Context, request applicationrelease.UploadAssetRequest) (applicationrelease.Asset, error) {
	s.called = true
	s.request = request
	return applicationrelease.Asset{ID: 11, Name: request.AssetName, Size: int64(len(request.Content))}, nil
}

type releaseAssetDeleterStub struct {
	called  bool
	request applicationrelease.DeleteAssetRequest
}

func (s *releaseAssetDeleterStub) DeleteAsset(_ context.Context, request applicationrelease.DeleteAssetRequest) error {
	s.called = true
	s.request = request
	return nil
}

func writeReleaseAssetFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReleaseAssetListOutputWithAssets(t *testing.T) {
	inspector := &releaseInspectorStub{detail: applicationrelease.ReleaseDetail{
		TagName: "v1.0.0",
		Assets:  []applicationrelease.Asset{{ID: 11, Name: "fj_darwin_arm64.tar.gz", Size: 1024}, {ID: 12, Name: "checksums.txt", Size: 64}},
	}}
	command := newReleaseAssetListCommand(inspector)
	command.SetArgs([]string{"alice/project", "v1.0.0"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if inspector.request.Owner != "alice" || inspector.request.Name != "project" || inspector.request.Tag != "v1.0.0" {
		t.Fatalf("unexpected request: %+v", inspector.request)
	}
	want := "Assets for v1.0.0:\n- #11 fj_darwin_arm64.tar.gz (1024 bytes)\n- #12 checksums.txt (64 bytes)\n"
	if output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleaseAssetListEmptyOutput(t *testing.T) {
	command := newReleaseAssetListCommand(&releaseInspectorStub{detail: applicationrelease.ReleaseDetail{TagName: "v1.0.0"}})
	command.SetArgs([]string{"alice/project", "v1.0.0"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if output.String() != "Assets for v1.0.0: none\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleaseAssetListRejectsInvalidInput(t *testing.T) {
	tests := [][]string{
		{},
		{"alice/project"},
		{"invalid", "v1.0.0"},
		{"alice/project", "  "},
	}
	for _, args := range tests {
		command := newReleaseAssetListCommand(&releaseInspectorStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}

func TestReleaseAssetUploadUsesFileNameByDefault(t *testing.T) {
	uploader := &releaseAssetUploaderStub{}
	path := writeReleaseAssetFile(t, "fj_darwin_arm64.tar.gz", "payload")
	command := newReleaseAssetUploadCommand(uploader)
	command.SetArgs([]string{"alice/project", "v1.0.0", path})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if uploader.request.Owner != "alice" || uploader.request.Name != "project" || uploader.request.Tag != "v1.0.0" {
		t.Fatalf("unexpected request: %+v", uploader.request)
	}
	if uploader.request.AssetName != "fj_darwin_arm64.tar.gz" || string(uploader.request.Content) != "payload" {
		t.Fatalf("unexpected asset: %+v", uploader.request)
	}
	if output.String() != "Asset uploaded to v1.0.0: #11 fj_darwin_arm64.tar.gz (7 bytes)\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleaseAssetUploadHonoursNameFlag(t *testing.T) {
	uploader := &releaseAssetUploaderStub{}
	path := writeReleaseAssetFile(t, "build.tar.gz", "payload")
	command := newReleaseAssetUploadCommand(uploader)
	command.SetArgs([]string{"alice/project", "v1.0.0", path, "--name", "fj_darwin_arm64.tar.gz"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if uploader.request.AssetName != "fj_darwin_arm64.tar.gz" {
		t.Fatalf("unexpected asset name: %+v", uploader.request)
	}
	if output.String() != "Asset uploaded to v1.0.0: #11 fj_darwin_arm64.tar.gz (7 bytes)\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleaseAssetUploadReportsUnreadableFile(t *testing.T) {
	uploader := &releaseAssetUploaderStub{}
	command := newReleaseAssetUploadCommand(uploader)
	command.SetArgs([]string{"alice/project", "v1.0.0", filepath.Join(t.TempDir(), "missing.tar.gz")})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err == nil {
		t.Fatal("expected an error when the asset file cannot be read")
	}
	if uploader.called {
		t.Fatal("uploader must not be called when the asset file cannot be read")
	}
}

func TestReleaseAssetUploadRejectsInvalidInput(t *testing.T) {
	path := writeReleaseAssetFile(t, "build.tar.gz", "payload")
	tests := [][]string{
		{},
		{"alice/project", "v1.0.0"},
		{"invalid", "v1.0.0", path},
		{"alice/project", "  ", path},
	}
	for _, args := range tests {
		command := newReleaseAssetUploadCommand(&releaseAssetUploaderStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}

func TestReleaseAssetDeleteOutputAndRequest(t *testing.T) {
	deleter := &releaseAssetDeleterStub{}
	command := newReleaseAssetDeleteCommand(deleter)
	command.SetArgs([]string{"alice/project", "v1.0.0", "fj_darwin_arm64.tar.gz", "--confirm"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if deleter.request.Owner != "alice" || deleter.request.Name != "project" || deleter.request.Tag != "v1.0.0" || deleter.request.AssetName != "fj_darwin_arm64.tar.gz" {
		t.Fatalf("unexpected request: %+v", deleter.request)
	}
	if output.String() != "Asset deleted from v1.0.0: fj_darwin_arm64.tar.gz\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleaseAssetDeleteRequiresConfirm(t *testing.T) {
	deleter := &releaseAssetDeleterStub{}
	command := newReleaseAssetDeleteCommand(deleter)
	command.SetArgs([]string{"alice/project", "v1.0.0", "fj_darwin_arm64.tar.gz"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err == nil {
		t.Fatal("expected an error when --confirm is absent")
	}
	if deleter.called {
		t.Fatal("deleter must not be called without --confirm")
	}
	if strings.Contains(output.String(), "Asset deleted from") {
		t.Fatalf("expected no success output on failure, got %q", output.String())
	}
}

func TestReleaseAssetDeleteRejectsInvalidInput(t *testing.T) {
	tests := [][]string{
		{},
		{"alice/project", "v1.0.0", "--confirm"},
		{"invalid", "v1.0.0", "fj.tar.gz", "--confirm"},
		{"alice/project", "  ", "fj.tar.gz", "--confirm"},
		{"alice/project", "v1.0.0", "  ", "--confirm"},
	}
	for _, args := range tests {
		command := newReleaseAssetDeleteCommand(&releaseAssetDeleterStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}

type releaseDeleterStub struct {
	called  bool
	request applicationrelease.DeleteRequest
}

func (s *releaseDeleterStub) Delete(_ context.Context, request applicationrelease.DeleteRequest) error {
	s.called = true
	s.request = request
	return nil
}

func TestReleaseDeleteOutputAndRequest(t *testing.T) {
	inspector := &releaseInspectorStub{detail: applicationrelease.ReleaseDetail{ID: 7, TagName: "v1.0.0", Title: "First release"}}
	deleter := &releaseDeleterStub{}
	command := newReleaseDeleteCommand(inspector, deleter)
	command.SetArgs([]string{"alice/project", "v1.0.0", "--confirm"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if !deleter.called || deleter.request.ID != 7 || deleter.request.Owner != "alice" || deleter.request.Name != "project" {
		t.Fatalf("unexpected request: %+v", deleter.request)
	}
	want := "Release deleted: v1.0.0 (First release)\nThe git tag v1.0.0 was not deleted.\n"
	if output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestReleaseDeleteRequiresConfirm(t *testing.T) {
	inspector := &releaseInspectorStub{detail: applicationrelease.ReleaseDetail{ID: 7, TagName: "v1.0.0", Title: "First release"}}
	deleter := &releaseDeleterStub{}
	command := newReleaseDeleteCommand(inspector, deleter)
	command.SetArgs([]string{"alice/project", "v1.0.0"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err == nil {
		t.Fatal("expected an error when --confirm is absent")
	}
	if deleter.called {
		t.Fatal("deleter must not be called without --confirm")
	}
	if inspector.request.Tag != "" {
		t.Fatal("inspector must not be called without --confirm")
	}
	if strings.Contains(output.String(), "Release deleted:") {
		t.Fatalf("expected no success output on failure, got %q", output.String())
	}
}

func TestReleaseDeleteRejectsInvalidInput(t *testing.T) {
	tests := [][]string{
		{},
		{"alice/project"},
		{"invalid", "v1.0.0", "--confirm"},
		{"alice/project", "  ", "--confirm"},
	}
	for _, args := range tests {
		command := newReleaseDeleteCommand(&releaseInspectorStub{}, &releaseDeleterStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}
