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

func TestIssueCreateUsesInspectOutput(t *testing.T) {
	command := newIssueCreateCommand(creatorStubForCLI{detail: applicationissue.IssueDetail{Number: 13, Title: "Created", State: applicationissue.StateOpen}})
	command.SetArgs([]string{"alice/project", "--title", "Created"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if want := "Issue: #13\nTitle: Created\nState: open\nBody: -\n"; output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestIssueUpdateOutputAndEmptyBody(t *testing.T) {
	command := newIssueUpdateCommand(updaterStubForCLI{detail: applicationissue.IssueDetail{Number: 12, Title: "Updated", State: applicationissue.StateOpen}})
	command.SetArgs([]string{"alice/project", "12", "--body", ""})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if want := "Issue: #12\nChanged fields: body\nTitle: Updated\nState: open\nBody: -\n"; output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestIssueStateOutput(t *testing.T) {
	command := newIssueStateCommand(stateChangerStubForCLI{detail: applicationissue.IssueDetail{Number: 12, State: applicationissue.StateClosed}})
	command.SetArgs([]string{"alice/project", "12", "--state", "closed"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if want := "Issue: #12\nState: closed\n"; output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestIssueCommentOutputs(t *testing.T) {
	list := newIssueCommentListCommand(commentViewerStubForCLI{comments: []applicationissue.Comment{{ID: 1, Body: "hello"}}})
	list.SetArgs([]string{"alice/project", "12"})
	var listOutput strings.Builder
	list.SetOut(&listOutput)
	if err := list.Execute(); err != nil {
		t.Fatal(err)
	}
	if want := "Comments:\n- #1 hello\n"; listOutput.String() != want {
		t.Fatalf("unexpected list output: %q", listOutput.String())
	}
	add := newIssueCommentAddCommand(commentCreatorStubForCLI{comment: applicationissue.Comment{ID: 2, Body: "hello"}})
	add.SetArgs([]string{"alice/project", "12", "--body", "hello"})
	var addOutput strings.Builder
	add.SetOut(&addOutput)
	if err := add.Execute(); err != nil {
		t.Fatal(err)
	}
	if want := "Comment:\n#2 hello\n"; addOutput.String() != want {
		t.Fatalf("unexpected add output: %q", addOutput.String())
	}
}

type commentViewerStubForCLI struct{ comments []applicationissue.Comment }

func (stub commentViewerStubForCLI) ListComments(context.Context, applicationissue.ListCommentsRequest) ([]applicationissue.Comment, error) {
	return stub.comments, nil
}

type commentCreatorStubForCLI struct{ comment applicationissue.Comment }

func (stub commentCreatorStubForCLI) AddComment(context.Context, applicationissue.AddCommentRequest) (applicationissue.Comment, error) {
	return stub.comment, nil
}

type inspectorStubForCLI struct{ detail applicationissue.IssueDetail }

func (stub inspectorStubForCLI) Inspect(context.Context, applicationissue.InspectRequest) (applicationissue.IssueDetail, error) {
	return stub.detail, nil
}

type creatorStubForCLI struct{ detail applicationissue.IssueDetail }

func (stub creatorStubForCLI) Create(context.Context, applicationissue.CreateRequest) (applicationissue.IssueDetail, error) {
	return stub.detail, nil
}

type updaterStubForCLI struct{ detail applicationissue.IssueDetail }

func (stub updaterStubForCLI) Update(context.Context, applicationissue.UpdateRequest) (applicationissue.IssueDetail, error) {
	return stub.detail, nil
}

type stateChangerStubForCLI struct{ detail applicationissue.IssueDetail }

func (stub stateChangerStubForCLI) ChangeState(context.Context, applicationissue.ChangeStateRequest) (applicationissue.IssueDetail, error) {
	return stub.detail, nil
}
