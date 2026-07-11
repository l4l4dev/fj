package cli

import (
	"context"
	"strings"
	"testing"

	applicationissue "github.com/l4l4dev/fj/internal/application/issue"
)

type issueListerStub struct{ page applicationissue.Page }

func (stub issueListerStub) List(context.Context, applicationissue.ListRequest) (applicationissue.Page, error) {
	return stub.page, nil
}

func TestIssueListPresenterOutput(t *testing.T) {
	command := newIssueListCommand(issueListerStub{page: applicationissue.Page{Issues: []applicationissue.Issue{{Number: 12, Title: "Fix authentication flow", State: applicationissue.StateOpen}}, Page: 1, Limit: 30}})
	var output strings.Builder
	command.SetOut(&output)
	command.SetArgs([]string{"alice/project"})
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	want := "Issues:\n- #12 Fix authentication flow [open]\n\nPage: 1\nLimit: 30\nMore pages: false\n"
	if output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestIssueListRejectsInvalidState(t *testing.T) {
	command := newIssueListCommand(issueListerStub{})
	command.SetArgs([]string{"alice/project", "--state", "invalid"})
	if err := command.Execute(); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestIssueListRejectsCombinedAndRepeatedFilters(t *testing.T) {
	for _, args := range [][]string{{"alice/project", "--assignee", "bob", "--label", "bug"}, {"alice/project", "--label", "bug", "--label", "feature"}} {
		command := newIssueListCommand(issueListerStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for args %v", args)
		}
	}
}

func TestIssueInspectOutput(t *testing.T) {
	command := newIssueInspectCommand(inspectorStubForCLI{detail: applicationissue.IssueDetail{Number: 12, Title: "Fix it", State: applicationissue.StateOpen, Body: "Details"}})
	command.SetArgs([]string{"alice/project", "12"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if want := "Issue: #12\nTitle: Fix it\nState: open\nBody: Details\n"; output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

type inspectorStubForCLI struct{ detail applicationissue.IssueDetail }

func (stub inspectorStubForCLI) Inspect(context.Context, applicationissue.InspectRequest) (applicationissue.IssueDetail, error) {
	return stub.detail, nil
}
